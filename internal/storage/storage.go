package storage

import (
	"context"
	"github.com/sonikq/gophermart/internal/config"
	"github.com/sonikq/gophermart/internal/models"
	"github.com/sonikq/gophermart/internal/storage/postgres"
)

type IStorage interface {
	RegisterUser(ctx context.Context, username string, password string) error
	GetCredentials(ctx context.Context, username string) (string, error)
	GetOrder(ctx context.Context, orderNumber string) (*models.Order, error)
	UploadOrder(ctx context.Context, orderNumber, username string) error
	ListUserOrders(ctx context.Context, username string) ([]models.Order, error)
	UpdateOrders(ctx context.Context, username string, infos []models.AccrualInfo) error
	GetBalance(ctx context.Context, username string) (*models.Balance, error)
	GetWithdrawals(ctx context.Context, username string) ([]models.Withdrawal, error)
	Withdraw(ctx context.Context, username, order string, sum, delta float64) error
	Close()
}

func New(ctx context.Context, cfg config.Config) (IStorage, error) {
	return postgres.NewStorage(ctx, cfg.DatabaseURI, cfg.DBPoolWorkers)
}
