package gophermart

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sonikq/gophermart/internal/config"
	"github.com/sonikq/gophermart/internal/service"
	"github.com/sonikq/gophermart/pkg/logger"
)

const (
	contentTypeHeaderKey = "Content-Type"
	contentTypeJSON      = "application/json"
	contentTypeTextPlain = "text/plain"
)

type Handler struct {
	config  config.Config
	logger  *logger.Logger
	service *service.Service
}

type HandlerConfig struct {
	Config  config.Config
	Logger  *logger.Logger
	Service *service.Service
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		config:  cfg.Config,
		logger:  cfg.Logger,
		service: cfg.Service,
	}
}

func getUsername(c *gin.Context) (string, error) {
	username, exists := c.Get("username")
	if !exists {
		return "", errors.New("failed to get username from context")
	}

	usernameValue, ok := username.(string)
	if !ok {
		return "", errors.New("failed to convert username to string")
	}

	return usernameValue, nil
}
