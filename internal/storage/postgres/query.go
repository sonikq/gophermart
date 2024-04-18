package postgres

const (
	registerUserQuery   = `INSERT INTO users(username, password) VALUES ($1, $2);`
	getCredentialsQuery = `SELECT password from users WHERE username = $1;`
	getOrder            = `SELECT number, username, status, accrual, uploaded_at, updated_at FROM orders WHERE number = $1;`
	uploadOrder         = `INSERT INTO orders(number, username, status, uploaded_at, updated_at) VALUES ($1, $2, $3, $4, $5);`
	listOrdersQuery     = `SELECT number, username, status, accrual, uploaded_at, updated_at FROM orders WHERE username = $1 ORDER BY uploaded_at ASC;`
	updateOrderQuery    = `UPDATE orders SET accrual = $1, status = $2, updated_at = $3 WHERE username = $4 AND number = $5;`
	updateBalanceQuery  = `INSERT INTO balances(username, current_balance)
								VALUES ($1, $2)
								ON CONFLICT (username) DO UPDATE
								SET current_balance = balances.current_balance + $2;`
	getUserBalanceQuery = `SELECT current_balance, withdrawn FROM balances WHERE username = $1;`
	getWithdrawalsQuery = `SELECT order_number, withdrawal_sum, processed_at FROM withdrawals WHERE username = $1 ORDER BY processed_at ASC;`
	withdrawnQuery      = `INSERT INTO withdrawals(order_number, username, withdrawal_sum, processed_at)
								VALUES ($1, $2, $3, $4);`
	decrementBalanceQuery = `INSERT INTO balances(username, current_balance, withdrawn)
								VALUES ($1, $2, $3)
								ON CONFLICT (username) DO UPDATE
								SET current_balance = $2, withdrawn = $3;`
)
