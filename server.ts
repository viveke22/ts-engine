import * as http from 'http';

const server: any = http.createServer(function (req: any, res: any) {
    console.log("Request received:", req.method, req.url);
    if (req.url === '/' && req.method === 'GET') {
        res.writeHead(200, { 'Content-Type': 'text/html' });
        res.end('<h1>Hello from TS Enginef!</h1> ');
    } else {
        res.writeHead(404);
        res.end('Not Found');
    }
});

server.listen(3000, function () {
    console.log('Server running on http://localhost:3000');
});
