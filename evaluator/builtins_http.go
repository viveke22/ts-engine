package evaluator

import (
	"fmt"
	"net/http"
	"strconv"
	"ts-engine/object"
)

func createHttpServer(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	handlerFn, ok := args[0].(*object.Function)
	if !ok {
		return newError("argument to `createServer` must be a function, got %s", args[0].Type())
	}

	// Create a Hash object representing the server
	serverMap := make(map[string]object.Object)

	// "listen" method
	serverMap["listen"] = &object.Builtin{
		Fn: func(listenArgs ...object.Object) object.Object {
			// Expected args: (port, callback?)
			if len(listenArgs) < 1 {
				return newError("listen expects at least 1 argument (port)")
			}

			portVal := listenArgs[0]
			var port int
			if intObj, ok := portVal.(*object.Integer); ok {
				port = int(intObj.Value)
			} else {
				return newError("port must be integer")
			}

			// Optional listen callback
			var listenCb *object.Function
			if len(listenArgs) > 1 {
				if fn, ok := listenArgs[1].(*object.Function); ok {
					listenCb = fn
				}
			}

			addr := ":" + strconv.Itoa(port)

			// Execute listen callback immediately just to simulate "started" if valid
			// In Node, this happens after 'listening' event. We'll just call it before blocking or async.
			// But ListenAndServe blocks. So we should probably call it before?
			// But if it fails, we shouldn't have called it.
			// Ideally, run ListenAndServe in goroutine.

			go func() {
				// We need a way to keep main alive?
				// For now, let's assume valid server script ends with .listen() and we want to BLOCK.
				// If we run in goroutine, the script finishes and program exits.
				// The user provided code has .listen() at the end.
			}()

			// If we block, we can't run the callback?
			// The callback is usually "Server running...".
			// If we block, we never run it.

			// Simple hack: Run callback, then block.
			if listenCb != nil {
				applyFunction(listenCb, []object.Object{})
			}

			fmt.Printf("Starting server on %s...\n", addr)
			err := http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 1. Convert Request
				tsReq := &object.Hash{
					Pairs: map[string]object.Object{
						"url":    &object.String{Value: r.URL.String()},
						"method": &object.String{Value: r.Method},
					},
				}

				// 2. Wrap Response
				// We need methods: writeHead, end
				// We can't use a simple Hash because it needs methods that close over 'w'.
				// But we can return a Hash full of Builtins!

				tsRes := &object.Hash{
					Pairs: map[string]object.Object{
						"writeHead": &object.Builtin{
							Fn: func(args ...object.Object) object.Object {
								if len(args) < 1 {
									return NULL
								}
								status := 200
								if s, ok := args[0].(*object.Integer); ok {
									status = int(s.Value)
								}

								// Handle headers (arg 1)
								if len(args) > 1 {
									if headers, ok := args[1].(*object.Hash); ok {
										for key, val := range headers.Pairs {
											// We only support String or Integer values for headers for now
											if strVal, ok := val.(*object.String); ok {
												w.Header().Set(key, strVal.Value)
											} else if intVal, ok := val.(*object.Integer); ok {
												w.Header().Set(key, strconv.Itoa(int(intVal.Value)))
											}
										}
									}
								}

								w.WriteHeader(status)
								return NULL
							},
						},
						"end": &object.Builtin{
							Fn: func(args ...object.Object) object.Object {
								// Support calling end with data: res.end("data")
								if len(args) > 0 {
									if s, ok := args[0].(*object.String); ok {
										w.Write([]byte(s.Value))
									} else {
										// Fallback for non-string, e.g. integer or just Inspect
										w.Write([]byte(args[0].Inspect()))
									}
								}
								return NULL
							},
						},
					},
				}

				// 3. Call Handler
				applyFunction(handlerFn, []object.Object{tsReq, tsRes})
			}))

			if err != nil {
				return newError("server error: %s", err)
			}
			return NULL
		},
	}

	return &object.Hash{Pairs: serverMap}
}
