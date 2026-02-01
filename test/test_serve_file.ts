import * as http from 'http';
import * as fs from 'fs';

const server = http.createServer(function (req: any, res: any) {
    console.log("Request received:", req.method, req.url);
    if (req.url === '/' && req.method === 'GET') {
        if (fs.existsSync("example.html")) {
            const content = fs.readFileSync("example.html");
            res.writeHead(200, { 'Content-Type': 'text/html' });
            res.end(content);
        } else {
            res.writeHead(404);
            res.end("example.html not found");
        }
    } else {
        res.writeHead(404);
        res.end('Not Found');
    }
});

server.listen(3000, function () {
    console.log('Server running on http://localhost:3000');
});
