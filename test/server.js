const express = require("express");
const app = express();

let statusCode = 200; 

app.get("/test", (req, res) => {
    res.status(statusCode).send(`Response with status ${statusCode}`);
});

pp.post("/set-status", (req, res) => {
    const { code } = req.query;
    statusCode = parseInt(code) || 500;
    res.send(`Updated status code to ${statusCode}`);
});

const PORT = 3000;
app.listen(PORT, () => {
    console.log(`Mock server running on http://localhost:${PORT}`);
});
