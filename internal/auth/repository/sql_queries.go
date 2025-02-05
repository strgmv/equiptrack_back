package repository

const (
	createUserQuery = `INSERT INTO users (login, password, role) VALUES ($1, $2, $3) RETURNING user_id`
	deleteUserQuery = `DELETE FROM users WHERE user_id = $1`
	getUserQuery    = `SELECT user_id, login, password, role FROM users WHERE user_id = $1`
	findUserByLogin = `SELECT user_id, login, password, role FROM users WHERE login = $1`

	setUserSession       = `INSERT INTO sessions (user_id, refresh_token) VALUES ($1, $2)`
	getUserSession       = `SELECT id, user_id, refresh_token FROM sessions WHERE user_id = $1 AND refresh_token = $2`
	deleteUserSession    = `DELETE FROM sessions WHERE id = $1`
	deleteSessionByToken = `DELETE FROM sessions WHERE user_id = $1 AND refresh_token = $2`

	qGetTotal = `SELECT COUNT(user_id) FROM users`
	qGetUsers = `SELECT user_id, login, role
			FROM users
			OFFSET $1 
			LIMIT $2`
)
