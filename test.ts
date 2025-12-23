// test.ts - Showcase of implemented JavaScript Language Features with TypeScript Types

// 1. Logic and Control Flow
let x: number = 10;
let y: number = 20;

if (x < y) {
    console.log("x is less than y");
} else {
    console.log("x is greater than or equal to y");
}

// 2. Logical Operators
let isTrue: boolean = true;
let isFalse: boolean = false;

if (isTrue && !isFalse) {
    console.log("Logical AND and NOT work!");
}

if (isFalse || isTrue) {
    console.log("Logical OR works!");
}

// 3. Functions and Recursion
function fibonacci(n: number): number {
    if (n == 0) {
        return 0;
    }
    if (n == 1) {
        return 1;
    }
    return fibonacci(n - 1) + fibonacci(n - 2);
}

console.log("Fibonacci(10): " + fibonacci(10));

// 4. Objects and Dot Notation
// Note: Object literals need 'any' type for now as we don't support interface definitions in parser yet
const user: any = {
    name: "Vivek",
    details: {
        role: "Admin",
        active: true
    }
};

console.log("User Name: " + user.name);
console.log("User Role: " + user.details.role);

// 5. Return values and implicits
function add(a: number, b: number): number {
    return a + b;
}
console.log("Add(5, 5): " + add(5, 5));

// 6. Variable Declarations
var oldVar: string = "legacy";
let newLet: string = "modern";
const constant: string = "immutable";

console.log(oldVar + ", " + newLet + ", " + constant);
export {}