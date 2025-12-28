// 1. Arrow Functions
const add = (a, b) => a + b;
console.log("Arrow Add (5+3):", add(5, 3));

const square = x => x * x;
console.log("Arrow Square (4):", square(4));

// 2. HTTP Server (using require)
console.log("Testing HTTP Module...");
const http = require("http");
if (http.createServer) {
    console.log("HTTP createServer exists");
} else {
    console.log("HTTP createServer missing");
}

// 3. Fetch API
console.log("Testing Fetch...");
const res = fetch("https://google.com");
console.log("Fetch returned: " + res);
console.log("Arrow.js tests complete");