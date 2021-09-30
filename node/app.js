const fetch = require("node-fetch");
const express = require('express');

const app = express();
const port = 80;

app.get('/hello', (req, res) => {
    res.contentType("text/plain");
    res.send('hello from node\n');
})

app.get("/hello/proxy/*", async (req, res) => {
    res.contentType("text/plain");

    const path = req.path.replace(/\/hello\/proxy\//, "");

    if (path.length === 0) {
        res.status(404)
        return res.send("NOT FOUND")
    }

    const parts = path.split(/\//);
    const proxyTarget = parts[0];
    const rest = parts.slice(1).join("/");

    const url = parts.length === 1 ? `http://${proxyTarget}/hello` : `http://${proxyTarget}/hello/proxy/${rest}`;

    try {
        const proxyResponse = await fetch(url)
        const responseBody = await proxyResponse.text();
        res.send(`${responseBody}hello from node\n`);
    } catch (e) {
        res.status(400);
        res.send(`ERROR ERROR THIS DOES NOT COMPUTE : ${e}`);
    }
})

app.listen(port, () => {
    console.log(`Example app listening at http://localhost:${port}`);
})