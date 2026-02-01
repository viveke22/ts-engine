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
const server = http.createServer(function (req, res) {
    console.log("Request received:", req.method, req.url);
    if (req.url === '/' && req.method === 'GET') {
        res.writeHead(200, { 'Content-Type': 'text/html' });
        res.end('<h1>Hello from TS Engine!</h1> ');
    } else {
        res.writeHead(404);
        res.end('Not Found');
    }
});

server.listen(3000, function () {
    console.log('Server running on http://localhost:3000');
});