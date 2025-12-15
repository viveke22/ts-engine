# Features of ts-engine

## Implemented Features

### 🚀 Direct Execution
Run TypeScript files directly without manual transpilation.
- **Command**: `tse <filename.ts>`

### 📦 Build System
Create standalone, distributable executables from your TypeScript code.
- **Command**: `tse build <filename.ts>`
- **Output**: A native `.exe` file that works without needing `tse` installed.

### 📝 Variables & Types
- **Declarations**: `let`, `const`, `var`.
- **Type Annotations**: Supports TypeScript syntax like `let x: number = 10;` or `function foo(a: string): void`.
  - *Note: Types are currently parsed for syntax validity but not strictly enforced at runtime.*

### 🛠️ Functions
- **Named Functions**: `function add(a, b) { ... }`
- **Function Expressions**: `let add = function(a, b) { ... }`
- ** Closures**: Full support for closures and lexical scoping.

### 🔀 Control Flow
- **Conditionals**: `if`, `else if`, `else`.
- **Operators**: 
  - Arithmetic: `+`, `-`, `*`, `/`, `%`
  - Comparison: `>`, `<`, `==`, `!=`, `===`, `!==`
  - Logical: `!` (prefix)

### 🖥️ Built-ins
- **Console**: `console.log(...)` for printing output.

### 📄 Comments
- Single-line: `// ...`
- Multi-line: `/* ... */`

### 🔧 IDE Support
- Supports `export {}` to treat files as modules, ensuring compatibility with standard TypeScript IDE tooling.

---

## 🔮 Upcoming Features

We are actively working on expanding `ts-engine`. Planned features include:

- **HTTP Support**: Native HTTP server and client capabilities for building web services.
- **File System API**: Read and write files directly.
- **Advanced Types**: Interfaces, generics, and strict type checking.
- **Type Inference**: Automatic type deduction.
- **Imports/Exports**: Real module loading between files.
