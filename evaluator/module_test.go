package evaluator

import (
	"os"
	"path/filepath"
	"testing"
	"ts-engine/lexer"
	"ts-engine/object"
	"ts-engine/parser"
)

func TestRelativeNamedImport(t *testing.T) {
	resetModules()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "math.ts"), `
		export function add(a: number, b: number) {
			return a + b;
		}
		export const pi = 3;
	`)
	entry := filepath.Join(root, "main.ts")
	writeFile(t, entry, `
		import { add, pi } from "./math";
		add(2, pi);
	`)

	intEquals(t, evalFile(t, entry), 5)
}

func TestRelativeStarImportAndRequire(t *testing.T) {
	resetModules()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "lib.ts"), `
		export function n() {
			return 4;
		}
	`)
	star := filepath.Join(root, "star.ts")
	writeFile(t, star, `
		import * as lib from "./lib";
		lib.n();
	`)
	req := filepath.Join(root, "req.ts")
	writeFile(t, req, `
		const lib = require("./lib");
		lib.n();
	`)

	intEquals(t, evalFile(t, star), 4)
	intEquals(t, evalFile(t, req), 4)
}

func TestNestedRelativeImportIgnoresCwd(t *testing.T) {
	resetModules()
	root := t.TempDir()
	pkg := filepath.Join(root, "pkg")
	if err := os.Mkdir(pkg, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(pkg, "c.ts"), `
		export function c() {
			return 3;
		}
	`)
	writeFile(t, filepath.Join(pkg, "b.ts"), `
		import { c } from "./c";
		export function b() {
			return c() + 1;
		}
	`)
	entry := filepath.Join(pkg, "main.ts")
	writeFile(t, entry, `
		import { b } from "./b";
		b();
	`)

	t.Chdir(os.TempDir())

	intEquals(t, evalFile(t, entry), 4)
}

func TestParentRelativeImport(t *testing.T) {
	resetModules()
	root := t.TempDir()
	pkg := filepath.Join(root, "pkg")
	if err := os.Mkdir(pkg, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(root, "util.ts"), `
		export const x = 7;
	`)
	entry := filepath.Join(pkg, "main.ts")
	writeFile(t, entry, `
		import { x } from "../util";
		x;
	`)

	intEquals(t, evalFile(t, entry), 7)
}

func TestJsExtensionAndTsWithoutExtension(t *testing.T) {
	resetModules()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a.js"), `
		export function v() {
			return 9;
		}
	`)
	entry := filepath.Join(root, "main.ts")
	writeFile(t, entry, `
		import { v } from "./a.js";
		v();
	`)

	intEquals(t, evalFile(t, entry), 9)
}

func TestModuleCacheSharesExports(t *testing.T) {
	resetModules()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "counter.ts"), `
		export let n = 0;
		export function inc() {
			n = n + 1;
			return n;
		}
	`)
	entry := filepath.Join(root, "main.ts")
	writeFile(t, entry, `
		import { inc } from "./counter";
		import { once } from "./more";
		inc();
		once();
	`)
	writeFile(t, filepath.Join(root, "more.ts"), `
		import { inc } from "./counter";
		export function once() {
			return inc();
		}
	`)

	intEquals(t, evalFile(t, entry), 2)
}

func TestMissingModuleAndExport(t *testing.T) {
	resetModules()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "lib.ts"), `
		export const a = 1;
	`)
	missing := filepath.Join(root, "missing.ts")
	writeFile(t, missing, `
		import { a } from "./nope";
		a;
	`)
	badExport := filepath.Join(root, "bad.ts")
	writeFile(t, badExport, `
		import { z } from "./lib";
		z;
	`)

	isErrorContaining(t, evalFile(t, missing), "module not found")
	isErrorContaining(t, evalFile(t, badExport), "export 'z' not found")
}

func TestModuleRuntimeError(t *testing.T) {
	resetModules()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "boom.ts"), `
		export const x = 1;
		foo;
	`)
	entry := filepath.Join(root, "main.ts")
	writeFile(t, entry, `
		import { x } from "./boom";
		x;
	`)

	isErrorContaining(t, evalFile(t, entry), "identifier not found")
}

func TestRequireFromExportedFunction(t *testing.T) {
	resetModules()
	root := t.TempDir()
	pkg := filepath.Join(root, "pkg")
	if err := os.Mkdir(pkg, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(pkg, "dep.ts"), `
		export const k = 11;
	`)
	writeFile(t, filepath.Join(pkg, "load.ts"), `
		export function load() {
			const m = require("./dep");
			return m.k;
		}
	`)
	entry := filepath.Join(root, "main.ts")
	writeFile(t, entry, `
		import { load } from "./pkg/load";
		load();
	`)

	intEquals(t, evalFile(t, entry), 11)
}

func TestCircularRelativeImports(t *testing.T) {
	resetModules()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a.ts"), `
		const b = require("./b");
		export const aVal = 1;
		export function getBVal() {
			return b.bVal;
		}
	`)
	writeFile(t, filepath.Join(root, "b.ts"), `
		const a = require("./a");
		export const bVal = 2;
		export function getAVal() {
			return a.aVal;
		}
	`)
	entry := filepath.Join(root, "main.ts")
	writeFile(t, entry, `
		import { getBVal } from "./a";
		import { getAVal } from "./b";
		getBVal() + getAVal();
	`)

	intEquals(t, evalFile(t, entry), 3)
}

func TestModuleParserError(t *testing.T) {
	resetModules()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "syntax.ts"), `
		export const = ;
	`)
	entry := filepath.Join(root, "main.ts")
	writeFile(t, entry, `
		import { x } from "./syntax";
		x;
	`)

	isErrorContaining(t, evalFile(t, entry), "parser errors in")
}

func TestExportClass(t *testing.T) {
	resetModules()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "point.ts"), `
		export class Point {
			constructor(x: number) {
				this.x = x;
			}
			n() {
				return this.x;
			}
		}
	`)
	entry := filepath.Join(root, "main.ts")
	writeFile(t, entry, `
		import { Point } from "./point";
		new Point(8).n();
	`)

	intEquals(t, evalFile(t, entry), 8)
}

func evalFile(t *testing.T, path string) object.Object {
	t.Helper()
	code, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	env := object.NewEnvironment()
	env.Dir = filepath.Dir(abs)
	l := lexer.New(string(code))
	p := parser.New(l, filepath.Ext(path) == ".ts")
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}
	return Eval(program, env)
}

func writeFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func resetModules() {
	moduleCache = map[string]*object.Hash{}
	moduleLoading = map[string]*object.Environment{}
}
