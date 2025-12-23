package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"ts-engine/evaluator"
	"ts-engine/lexer"
	"ts-engine/object"
	"ts-engine/parser"
)

const magicHeaderStart = "#####"
const magicHeaderEnd = "TSE_DATA#####"

func main() {
	// Construct magic string dynamically so it doesn't appear as a literal in the binary
	magicMarker := []byte(magicHeaderStart + magicHeaderEnd)

	// 1. Check if we are running as a bundled executable
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return
	}

	exeBytes, err := ioutil.ReadFile(exePath)
	if err != nil {
		fmt.Println("Error reading executable:", err)
		return
	}

	if bytes.Contains(exeBytes, magicMarker) {
		parts := bytes.Split(exeBytes, magicMarker)
		// The last part is the source code
		if len(parts) > 1 {
			sourceCode := string(parts[len(parts)-1])
			runCode(sourceCode, true) // Embedded code is assumed to be TS/Strict
			return
		}
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: ts-engine <filename.ts> OR ts-engine build <filename.ts>")
		return
	}

	command := os.Args[1]

	if command == "build" {
		if len(os.Args) < 3 {
			fmt.Println("Usage: ts-engine build <filename.ts>")
			return
		}
		sourceFile := os.Args[2]
		buildExecutable(exePath, sourceFile, magicMarker)
		return
	}

	// Normal execution
	filename := os.Args[1]
	code, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}
	isStrict := strings.HasSuffix(filename, ".ts")
	runCode(string(code), isStrict)
}

func runCode(code string, isStrict bool) {
	env := object.NewEnvironment()
	l := lexer.New(code)
	p := parser.New(l, isStrict)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
		if evaluated.Type() == object.ERROR_OBJ {
			fmt.Println(evaluated.Inspect())
		}
	}
}

func buildExecutable(selfPath, sourcePath string, magicMarker []byte) {
	// Read source
	sourceCode, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		fmt.Printf("Error reading source file: %s\n", err)
		return
	}

	// Read self
	selfBytes, err := ioutil.ReadFile(selfPath)
	if err != nil {
		fmt.Printf("Error reading executable: %s\n", err)
		return
	}

	// Output filename
	baseName := strings.TrimSuffix(filepath.Base(sourcePath), filepath.Ext(sourcePath))
	outName := baseName + ".exe"

	// Create new file
	outFile, err := os.Create(outName)
	if err != nil {
		fmt.Printf("Error creating output file: %s\n", err)
		return
	}
	defer outFile.Close()

	// Write self + magic + source
	if _, err := outFile.Write(selfBytes); err != nil {
		fmt.Printf("Error writing executable bytes: %s\n", err)
		return
	}

	if _, err := outFile.Write(magicMarker); err != nil {
		fmt.Printf("Error writing magic marker: %s\n", err)
		return
	}

	if _, err := outFile.Write(sourceCode); err != nil {
		fmt.Printf("Error writing source code: %s\n", err)
		return
	}

	fmt.Printf("Built %s successfully.\n", outName)
}

func printParserErrors(errors []string) {
	fmt.Println("Parser errors:")
	for _, msg := range errors {
		fmt.Printf("\t%s\n", msg)
	}
}
