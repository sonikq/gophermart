package gophermart

import (
	"context"
	"errors"
	"github.com/sonikq/gophermart/internal/app/accrual"
	"github.com/sonikq/gophermart/internal/config"
	"github.com/sonikq/gophermart/internal/handler"
	"github.com/sonikq/gophermart/internal/service"
	"github.com/sonikq/gophermart/internal/storage"
	"github.com/sonikq/gophermart/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Run() {
	conf, err := config.Load("internal/config/.env")
	if err != nil {
		log.Fatal("failed to initialize config")
	}

	lg := logger.NewLogger(conf.ServiceName)

	ctx, cancel := context.WithTimeout(context.Background(), conf.CtxTimeOut)
	defer cancel()

	store, err := storage.New(ctx, conf)
	if err != nil {
		lg.Fatal().Err(err).Msg("failed to initialize storage")
	}
	defer store.Close()

	accrualClient := accrual.NewClient(conf.AccrualSystemAddress)
	serviceManager := service.New(store, accrualClient)
	router := handler.NewRouter(handler.Option{
		Conf:    conf,
		Logger:  lg,
		Service: serviceManager,
	})

	server := &http.Server{
		Addr:           conf.RunAddress,
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Info().Err(err).Msg("failed to run http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel = context.WithTimeout(context.Background(), conf.CtxTimeOut)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		lg.Error().Err(err).Msg("error in shutting down server")
	}

	lg.Info().Msg("server stopped successfully")
}
