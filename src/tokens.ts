import jwt from "jsonwebtoken"

const SECRET_KEY = 'your_secret_key';  // Store this securely in production!

interface TokenPayload {
	clientId: string
	userId: string
}

/**
 * Creates an auth code for a specific user on a specific client
 * @returns the created code
 */
export function generateAuthCode(clientId: string, userId: string): string {
	// In real implementation, store this in DB with expiration
	return jwt.sign({ clientId, userId } as TokenPayload, SECRET_KEY, { expiresIn: '2m' });
}

/**
 * Creates an access token for a specific user on a specific client
 * @returns the created token
 */
export function generateAccessToken(clientId: string, userId: string): string {
	return jwt.sign({ clientId, userId } as TokenPayload, SECRET_KEY, { expiresIn: '30d' });
}

/**
 * Creates an auth token for the authentication page
 * @returns the created token
 */
export function generateAuthToken(userId: string): string {
	return jwt.sign({ userId }, SECRET_KEY, { expiresIn: '30d' });
}

/**
 * Checks whether or not the token was created by an authorized instance
 * @param token to check
 */
export function validateToken(token: string): TokenPayload | null {
	try {
		return jwt.verify(token, SECRET_KEY) as TokenPayload;
	} catch (e) {
		return null;
	}
}