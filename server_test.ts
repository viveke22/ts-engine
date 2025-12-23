// server_test.ts - Showcase of implemented HTTP Server Features

// 1. Import Syntax
import * as http from 'http';

// 2. Create Server with Typed Handler
const server: any = http.createServer(function (req: http.IncomingMessage, res: http.ServerResponse) {
    console.log("Received request: " + req.method + " " + req.url);

    // 3. Routing Logic
    if (req.url === '/') {
        // 4. Response Headers and Body
        res.writeHead(200, { 'Content-Type': 'text/html' });
        res.end('<h1>Welcome to TS-Engine Server!</h1><p>Running on native Go.</p>');
    } else if (req.url === '/json') {
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end('{ "status": "ok", "engine": "ts-engine" }');
    } else {
        res.writeHead(404, { 'Content-Type': 'text/plain' });
        res.end('Page Not Found');
    }
});

// 5. Start Server
// Note: This blocks the execution loop.
console.log("Starting server on http://localhost:3000");
server.listen(3000, function () {
    // Callback support (executed before blocking)
    console.log("Server is listening!");
});
