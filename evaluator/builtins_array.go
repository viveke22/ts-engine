package evaluator

import (
	"strings"
	"ts-engine/object"
)

func arrayProperty(arr *object.Array, name string) object.Object {
	switch name {
	case "length":
		return &object.Integer{Value: float64(len(arr.Elements))}
	case "push":
		return &object.Builtin{Fn: func(args ...object.Object) object.Object {
			arr.Elements = append(arr.Elements, args...)
			return &object.Integer{Value: float64(len(arr.Elements))}
		}}
	case "pop":
		return &object.Builtin{Fn: func(args ...object.Object) object.Object {
			n := len(arr.Elements)
			if n == 0 {
				return UNDEFINED
			}
			last := arr.Elements[n-1]
			arr.Elements = arr.Elements[:n-1]
			return last
		}}
	case "map":
		return &object.Builtin{Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return newError("map expects a function")
			}
			mapped := make([]object.Object, 0, len(arr.Elements))
			for i, el := range arr.Elements {
				idx := &object.Integer{Value: float64(i)}
				result := applyFunction(args[0], []object.Object{el, idx, arr})
				if isError(result) {
					return result
				}
				mapped = append(mapped, result)
			}
			return &object.Array{Elements: mapped}
		}}
	case "join":
		return &object.Builtin{Fn: func(args ...object.Object) object.Object {
			sep := ","
			if len(args) > 0 && args[0] != UNDEFINED {
				if s, ok := args[0].(*object.String); ok {
					sep = s.Value
				} else {
					sep = args[0].Inspect()
				}
			}
			parts := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				parts[i] = joinElement(el)
			}
			return &object.String{Value: strings.Join(parts, sep)}
		}}
	default:
		return UNDEFINED
	}
}

func joinElement(obj object.Object) string {
	if obj == nil || obj == NULL || obj == UNDEFINED {
		return ""
	}
	if s, ok := obj.(*object.String); ok {
		return s.Value
	}
	return obj.Inspect()
}
