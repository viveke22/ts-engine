package evaluator

import (
	"fmt"
	"strings"
	"ts-engine/ast"
	"ts-engine/http"
	"ts-engine/object"
	"ts-engine/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		// Special handling for dot operator to avoid evaluating the property as a variable
		if node.Operator == "." {
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			return evalDotIndexExpression(left, node.Right)
		}

		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.ExportStatement:
		return NULL

	case *ast.ImportStatement:
		return evalImportStatement(node, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		if node.Token.Type == token.VAR {
			// VAR allows redeclaration, so we don't check GetCurrent
			// Ideally VAR is function-scoped, but for now we treat it as block-scoped or whatever env is.
		} else {
			// LET and CONST do not allow redeclaration
			if _, ok := env.GetCurrent(node.Name.Value); ok {
				return newError("cannot redeclare block-scoped variable '%s'", node.Name.Value)
			}
		}

		if node.Name.Type != "" {
			if err := checkType(val, node.Name.Type); err != nil {
				return err
			}
		}

		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		fn := &object.Function{Parameters: params, Env: env, Body: body}
		if node.Name != "" {
			env.Set(node.Name, fn)
		}
		return fn

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	case "await":
		return right
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "+" && (left.Type() == object.STRING_OBJ || right.Type() == object.STRING_OBJ):
		return evalStringConcatenation(left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case operator == "===":
		return nativeBoolToBooleanObject(left == right && left.Type() == right.Type())
	case operator == "!==":
		return nativeBoolToBooleanObject(left != right || left.Type() != right.Type())
	case operator == "!==":
		return nativeBoolToBooleanObject(left != right || left.Type() != right.Type())
	case operator == "&&":
		return nativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))
	case operator == "||":
		return nativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalImportStatement(node *ast.ImportStatement, env *object.Environment) object.Object {
	source := node.Source.Value

	// We reuse the 'require' builtin logic
	requireObj, ok := builtins["require"]
	if !ok {
		return newError("require builtin not found")
	}

	requireFn, ok := requireObj.(*object.Builtin)
	if !ok {
		return newError("require is not a function")
	}

	// Call require("source")
	module := requireFn.Fn(&object.String{Value: source})
	if isError(module) {
		return module
	}

	// Bind result to alias
	env.Set(node.Alias.Value, module)

	return NULL
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[string]object.Object)

	for keyNode, valueNode := range node.Pairs {
		var keyStr string

		// If key is Identifier, take the name as string literal (e.g. { name: "val" })
		if ident, ok := keyNode.(*ast.Identifier); ok {
			keyStr = ident.Value
		} else {
			key := Eval(keyNode, env)
			if isError(key) {
				return key
			}
			if s, ok := key.(*object.String); ok {
				keyStr = s.Value
			} else {
				keyStr = key.Inspect() // Fallback
			}
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		pairs[keyStr] = value
	}

	return &object.Hash{Pairs: pairs}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "===":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!==":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "===":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!==":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalDotIndexExpression(left object.Object, rightNode ast.Node) object.Object {
	if left.Type() == object.HASH_OBJ {
		hash := left.(*object.Hash)

		// Right node should be an identifier for dot notation
		ident, ok := rightNode.(*ast.Identifier)
		if !ok {
			return newError("expected identifier after dot, got %T", rightNode)
		}

		key := ident.Value
		val, ok := hash.Pairs[key]
		if !ok {
			return NULL
		}
		return val
	}
	return newError("property access not supported on %s", left.Type())
}

func evalStringConcatenation(left, right object.Object) object.Object {
	var leftVal, rightVal string

	if left.Type() == object.STRING_OBJ {
		leftVal = left.(*object.String).Value
	} else {
		leftVal = left.Inspect()
	}

	if right.Type() == object.STRING_OBJ {
		rightVal = right.(*object.String).Value
	} else {
		rightVal = right.Inspect()
	}

	return &object.String{Value: leftVal + rightVal}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := evalBlockStatement(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

var builtins map[string]object.Object

func init() {
	builtins = map[string]object.Object{
		"console": &object.Hash{
			Pairs: map[string]object.Object{
				"log": &object.Builtin{
					Fn: func(args ...object.Object) object.Object {
						var out []string
						for _, arg := range args {
							out = append(out, arg.Inspect())
						}
						fmt.Println(strings.Join(out, " "))
						return NULL
					},
				},
			},
		},
		"fetch": &object.Builtin{
			Fn: http.Fetch,
		},
		"require": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				name, ok := args[0].(*object.String)
				if !ok {
					return newError("argument to `require` must be STRING, got %s", args[0].Type())
				}

				if name.Value == "http" {
					// Return the HTTP module object
					return &object.Hash{
						Pairs: map[string]object.Object{
							"createServer": &object.Builtin{
								Fn: createHttpServer,
							},
						},
					}
				}

				return newError("module not found: %s", name.Value)
			},
		},
	}
}

func checkType(obj object.Object, typeName string) *object.Error {
	switch typeName {
	case "any", "unknown":
		return nil
	case "number":
		if obj.Type() != object.INTEGER_OBJ {
			return newError("type mismatch: expected number, got %s", obj.Type())
		}
	case "string":
		if obj.Type() != object.STRING_OBJ {
			return newError("type mismatch: expected string, got %s", obj.Type())
		}
	case "boolean":
		if obj.Type() != object.BOOLEAN_OBJ {
			return newError("type mismatch: expected boolean, got %s", obj.Type())
		}
	case "never":
		return newError("type mismatch: cannot assign to never")
	default:
		// If type contains '.', treat as named interface/class and allow (for now)
		// e.g. http.IncomingMessage
		// Also allow simple identifiers that are not primitives (e.g. MyClass)
		return nil
	}
	return nil
}
