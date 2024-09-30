import Database from "better-sqlite3";

// Initializes the connection to the sqlite database
export const database = new Database("./db.sqlite")

interface ClientModel {
	id: string
	name: string
	secret: string
	created_at: Date
	redirect_uri: string
}

interface UserModel {
	id: string
	email: string
	password: string
	created_at: Date
}

/**
 * Selects a single client by their id
 * @param clientId the id to query for
 * @returns the client as {ClientModel}
 */
export const findClientById = (clientId: string): ClientModel | null => {
	const row = database.prepare('SELECT * FROM clients WHERE id = ?').get(clientId) as ClientModel | null;

	if (!row)
		return null

	return { ...row, created_at: new Date(row.created_at) }
}

/**
 * Selects a single client by their id
 * @param clientId the id to query for
 * @returns the user as {UserModel}
 */
export const findUserById = (userId: string): UserModel | null => {
	const row = database.prepare('SELECT * FROM users WHERE id = ?').get(userId) as UserModel | null;

	if (!row)
		return null

	return { ...row, created_at: new Date(row.created_at) }
}