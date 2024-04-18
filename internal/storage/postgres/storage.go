package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sonikq/gophermart/internal/models"
	"log"
	"time"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(ctx context.Context, uri string, dbPoolWorkers int) (*Storage, error) {
	t1 := time.Now()
	var pool *pgxpool.Pool
	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, err
	}
	config.MaxConns = int32(dbPoolWorkers)

	pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	if err = migrate(pool); err != nil {
		return nil, err
	}

	log.Printf("connection to database took: %v\n", time.Since(t1))

	return &Storage{pool: pool}, nil

}

func (ps *Storage) Close() {
	ps.pool.Close()
}

// RegisterUser - creates a new user
func (ps *Storage) RegisterUser(ctx context.Context, username string, password string) error {
	_, err := ps.pool.Exec(ctx, registerUserQuery, username, password)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return models.ErrInvalidCredentials
			}
		}
	}
	return err
}

// GetCredentials - receiving a password hash by username
func (ps *Storage) GetCredentials(ctx context.Context, username string) (string, error) {
	var pwdHash string
	if err := ps.pool.QueryRow(ctx, getCredentialsQuery, username).Scan(&pwdHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", models.ErrInvalidCredentials
		}

		return "", err
	}
	return pwdHash, nil
}

// GetOrder - getting order info by order_id
func (ps *Storage) GetOrder(ctx context.Context, orderNumber string) (*models.Order, error) {
	var order models.Order
	if err := ps.pool.QueryRow(ctx, getOrder, orderNumber).
		Scan(&order.Number, &order.Username, &order.Status, &order.Accrual,
			&order.UploadedAt, &order.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("error in getting order from db: %w", err)
	}
	return &order, nil
}

// UploadOrder - uploading order
func (ps *Storage) UploadOrder(ctx context.Context, orderNumber, username string) error {
	now := time.Now()
	_, err := ps.pool.Exec(ctx, uploadOrder, orderNumber, username, models.NewOrder, now, now)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return models.ErrNotUniqueOrderNum
			}
		}
	}
	return err
}

// ListUserOrders - obtaining a list of information about orders by username
func (ps *Storage) ListUserOrders(ctx context.Context, username string) ([]models.Order, error) {
	var orders []models.Order

	rows, err := ps.pool.Query(ctx, listOrdersQuery, username)
	if err != nil {
		return nil, fmt.Errorf("error in executing pool.Query: %w", err)
	}

	for rows.Next() {
		var order models.Order

		err = rows.Scan(&order.Number, &order.Username, &order.Status, &order.Accrual, &order.UploadedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error in reading rows: %w", err)
	}

	return orders, nil

}

// UpdateOrders - updating information on orders by username
func (ps *Storage) UpdateOrders(ctx context.Context, username string, infos []models.AccrualInfo) error {
	tx, err := ps.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error while begin transaction: %w", err)
	}
	defer func() {
		if errRollBack := tx.Rollback(ctx); errRollBack != nil {
			fmt.Printf("rollback error: %v", errRollBack)
		}
	}()

	for _, info := range infos {
		_, err = tx.Exec(ctx, updateOrderQuery, info.Accrual, info.Status, time.Now(), username, info.Order)
		if err != nil {
			return fmt.Errorf("error in executing update order query: %w", err)
		}

		_, err = tx.Exec(ctx, updateBalanceQuery, username, info.Accrual)
		if err != nil {
			return fmt.Errorf("error in executing update balance query: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// GetBalance - Getting the user's current balance
func (ps *Storage) GetBalance(ctx context.Context, username string) (*models.Balance, error) {
	var balance models.Balance
	err := ps.pool.QueryRow(ctx, getUserBalanceQuery, username).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error in executing get balance query: %w", err)
	}
	return nil, nil
}

// GetWithdrawals - Receiving information about withdrawal of funds
func (ps *Storage) GetWithdrawals(ctx context.Context, username string) ([]models.Withdrawal, error) {
	rows, err := ps.pool.Query(ctx, getWithdrawalsQuery, username)
	if err != nil {
		return nil, fmt.Errorf("error in executing pool.Query(): %w", err)
	}

	var withdrawals []models.Withdrawal
	for rows.Next() {
		var withdrawal models.Withdrawal
		if err = rows.Scan(&withdrawal.Order, &withdrawal.Sum, &withdrawal.ProcessedAt); err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error in reading rows: %w", rows.Err())
	}

	return withdrawals, nil
}

// Withdraw - Request for debiting funds
func (ps *Storage) Withdraw(ctx context.Context, username, order string, sum, delta float64) error {
	tx, err := ps.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error while begin transaction: %w", err)
	}
	defer func() {
		if errRollBack := tx.Rollback(ctx); errRollBack != nil {
			fmt.Printf("rollback error: %v", errRollBack)
		}
	}()

	_, err = tx.Exec(ctx, decrementBalanceQuery, username, delta, sum)
	if err != nil {
		return fmt.Errorf("error in decrement balance query: %w", err)
	}

	_, err = tx.Exec(ctx, withdrawnQuery, order, username, sum, time.Now())
	if err != nil {
		return fmt.Errorf("error in withdrawn query: %w", err)
	}

	return tx.Commit(ctx)
}
