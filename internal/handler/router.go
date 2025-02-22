package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sonikq/gophermart/internal/config"
	"github.com/sonikq/gophermart/internal/handler/gophermart"
	"github.com/sonikq/gophermart/internal/middleware"
	"github.com/sonikq/gophermart/internal/service"
	"github.com/sonikq/gophermart/pkg/logger"
	"net/http"
)

type Handler struct {
	Gophermart *gophermart.Handler
}

type Option struct {
	Conf    config.Config
	Logger  *logger.Logger
	Service *service.Service
}

func NewRouter(option Option) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CompressResponse(option.Logger), middleware.DecompressRequest(option.Logger))
	router.Use(middleware.RequestResponseLogger(option.Logger))

	h := Handler{Gophermart: gophermart.New(&gophermart.HandlerConfig{
		Config:  option.Conf,
		Logger:  option.Logger,
		Service: option.Service,
	}),
	}

	router.GET("/ping_gophermart_service", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Pong!",
		})
	})

	// Unauthorized routes
	router.POST("/api/user/register", h.Gophermart.Register)
	router.POST("/api/user/login", h.Gophermart.Login)

	// Authorized routes
	authorized := router.Group("/api/user")
	authorized.Use(middleware.IsAuthorized(option.Logger, option.Conf.TokenSecretKey))
	{
		authorized.POST("/orders", h.Gophermart.UploadOrder)
		authorized.GET("/orders", h.Gophermart.ListOrders)
		authorized.GET("/balance", h.Gophermart.GetBalance)
		authorized.POST("/balance/withdraw", h.Gophermart.Withdraw)
		authorized.GET("/withdrawals", h.Gophermart.GetWithdrawals)
	}

	return router
}
