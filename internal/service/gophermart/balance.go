package gophermart

import (
	"context"
	"fmt"
	"github.com/sonikq/gophermart/internal/models"
)

func (s *Service) GetBalance(ctx context.Context, username string) (models.Balance, error) {
	if err := s.UpdateUserOrders(ctx, username); err != nil {
		return models.Balance{}, err
	}
	return s.storage.GetBalance(ctx, username)
}

func (s *Service) Withdraw(ctx context.Context, request models.WithdrawRequest) error {
	if err := validateOrderNum(request.Order); err != nil {
		return fmt.Errorf("%w: %w", models.ErrInvalidOrderNum, err)
	}

	if err := s.UpdateUserOrders(ctx, request.Username); err != nil {
		return err
	}

	balance, err := s.GetBalance(ctx, request.Username)
	if err != nil {
		return err
	}

	if request.Sum > balance.Current {
		return models.ErrInsufficientFunds
	}

	return s.storage.Withdraw(ctx,
		request.Username,
		request.Order,
		request.Sum,
		balance.Current-request.Sum,
	)
}

func (s *Service) GetWithdrawals(ctx context.Context, username string) ([]models.Withdrawal, error) {
	return s.storage.GetWithdrawals(ctx, username)
}
