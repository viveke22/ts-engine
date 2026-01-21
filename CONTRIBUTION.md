# Contributing to ts-engine

Thank you for your interest in contributing to **ts-engine**! We welcome contributions from the community to help make this TypeScript-to-native execution engine even better.

## 🚀 Getting Started

### Prerequisites
- **Go**: Version 1.25 or higher.
- **Git**: For version control.

### Setup
1.  **Clone the repository**:
    ```bash
    git clone https://github.com/yourusername/ts-engine.git
    cd ts-engine
    ```
2.  **Install dependencies**:
    ```bash
    go mod download
    ```

## 🏗️ Project Structure

The codebase follows a standard interpreter architecture:

- **`token/`**: Defines tokens (keywords, operators, identifiers).
- **`lexer/`**: Tokenizes the input source code.
- **`ast/`**: Defines the Abstract Syntax Tree nodes.
- **`parser/`**: Parses tokens into an AST (Pratt Parser).
- **`evaluator/`**: Walks the AST and executes code.
- **`object/`**: Defines runtime values (Integers, Strings, Functions including Environment/Scope).
- **`main.go`**: Entry point and CLI.

## 🏃 Running and Building

### Running from source
You can run TypeScript files directly using `go run`:
```bash
go run main.go path/to/file.ts
```

### Building the binary
To build a standalone executable:
```bash
# Windows
go build -o tse.exe main.go

# Linux/Mac
go build -o tse main.go
```
Then run it:
```bash
./tse.exe test_loops.ts
```

## 🧪 Testing

### Unit Tests
The parser package has unit tests written in Go. Run them using:
```bash
go test ./...
```

### Integration Tests
We verify features using TypeScript/JavaScript files in the root directory. Key test files include:
- `test.js`: General feature test.
- `test_loops.js`: For loops.
- `test_while.ts`: While loops.
- `loops.ts`: Increment/Decrement operators.
- `server_test.ts`: HTTP server features.

When adding a new feature, please create a corresponding `.ts` or `.js` file to verify it works as expected.

## 🤝 Contribution Workflow

1.  **Fork** the repository.
2.  **Create a branch** for your feature or bug fix:
    ```bash
    git checkout -b feature/my-new-feature
    ```
3.  **Implement** your changes.
    - If adding a language feature, remember to update `token`, `ast`, `parser`, and `evaluator` layers.
    - Add a test case file to verify.
4.  **Run tests** to ensure no regressions.
5.  **Commit** your changes with clear messages.
6.  **Push** to your fork and submit a **Pull Request**.

## ✨ Coding Standards
- Follow standard **Go formatting** (`go fmt`).
- Keep code readable and well-commented where complex.
- Stick to the existing style for AST and Evaluator logic.

We look forward to your contributions!
