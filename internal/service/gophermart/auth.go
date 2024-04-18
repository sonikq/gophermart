package gophermart

import (
	"context"
	"github.com/sonikq/gophermart/internal/models"
	"github.com/sonikq/gophermart/pkg/hash"
)

func (s *Service) Register(ctx context.Context, request models.RegisterUserRequest) error {
	pwdHash, err := hash.Generate(request.Password)
	if err != nil {
		return err
	}

	return s.storage.RegisterUser(ctx, request.Login, pwdHash)
}

func (s *Service) Login(ctx context.Context, request models.LoginRequest) error {
	pwdHash, err := s.storage.GetCredentials(ctx, request.Login)
	if err != nil {
		return err
	}

	return hash.Compare(pwdHash, request.Password)
}
