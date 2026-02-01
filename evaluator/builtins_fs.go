package evaluator

import (
	"os"
	"ts-engine/object"
)

func createFsModule() *object.Hash {
	pairs := make(map[string]object.Object)

	pairs["writeFileSync"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("writeFileSync expects 2 arguments (path, content)")
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument 1 must be string")
			}
			content, ok := args[1].(*object.String)
			if !ok {
				return newError("argument 2 must be string")
			}

			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("fs error: %s", err)
			}
			return NULL
		},
	}

	pairs["readFileSync"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("readFileSync expects 1 argument (path)")
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument 1 must be string")
			}

			data, err := os.ReadFile(path.Value)
			if err != nil {
				return newError("fs error: %s", err)
			}
			return &object.String{Value: string(data)}
		},
	}

	pairs["removeSync"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("removeSync expects 1 argument (path)")
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument 1 must be string")
			}

			err := os.Remove(path.Value)
			if err != nil {
				return newError("fs error: %s", err)
			}
			return NULL
		},
	}

	pairs["existsSync"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("existsSync expects 1 argument (path)")
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument 1 must be string")
			}

			_, err := os.Stat(path.Value)
			if err == nil {
				return TRUE
			}
			if os.IsNotExist(err) {
				return FALSE
			}
			return FALSE // Error implies not accessible/exist usually
		},
	}

	return &object.Hash{Pairs: pairs}
}
