# Features of ts-engine

## Implemented Features

### 🚀 Direct Execution
Run TypeScript files directly without manual transpilation.
- **Command**: `tse <filename.ts>`
- **Execution**: AST-walking interpreter written in Go.

### 📦 Build System
Create standalone, distributable executables from your TypeScript code.
- **Command**: `tse build <filename.ts>`
- **Output**: A native `.exe` file that works without needing `tse` installed.

### 🌐 HTTP Server & Client
Native support for building web servers and making requests.
- **Server**: `http.createServer((req, res) => { ... })`
- **Listen**: `server.listen(port, callback)`
- **Request**: `req.method`, `req.url` (Dotted access)
- **Response**: 
    - `res.writeHead(status, headers)`
    - `res.end(body)`
- **Headers**: Full support for setting response headers (e.g. `{ 'Content-Type': 'text/html' }`).
- **Client**: Global `fetch()` API with `await` support.

### 📦 Modules & Imports
- **Import Syntax**: `import * as http from 'http';` supported.
- **Require**: Legacy `require('http')` supported.
- **File Isolation**: `export {}` makes a file a module.
- **Built-in Modules**: `http` (internal).

### 🔒 Strict Mode & Types
- **Strict Mode**: Implicitly enabled for `.ts` files. Enforces mandatory type annotations.
- **Loose Mode**: `.js` files allow missing types.
- **Supported Types**: `number`, `string`, `boolean`, `any`, `unknown`, `never`.
- **Complex Types**: Dotted types like `http.IncomingMessage` are accepted (treated as `any` at runtime).
- **IDE Support**: `ts-engine.d.ts` provided for full IntelliSense.

### 📝 Objects & Variables
- **Declarations**: `let`, `const`, `var`.
- **Object Literals**: `{ key: "value", nested: { data: 1 } }`.
- **Dot Notation**: `obj.key`, `obj.nested.data` (Read access).
- **Arrays**: Basic array support (via host integration).

### 🛠️ Functions & Control Flow
- **Functions**: Named `function foo() {}` and expressions `let f = function() {}`.
- **Recursion**: Fully supported.
- **Control Flow**: `if`, `else if`, `else`, `while` loops.
- **Operators**: Arithmetic, Logical (`&&`, `||`, `!`), Comparison (`===`, `!==`, etc.).

### 🖥️ Built-ins
- **Console**: `console.log(...)`.
- **Fetch**: `fetch(url)`.

---

## 🔮 Upcoming Features (Roadmap)

We are actively working on expanding `ts-engine`. Planned features include:

- **Arrow Functions**: `() => {}` syntax support.
- **Classes**: `class MyClass {}` support.
- **Property Assignment**: `obj.prop = value` support.
- **Template Literals**: Backtick strings with interpolation.
- **Advanced Array Support**: Array literals `[1, 2]` and array methods.
- **Full Module System**: Relative imports `import { x } from './file'`.
- **File System API**: `fs.readFile`, `fs.writeFile`.
- **Async/Await**: Function syntax support (async function ...).
