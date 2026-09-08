package parser

import (
	"testing"
	"ts-engine/ast"
	"ts-engine/lexer"
)

func TestArrowInfix(t *testing.T) {
	input := "x => x"
	l := lexer.New(input)
	p := New(l, false)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Errorf("parser has %d errors", len(p.Errors()))
		for _, msg := range p.Errors() {
			t.Errorf("parser error: %s", msg)
		}
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	_, ok = stmt.Expression.(*ast.ArrowFunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ArrowFunctionLiteral. got=%T",
			stmt.Expression)
	}
}

func TestBreakContinueNullUndefined(t *testing.T) {
	input := `
		break;
		continue;
		null;
		undefined;
		obj.x = 1;
		arr[0] = 2;
	`
	l := lexer.New(input)
	p := New(l, true)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}
	if len(program.Statements) != 6 {
		t.Fatalf("got %d statements", len(program.Statements))
	}
	if _, ok := program.Statements[0].(*ast.BreakStatement); !ok {
		t.Fatalf("stmt 0 not BreakStatement: %T", program.Statements[0])
	}
	if _, ok := program.Statements[1].(*ast.ContinueStatement); !ok {
		t.Fatalf("stmt 1 not ContinueStatement: %T", program.Statements[1])
	}
	es, ok := program.Statements[2].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt 2 not ExpressionStatement: %T", program.Statements[2])
	}
	if _, ok := es.Expression.(*ast.NullLiteral); !ok {
		t.Fatalf("stmt 2 not NullLiteral: %T", es.Expression)
	}
	es, ok = program.Statements[3].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt 3 not ExpressionStatement: %T", program.Statements[3])
	}
	if _, ok := es.Expression.(*ast.UndefinedLiteral); !ok {
		t.Fatalf("stmt 3 not UndefinedLiteral: %T", es.Expression)
	}
	es, ok = program.Statements[4].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt 4 not ExpressionStatement: %T", program.Statements[4])
	}
	assign, ok := es.Expression.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("stmt 4 not AssignmentExpression: %T", es.Expression)
	}
	if infix, ok := assign.Left.(*ast.InfixExpression); !ok || infix.Operator != "." {
		t.Fatalf("assignment left not dot infix: %T", assign.Left)
	}
	es, ok = program.Statements[5].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt 5 not ExpressionStatement: %T", program.Statements[5])
	}
	assign, ok = es.Expression.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("stmt 5 not AssignmentExpression: %T", es.Expression)
	}
	if _, ok := assign.Left.(*ast.IndexExpression); !ok {
		t.Fatalf("assignment left not IndexExpression: %T", assign.Left)
	}
}

func TestTemplateLiteral(t *testing.T) {
	input := "`hello ${x + 1}!`"
	l := lexer.New(input)
	p := New(l, true)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}
	if len(program.Statements) != 1 {
		t.Fatalf("got %d statements", len(program.Statements))
	}
	es, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not ExpressionStatement: %T", program.Statements[0])
	}
	tl, ok := es.Expression.(*ast.TemplateLiteral)
	if !ok {
		t.Fatalf("not TemplateLiteral: %T", es.Expression)
	}
	if len(tl.Quasis) != 2 || tl.Quasis[0] != "hello " || tl.Quasis[1] != "!" {
		t.Fatalf("quasis=%v", tl.Quasis)
	}
	if len(tl.Expressions) != 1 {
		t.Fatalf("expressions=%d", len(tl.Expressions))
	}
}
