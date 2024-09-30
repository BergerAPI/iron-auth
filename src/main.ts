import express, { Request, Response } from "express";

const port: number = 3000;

(async () => {
	// Create an Express application
	const app = express();

	// Define a route for the root path ('/')
	app.get('/', (req: Request, res: Response) => {
		res.send('Hello, TypeScript + Node.js + Express!');
	});

	// Start the server and listen on the specified port
	app.listen(port, () => {
		// Log a message when the server is successfully running
		console.log(`Server is running on http://localhost:${port}`);
	});
})().then(_ => { /* */ })