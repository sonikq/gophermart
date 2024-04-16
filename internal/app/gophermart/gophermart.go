package gophermart

import (
	"context"
	"errors"
	"github.com/sonikq/gophermart/internal/app/accrual"
	"github.com/sonikq/gophermart/internal/config"
	"github.com/sonikq/gophermart/internal/handler"
	httpserv "github.com/sonikq/gophermart/internal/server/http"
	"github.com/sonikq/gophermart/internal/service"
	"github.com/sonikq/gophermart/internal/storage"
	"github.com/sonikq/gophermart/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
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

	workerPool := make(chan accrual.WorkerPool)
	defer close(workerPool)

	accrualClient := accrual.NewClient(conf.AccrualSystemAddress, workerPool)
	go accrualClient.Run()

	serviceManager := service.New(store, accrualClient)
	router := handler.NewRouter(handler.Option{
		Conf:    conf,
		Logger:  lg,
		Service: serviceManager,
	})

	server := httpserv.NewServer(conf.RunAddress, router)

	go func() {
		err = server.Run()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Info().Err(err).Msg("failed to run http server")
		}
	}()

	lg.Info().Msg("Server listening on " + conf.RunAddress)

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
