
console.log("Testing Increment/Decrement...");

// Test basic loop with i++
let sum = 0;
for (let i = 0; i < 5; i++) {
    sum = sum + i;
}
console.log("Sum with i++: " + sum); // Expected 10

if (sum == 10) {
    console.log("PASS: Loop with i++");
} else {
    console.log("FAIL: Loop with i++");
}

// Test postfix return value
let x = 5;
let y = x++;
console.log("x (after x++): " + x); // Expected 6
console.log("y (value of x++): " + y); // Expected 5

if (x == 6 && y == 5) {
    console.log("PASS: Postfix value");
} else {
    console.log("FAIL: Postfix value");
}

// Test decrement
let z = 10;
z--;
console.log("z (after z--): " + z); // Expected 9
if (z == 9) {
    console.log("PASS: Decrement");
} else {
    console.log("FAIL: Decrement");
}

// Test loop with i--
let countdown = 3;
for (let i = 3; i > 0; i--) {
    console.log("Countdown: " + i);
    countdown = countdown - 1;
}
if (countdown == 0) {
    console.log("PASS: Loop with i--");
} else {
    console.log("FAIL: Loop with i--");
}
export {}