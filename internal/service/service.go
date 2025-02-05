package service

import (
	"context"
	"github.com/sonikq/gophermart/internal/app/accrual"
	"github.com/sonikq/gophermart/internal/models"
	"github.com/sonikq/gophermart/internal/service/gophermart"
	"github.com/sonikq/gophermart/internal/storage"
)

type IGophermartService interface {
	Register(ctx context.Context, request models.RegisterUserRequest) error
	Login(ctx context.Context, request models.LoginRequest) error
	UploadOrder(ctx context.Context, orderNum string, username string) error
	ListUserOrders(ctx context.Context, username string) ([]models.Order, error)
	GetBalance(ctx context.Context, username string) (models.Balance, error)
	Withdraw(ctx context.Context, request models.WithdrawRequest) error
	GetWithdrawals(ctx context.Context, username string) ([]models.Withdrawal, error)
}

type Service struct {
	IGophermartService
}

func New(store storage.IStorage, accrualClient *accrual.Client) *Service {
	return &Service{
		IGophermartService: gophermart.NewService(store, accrualClient),
	}
}
