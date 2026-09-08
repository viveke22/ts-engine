package evaluator

import (
	"strings"
	"testing"
	"ts-engine/lexer"
	"ts-engine/object"
	"ts-engine/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l, true)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}

func TestPropertyAssignment(t *testing.T) {
	result := testEval(`
		let obj: any = { x: 0, a: { b: 0 } };
		obj.x = 1;
		obj.a.b = 2;
		obj.x + obj.a.b;
	`)
	intEquals(t, result, 3)

	err := testEval(`
		let n: number = 1;
		n.x = 2;
	`)
	isErrorContaining(t, err, "cannot assign property")
}

func TestIndexAssignment(t *testing.T) {
	result := testEval(`
		let arr: number[] = [1, 2];
		arr[0] = 9;
		arr[2] = 3;
		arr[0] + arr[1] + arr[2];
	`)
	intEquals(t, result, 14)

	padded := testEval(`
		let arr: any = [1];
		arr[3] = 4;
		arr[1] === undefined && arr[2] === undefined && arr[3] === 4;
	`)
	boolEquals(t, padded, true)

	neg := testEval(`
		let arr: number[] = [1];
		arr[-1] = 0;
	`)
	isErrorContaining(t, neg, "negative")

	hash := testEval(`
		let obj: any = {};
		obj["k"] = 5;
		obj["k"];
	`)
	intEquals(t, hash, 5)
}

func TestBreakContinue(t *testing.T) {
	brk := testEval(`
		let s: number = 0;
		for (let i: number = 0; i < 10; i++) {
			if (i == 3) {
				break;
			}
			s = s + i;
		}
		s;
	`)
	intEquals(t, brk, 3)

	cont := testEval(`
		let s: number = 0;
		for (let i: number = 0; i < 5; i++) {
			if (i == 2) {
				continue;
			}
			s = s + i;
		}
		s;
	`)
	intEquals(t, cont, 8)

	wbreak := testEval(`
		let i: number = 0;
		while (i < 10) {
			i = i + 1;
			if (i == 2) {
				break;
			}
		}
		i;
	`)
	intEquals(t, wbreak, 2)

	wcont := testEval(`
		let i: number = 0;
		let s: number = 0;
		while (i < 5) {
			i = i + 1;
			if (i == 3) {
				continue;
			}
			s = s + i;
		}
		s;
	`)
	intEquals(t, wcont, 12)

	outside := testEval(`break;`)
	isErrorContaining(t, outside, "break outside of loop")

	outsideCont := testEval(`continue;`)
	isErrorContaining(t, outsideCont, "continue outside of loop")

	ifBreak := testEval(`if (true) { break; }`)
	isErrorContaining(t, ifBreak, "break outside of loop")

	ifCont := testEval(`if (true) { continue; }`)
	isErrorContaining(t, ifCont, "continue outside of loop")

	fnBreak := testEval(`
		function f() {
			break;
		}
		f();
	`)
	isErrorContaining(t, fnBreak, "break outside of loop")

	nested := testEval(`
		let outerCount: number = 0;
		for (let i: number = 0; i < 3; i++) {
			for (let j: number = 0; j < 3; j++) {
				break;
			}
			outerCount = outerCount + 1;
		}
		outerCount;
	`)
	intEquals(t, nested, 3)
}

func TestConst(t *testing.T) {
	ok := testEval(`
		const n: number = 1;
		n;
	`)
	intEquals(t, ok, 1)

	reassign := testEval(`
		const n: number = 1;
		n = 2;
	`)
	isErrorContaining(t, reassign, "cannot assign to const")

	inc := testEval(`
		const n: number = 1;
		n++;
	`)
	isErrorContaining(t, inc, "cannot assign to const")

	mutable := testEval(`
		const obj: any = { x: 1 };
		obj.x = 2;
		obj.x;
	`)
	intEquals(t, mutable, 2)

	letOk := testEval(`
		let n: number = 1;
		n = 2;
		n;
	`)
	intEquals(t, letOk, 2)

	outerConstInnerReassign := testEval(`
		const c: number = 10;
		function f() {
			c = 20;
		}
		f();
	`)
	isErrorContaining(t, outerConstInnerReassign, "cannot assign to const")

	shadowing := testEval(`
		const x: number = 10;
		function f(): number {
			let x: number = 20;
			x = 30;
			return x;
		}
		f() + x;
	`)
	intEquals(t, shadowing, 40)

	varRedeclare := testEval(`
		var v: number = 1;
		var v: number = 2;
		v;
	`)
	intEquals(t, varRedeclare, 2)

	missingInit := testEval(`const uninit;`)
	isErrorContaining(t, missingInit, "missing initializer")
}

func TestUndefinedNull(t *testing.T) {
	u := testEval(`let x; x;`)
	if u != UNDEFINED {
		t.Fatalf("uninitialized let should be undefined, got %s", inspect(u))
	}

	typed := testEval(`let y: number; y;`)
	if typed != UNDEFINED {
		t.Fatalf("uninitialized typed let should be undefined, got %s", inspect(typed))
	}

	n := testEval(`null;`)
	if n != NULL {
		t.Fatalf("expected null, got %s", inspect(n))
	}

	eq := testEval(`null === null && undefined === undefined && !(null === undefined) && !(null == undefined);`)
	boolEquals(t, eq, true)

	neq := testEval(`(null != undefined) && (null !== undefined) && !(null != null) && !(undefined != undefined);`)
	boolEquals(t, neq, true)

	truth := testEval(`!null && !undefined;`)
	boolEquals(t, truth, true)

	falsyIf := testEval(`
		let n: number = 0;
		if (null) { n = 1; }
		if (undefined) { n = 1; }
		n;
	`)
	intEquals(t, falsyIf, 0)

	ann := testEval(`
		let a: null = null;
		let b: undefined = undefined;
		a === null && b === undefined;
	`)
	boolEquals(t, ann, true)

	typeErrNull := testEval(`let x: number = null;`)
	isErrorContaining(t, typeErrNull, "type mismatch")

	typeErrUndefined := testEval(`let x: number = undefined;`)
	isErrorContaining(t, typeErrUndefined, "type mismatch")

	typeErrAssign := testEval(`let x: null = 1;`)
	isErrorContaining(t, typeErrAssign, "type mismatch")
}

func TestForLoopBodyScope(t *testing.T) {
	result := testEval(`
		let sum: number = 0;
		for (let i: number = 0; i < 4; i++) {
			let factor: number = 2;
			const offset: number = 1;
			sum = sum + (i * factor) + offset;
		}
		sum;
	`)
	intEquals(t, result, 16)

	shadow := testEval(`
		let total: number = 0;
		for (let i: number = 0; i < 3; i++) {
			let i: number = 10;
			total = total + i;
		}
		total;
	`)
	intEquals(t, shadow, 30)
}

func TestEmptyReturnUndefined(t *testing.T) {
	result := testEval(`
		function f() {
			return;
		}
		f();
	`)
	if result != UNDEFINED {
		t.Fatalf("expected undefined from empty return, got %s", inspect(result))
	}

	arrowResult := testEval(`
		let g = () => {
			return;
		};
		g();
	`)
	if arrowResult != UNDEFINED {
		t.Fatalf("expected undefined from arrow empty return, got %s", inspect(arrowResult))
	}
}

func TestNestedIndexAndPropertyAssignment(t *testing.T) {
	nestedArr := testEval(`
		let grid: any = [[0, 0], [0, 0]];
		grid[0][1] = 42;
		grid[0][1];
	`)
	intEquals(t, nestedArr, 42)

	objArr := testEval(`
		let data: any = { items: [1, 2] };
		data.items[1] = 99;
		data.items[1];
	`)
	intEquals(t, objArr, 99)

	arrObj := testEval(`
		let list: any = [{ val: 10 }];
		list[0].val = 20;
		list[0].val;
	`)
	intEquals(t, arrObj, 20)
}

func TestArrayMethods(t *testing.T) {
	pushed := testEval(`
		let a: any = [1];
		let n: number = a.push(2, 3);
		n === 3 && a[2] === 3;
	`)
	boolEquals(t, pushed, true)

	popped := testEval(`
		let a: any = [1, 2];
		let x: any = a.pop();
		x === 2 && a.length === 1;
	`)
	boolEquals(t, popped, true)

	emptyPop := testEval(`
		let a: any = [];
		a.pop();
	`)
	if emptyPop != UNDEFINED {
		t.Fatalf("pop empty expected undefined, got %s", inspect(emptyPop))
	}

	mapped := testEval(`
		let a: any = [1, 2, 3];
		let b: any = a.map((x: number) => x * 2);
		b[0] + b[1] + b[2];
	`)
	intEquals(t, mapped, 12)

	joined := testEval(`
		let a: any = [1, null, undefined, "x"];
		a.join("-");
	`)
	stringEquals(t, joined, "1---x")

	joinedDefault := testEval(`
		let a: any = [1, 2];
		a.join();
	`)
	stringEquals(t, joinedDefault, "1,2")

	joinedUndef := testEval(`
		let a: any = [1, 2];
		a.join(undefined);
	`)
	stringEquals(t, joinedUndef, "1,2")
}

func TestTemplateInterpolation(t *testing.T) {
	plain := testEval("let s: string = `hello`; s;")
	stringEquals(t, plain, "hello")

	interp := testEval("let n: number = 2; `a ${n + 1} b`;")
	stringEquals(t, interp, "a 3 b")

	multi := testEval("let x: string = `x`; let y: string = `y`; `${x}-${y}`;")
	stringEquals(t, multi, "x-y")

	nested := testEval("let n: number = 1; `a ${`b ${n}`} c`;")
	stringEquals(t, nested, "a b 1 c")

	nullish := testEval("`n=${null} u=${undefined}`;")
	stringEquals(t, nullish, "n=null u=undefined")

	braced := testEval("`val=${(() => { return 42; })()}`;")
	stringEquals(t, braced, "val=42")

	chained := testEval("`list=${[1, 2, 3].map((x: number) => x * 10).join(\"-\")}`;")
	stringEquals(t, chained, "list=10-20-30")
}

func TestFunctionFewArguments(t *testing.T) {
	result := testEval(`
		function testArgs(a: any, b: any) {
			return b;
		}
		testArgs(10);
	`)
	if result != UNDEFINED {
		t.Fatalf("expected undefined for missing argument, got %s", inspect(result))
	}
}

func intEquals(t *testing.T, obj object.Object, expected float64) {
	t.Helper()
	i, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("expected Integer %g, got %s (%T)", expected, inspect(obj), obj)
	}
	if i.Value != expected {
		t.Fatalf("expected %g, got %g", expected, i.Value)
	}
}

func boolEquals(t *testing.T, obj object.Object, expected bool) {
	t.Helper()
	b, ok := obj.(*object.Boolean)
	if !ok {
		t.Fatalf("expected Boolean %t, got %s (%T)", expected, inspect(obj), obj)
	}
	if b.Value != expected {
		t.Fatalf("expected %t, got %t", expected, b.Value)
	}
}

func stringEquals(t *testing.T, obj object.Object, expected string) {
	t.Helper()
	s, ok := obj.(*object.String)
	if !ok {
		t.Fatalf("expected String %q, got %s (%T)", expected, inspect(obj), obj)
	}
	if s.Value != expected {
		t.Fatalf("expected %q, got %q", expected, s.Value)
	}
}

func isErrorContaining(t *testing.T, obj object.Object, sub string) {
	t.Helper()
	err, ok := obj.(*object.Error)
	if !ok {
		t.Fatalf("expected Error containing %q, got %s (%T)", sub, inspect(obj), obj)
	}
	if !strings.Contains(err.Message, sub) {
		t.Fatalf("error %q does not contain %q", err.Message, sub)
	}
}

func inspect(obj object.Object) string {
	if obj == nil {
		return "nil"
	}
	return obj.Inspect()
}
