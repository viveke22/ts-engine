package evaluator

import (
	"os"
	"path/filepath"
	"strings"
	"ts-engine/lexer"
	"ts-engine/object"
	"ts-engine/parser"
)

var requireBuiltin = &object.Builtin{}

func init() {
	requireBuiltin.Fn = func(args ...object.Object) object.Object {
		return requireModule(nil, args)
	}
}

var (
	moduleCache   = map[string]*object.Hash{}
	moduleLoading = map[string]*object.Environment{}
)

func requireModule(env *object.Environment, args []object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	name, ok := args[0].(*object.String)
	if !ok {
		return newError("argument to `require` must be STRING, got %s", args[0].Type())
	}

	if name.Value == "http" {
		return &object.Hash{
			Pairs: map[string]object.Object{
				"createServer": &object.Builtin{
					Fn: createHttpServer,
				},
			},
		}
	}

	if name.Value == "fs" {
		return createFsModule()
	}

	fromDir := ""
	if env != nil {
		fromDir = env.Dir
	}
	return loadFileModule(name.Value, fromDir)
}

func loadFileModule(spec, fromDir string) object.Object {
	filename, err := resolveModulePath(spec, fromDir)
	if err != nil {
		return newError("module not found: %s", spec)
	}

	if cached, ok := moduleCache[filename]; ok {
		return cached
	}
	if loading, ok := moduleLoading[filename]; ok {
		return &object.Hash{Pairs: loading.Exports}
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return newError("module not found: %s", spec)
	}

	l := lexer.New(string(content))
	p := parser.New(l, strings.HasSuffix(filename, ".ts"))
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		return newError("parser errors in %s: %s", filename, strings.Join(p.Errors(), "; "))
	}

	newEnv := object.NewEnvironment()
	newEnv.Dir = filepath.Dir(filename)
	moduleLoading[filename] = newEnv
	result := Eval(prog, newEnv)
	delete(moduleLoading, filename)
	if isError(result) {
		return result
	}

	mod := &object.Hash{Pairs: newEnv.Exports}
	moduleCache[filename] = mod
	return mod
}

func resolveModulePath(spec, fromDir string) (string, error) {
	base := spec
	if !filepath.IsAbs(spec) {
		dir := fromDir
		if dir == "" {
			wd, err := os.Getwd()
			if err != nil {
				return "", err
			}
			dir = wd
		}
		base = filepath.Join(dir, spec)
	}

	candidates := []string{base}
	ext := filepath.Ext(base)
	if ext != ".ts" && ext != ".js" {
		candidates = []string{base + ".ts", base + ".js", base}
	}

	for _, c := range candidates {
		info, err := os.Stat(c)
		if err != nil || info.IsDir() {
			continue
		}
		abs, err := filepath.Abs(c)
		if err != nil {
			return "", err
		}
		return abs, nil
	}
	return "", os.ErrNotExist
}
