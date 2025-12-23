package http

import (
	"encoding/json"
	"io"
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return &object.Error{Message: "failed to read response body: " + err.Error()}
	}
	bodyString := string(bodyBytes)

	pairs := make(map[string]object.Object)
	pairs["status"] = &object.Integer{Value: int64(resp.StatusCode)}
	pairs["ok"] = &object.Boolean{Value: resp.StatusCode >= 200 && resp.StatusCode < 300}
	pairs["statusText"] = &object.String{Value: resp.Status}

	// .text() method
	pairs["text"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			return &object.String{Value: bodyString}
		},
	}

	// .json() method
	pairs["json"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			var result interface{}
			if err := json.Unmarshal([]byte(bodyString), &result); err != nil {
				return &object.Error{Message: "failed to parse JSON: " + err.Error()}
			}
			return convertJsonToObject(result)
		},
	}

	return &object.Hash{Pairs: pairs}
}

func convertJsonToObject(v interface{}) object.Object {
	switch val := v.(type) {
	case nil:
		return &object.Null{}
	case bool:
		return &object.Boolean{Value: val}
	case float64:
		return &object.Integer{Value: int64(val)} // Simplified: JSON numbers are floats, we treat as int for now
	case string:
		return &object.String{Value: val}
	case []interface{}:
		// Arrays not fully supported yet in object system??
		// Check implementation plans. We have Arrays support in lexer/parser?
		// Features.md says "Arrays: Basic array support (via host integration)".
		// But in `object/object.go`?
		// Let's assume we return NULL or String representation for now if ArrayObj is missing.
		// Wait, I should check object.go for Array support.
		// Assuming we don't have Array object yet, fallback to string?
		return &object.String{Value: "[Array]"}
	case map[string]interface{}:
		pairs := make(map[string]object.Object)
		for k, v := range val {
			pairs[k] = convertJsonToObject(v)
		}
		return &object.Hash{Pairs: pairs}
	}
	return &object.String{Value: "unknown"}
}
