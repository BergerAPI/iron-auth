import express, { Request, Response } from "express";
import { findClientById, findUserById } from "./database";
import { generateAccessToken, generateAuthCode, validateToken } from "./tokens";

const port: number = 3000;

(async () => {
	// Create an Express application
	const app = express();

	// Authorization Endpoint
	app.get('/oauth/authorize', (req, res): void => {
		const { client_id, redirect_uri, state } = req.query;

		if (!client_id || !redirect_uri)
			return res.status(400).send({ error: 'Query Parameters are missing. Please provide client_id and redirect_uri' }) as never

		const client = findClientById(client_id as string);

		if (client === null || client.redirect_uri !== (redirect_uri as string))
			return res.status(400).json({ error: 'Invalid client or redirect_uri' }) as never;

		// TODO: validate user

		// Generate authorization code
		const authCode = generateAuthCode(client.id, "user.id");

		// Redirect back with code and state
		const redirectWithCode = `${client.redirect_uri}?code=${authCode}&state=${state}`;

		return res.redirect(redirectWithCode);
	});

	// Token Exchange Endpoint
	app.post('/oauth/token', express.json(), (req: Request, res: Response) => {
		const { grant_type, code, client_id, client_secret } = req.body;

		if (grant_type !== 'authorization_code') {
			return res.status(400).json({ error: 'Unsupported grant_type' }) as never;
		}

		const client = findClientById(client_id);

		if (!client || client_secret !== client.secret) {
			return res.status(400).json({ error: 'Invalid client credentials' }) as never;
		}

		// Here, we skip validating the code (should be done with a DB lookup)
		try {
			const decoded = validateToken(code);

			if (decoded === null || decoded.clientId !== client_id)
				return res.status(400).json({ error: 'Invalid authorization code' }) as never;

			// Generate access token
			const accessToken = generateAccessToken(client_id, decoded.userId);

			return res.json({ access_token: accessToken, token_type: 'Bearer', expires_in: 3600 }) as never;
		} catch (err) {
			return res.status(400).json({ error: 'Invalid authorization code' }) as never;
		}
	});

	// Start the server and listen on the specified port
	app.listen(port, () => {
		// Log a message when the server is successfully running
		console.log(`Server is running on http://localhost:${port}`);
	});
})().then(_ => { /* */ })