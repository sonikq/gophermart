package postgres

const (
	registerUserQuery   = `INSERT INTO users(username, password) VALUES ($1, $2);`
	getCredentialsQuery = `SELECT password from users WHERE username = $1;`
	getOrder            = `SELECT number, username, status, accrual, uploaded_at, updated_at FROM orders WHERE number = $1;`
	uploadOrder         = `INSERT INTO orders(number, username, status, uploaded_at, updated_at) VALUES ($1, $2, $3, $4, $5);`
	listOrdersQuery     = `SELECT number, username, status, accrual, uploaded_at, updated_at FROM orders WHERE username = $1 ORDER BY uploaded_at ASC;`
	updateOrderQuery    = `UPDATE orders SET accrual = $1, status = $2, updated_at = $3 WHERE username = $4 AND number = $5;`
	updateBalanceQuery  = `INSERT INTO balances(username, current_balance) VALUES ($1, $2) ON CONFLICT (username) DO UPDATE SET current_balance = balances.current_balance + $2;`
	getUserBalanceQuery = `SELECT current_balance, withdrawn FROM balances WHERE username = $1;`
	getWithdrawalsQuery = `SELECT order_num, withdrawal_sum, withdraw_processed_at FROM balances WHERE username = $1 ORDER BY withdraw_processed_at ASC;`
	withdrawnQuery      = `INSERT INTO balances(order_num, username, current_balance, withdrawal_sum, withdraw_processed_at)
								VALUES ($1, $2, $3, $4, $5)
								ON CONFLICT (username) DO UPDATE
								SET order_num = $1, current_balance = $3, withdrawn = $4, withdraw_processed_at = $5;`
)
