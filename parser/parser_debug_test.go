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
