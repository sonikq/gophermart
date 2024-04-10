package postgres

const (
	registerUserQuery   = `INSERT INTO users(username, password) VALUES ($1, $2);`
	getCredentialsQuery = `SELECT password from users WHERE username = $1;`
	getOrder            = `SELECT number, username, status, accrual, uploaded_at, updated_at FROM orders WHERE id = $1;`
	uploadOrder         = `INSERT INTO orders(number, username, status, uploaded_at, updated_at) VALUES ($1, $2, $3, $4, $5);`
	listOrdersQuery     = `SELECT number, username, status, accrual, uploaded_at, updated_at FROM orders WHERE username = $1 ORDER BY uploaded_at ASC;`
	updateOrderQuery    = `UPDATE orders SET accrual = $1, status = $2, updated_at = $3 WHERE username = $4 AND number = $5;`
)
