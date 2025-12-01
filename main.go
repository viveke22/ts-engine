package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"ts-engine/evaluator"
	"ts-engine/lexer"
	"ts-engine/object"
	"ts-engine/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ts-engine <filename.ts>")
		return
	}

	filename := os.Args[1]
	code, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	env := object.NewEnvironment()
	l := lexer.New(string(code))

	// Debug: print tokens
	// for {
	// 	tok := l.NextToken()
	// 	fmt.Printf("%+v\n", tok)
	// 	if tok.Type == token.EOF {
	// 		break
	// 	}
	// }
	// l = lexer.New(string(code)) // Reset lexer

	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
		// Don't print the result of the last expression unless it's an error,
		// as standard TS/JS runtimes don't output the last value like a REPL.
		// But for debugging/verification it might be useful.
		// The user asked for "console.log" so we should rely on that for output.
		if evaluated.Type() == object.ERROR_OBJ {
			fmt.Println(evaluated.Inspect())
		}
	}
}

func printParserErrors(errors []string) {
	fmt.Println("Parser errors:")
	for _, msg := range errors {
		fmt.Printf("\t%s\n", msg)
	}
}
