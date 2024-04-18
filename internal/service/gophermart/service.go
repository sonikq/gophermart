package gophermart

import (
	"github.com/sonikq/gophermart/internal/app/accrual"
	"github.com/sonikq/gophermart/internal/storage"
)

type Service struct {
	storage       storage.IStorage
	accrualClient *accrual.Client
}

func NewService(store storage.IStorage, accrualClient *accrual.Client) *Service {
	return &Service{
		storage:       store,
		accrualClient: accrualClient,
	}
}
