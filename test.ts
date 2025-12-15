// 1. Comments verification
// This is a single line comment
/* This is a 
   multi-line comment */

// 2. Variables (let, const, var) and Types
let x: number = 42;
const pi: number = 3;
var message: string = "Hello TS Engine";

console.log(x);
console.log(pi);
console.log(message);

// 3. Control Flow (if, else if, else) and Comparisons
if (x > 50) {
    console.log("x is large");
} else if (x > 40) {
    console.log("x is greater than 40");
} else {
    console.log("x is small");
}

if (x === 42) {
    console.log("Strict equality works");
}

// 4. Functions (Declarations, Expressions, Type Annotations)
function add(a: number, b: number): number {
    return a + b;
}

let sum = add(10, 20);
console.log(sum);

let multiply = function (a: number, b: number) {
    return a * b;
};

console.log(multiply(5, 5));

// 5. Final check
console.log("Test suite completed successfully");

export { };