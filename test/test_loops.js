
console.log("Testing For Loop...");

let sum = 0;
for (let i = 0; i < 5; i = i + 1) {
    console.log("i =", i);
    sum = sum + i;
}
console.log("Sum:", sum);

if (sum == 10) {
    console.log("Loop Test Passed!");
} else {
    console.log("Loop Test Failed! Expected 10, got", sum);
}

// Test with external variable
let j = 0;
for (; j < 3; j = j + 1) {
    console.log("j =", j);
}
console.log("Final j:", j); // Should be 3
