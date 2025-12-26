// server_test.ts - Showcase of implemented HTTP Server Features

import * as http from 'http';

// 2. Create Server with Typed Handler
const server: any = http.createServer(function (req: http.IncomingMessage, res: http.ServerResponse) {
    console.log("Received request: " + req.method + " " + req.url);

    // 3. Routing Logic
    if (req.url === '/favicon.ico') {
        res.writeHead(204);  // No Content – tells browser "nothing here"
        res.end();
        return;
    }
    else if (req.url === '/') {
        // 4. Response Headers and Body
        res.writeHead(200, { 'Content-Type': 'text/html' });
        res.end(`
  <h1>Welcome to TS Engine Server!</h1>
  <button onclick="window.location.href = '/home'">Go to /home</button>
`);
    }else if (req.url === '/home') {
        res.writeHead(200, { 'Content-Type': 'text/html' });
        res.end(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>TS Engine Demo</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      text-align: center;
      padding: 50px;
      background-color: #f0f0f0;
    }
    h1 {
      font-size: 3em;
      margin-bottom: 20px;
      transition: color 0.5s ease;
    }
    button {
      padding: 15px 30px;
      font-size: 1.2em;
      background-color: #007bff;
      color: white;
      border: none;
      border-radius: 8px;
      cursor: pointer;
    }
    button:hover {
      background-color: #0056b3;
    }
  </style>
</head>
<body>
  <h1 id="title">Welcome to TS Engine Server!</h1>
<p>Running on native Go. 🚀</p>

<button onclick="toggleColor()">Toggle Title Color</button>

<script>
  const title = document.getElementById('title');
  let isBlue = false;

  function toggleColor() {
    isBlue = !isBlue;
    title.style.color = isBlue ? 'blue' : 'black';

    // Log red or blue (as requested)
    if (isBlue) {
      console.log("blue");
    } else {
      console.log("red");
    }
  }

  // Initial log
  console.log("red");
</script>
</body>
</html>`);
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
