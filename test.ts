// test.ts - Showcase of implemented TypeScript Language Features
// (Run with strict mode enabled implicitly for .ts files)

// 1. Strict Variable Typing
let count: number = 42;
let message: string = "Hello TypeScript";
let isActive: boolean = true;
let anything: any = "Could be anything";

console.log(message + ", Count: " + count);

// 2. Typed Function Parameters and Return Type
function multiply(a: number, b: number): number {
    return a * b;
}

console.log("Multiply(10, 2): " + multiply(10, 2));

// 3. Complex/Dotted Types (treated as 'any' at runtime but parsed correctly)
// This simulates types like http.IncomingMessage
function handleRequest(req: http.IncomingMessage, res: http.ServerResponse) {
    console.log("Handing request with types...");
}

// 4. Declare Statements (parsed and ignored to satisfy IDE)
declare var console: any;
declare var http: any;

// 5. Type Enforcement verification
// Uncommenting the line below should cause a runtime error in strict mode:
// let failure: number = "this is a string"; 