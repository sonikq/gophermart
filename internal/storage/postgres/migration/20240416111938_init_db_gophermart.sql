-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
                                     username TEXT PRIMARY KEY,
                                     password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
                                      number TEXT PRIMARY KEY UNIQUE,
                                      username TEXT NOT NULL,
                                      status TEXT NOT NULL,
                                      accrual DOUBLE PRECISION,
                                      uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL,
                                      updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

ALTER TABLE orders ADD CONSTRAINT fk_username FOREIGN KEY (username) REFERENCES users(username);

CREATE TABLE IF NOT EXISTS balances (
                                        id SERIAL PRIMARY KEY,
                                        order_num TEXT NOT NULL,
                                        username TEXT NOT NULL UNIQUE,
                                        current_balance DOUBLE PRECISION DEFAULT 0,
                                        withdrawal_sum DOUBLE PRECISION DEFAULT 0,
                                        withdraw_processed_at TIMESTAMP WITH TIME ZONE NOT NULL
);

ALTER TABLE balances ADD CONSTRAINT fk_username FOREIGN KEY (username) REFERENCES users(username);
ALTER TABLE balances ADD CONSTRAINT fk_order_num FOREIGN KEY (order_num) REFERENCES orders(number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS balances;
-- +goose StatementEnd