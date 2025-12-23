package http

import (
	"net/http"
	"strconv"
	"ts-engine/object"
)

func Fetch(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{Message: "wrong number of arguments. got=" + strconv.Itoa(len(args)) + ", want=1"}
	}

	url, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{Message: "argument to `fetch` must be STRING, got " + string(args[0].Type())}
	}

	resp, err := http.Get(url.Value)
	if err != nil {
		return &object.Error{Message: "http error: " + err.Error()}
	}
	defer resp.Body.Close()

	pairs := make(map[string]object.Object)
	pairs["status"] = &object.Integer{Value: int64(resp.StatusCode)}
	pairs["ok"] = &object.Boolean{Value: resp.StatusCode >= 200 && resp.StatusCode < 300}
	pairs["statusText"] = &object.String{Value: resp.Status}

	return &object.Hash{Pairs: pairs}
}
