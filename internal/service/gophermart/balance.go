package gophermart

import (
	"context"
	"github.com/sonikq/gophermart/internal/models"
)

func (s *Service) GetBalance(ctx context.Context, username string) (models.Balance, error) {
	if err := s.UpdateUserOrders(ctx, username); err != nil {
		return models.Balance{}, err
	}
	return s.storage.GetBalance(ctx, username)
}
