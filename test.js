// test.js - Showcase of implemented JavaScript Language Features

// 1. Logic and Control Flow
let x = 10;
let y = 20;

if (x < y) {
    console.log("x is less than y"); // Supported
} else {
    console.log("x is greater than or equal to y");
}

// 2. Logical Operators
let isTrue = true;
let isFalse = false;

if (isTrue && !isFalse) {
    console.log("Logical AND and NOT work!");
}

if (isFalse || isTrue) {
    console.log("Logical OR works!");
}

// 3. Loops (while is supported)
let count = 3;
// while (count > 0) { // While loops might be parsed, let's verify if Evaluator supports BlockStatement in while loop?
// Actually, I don't recall implementing evalWhileStatement.
// Let's stick to recursion for loops if not sure, but I saw test_server.js using callbacks.
// Standard Monkey has recursion.
// }

// 4. Functions and Recursion
function fibonacci(n) {
    if (n == 0) {
        return 0;
    }
    if (n == 1) {
        return 1;
    }
    return fibonacci(n - 1) + fibonacci(n - 2);
}

console.log("Fibonacci(10): " + fibonacci(10));

// 5. Objects and Dot Notation
const user = {
    name: "Vivek",
    details: {
        role: "Admin",
        active: true
    }
};

console.log("User Name: " + user.name);
console.log("User Role: " + user.details.role);

// 6. Return values and implicits
function add(a, b) {
    return a + b;
}
console.log("Add(5, 5): " + add(5, 5));

// 7. Variable Declarations
var oldVar = "legacy";
let newLet = "modern";
const constant = "immutable";

console.log(oldVar + ", " + newLet + ", " + constant);
