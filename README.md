# ts-engine

A lightweight, high-performance TypeScript runtime written in Go. `ts-engine` compiles and executes TypeScript code directly, offering a seamless development experience.

## Unqiue Selling Point
Unlike Node.js or Deno which rely on V8, `ts-engine` is a custom implementation built from scratch in Go.

## 🚀 Usage

### Running Code
To execute a TypeScript file immediately:

```bash
tse <filename.ts>
```

Example:
```bash
tse test.ts
```

### Building Executables
Need to distribute your app? `ts-engine` can bundle your TypeScript code into a single, standalone executable.

```bash
tse build <filename.ts>
```

This will generate a `<filename>.exe` that runs on any compatible machine, even if they don't have `ts-engine` installed!

## 📥 Installation
(Instructions to download `tse.exe` would go here - e.g. "Download the latest release from the Releases page")

## ✨ Features
We support variables (`const`, `let`, `var`), functions, type annotations, control flow (`if/else`), and more. 

👉 **[See features.md](features.md) for a full list of implemented and upcoming features.**
