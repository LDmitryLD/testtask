package run

import (
	"context"
	"os"

	"github.com/LDmitryLD/testtask/config"
	"github.com/LDmitryLD/testtask/grpc/rates"
	"github.com/LDmitryLD/testtask/internal/db"
	"github.com/LDmitryLD/testtask/internal/infrastructure/server"
	"github.com/LDmitryLD/testtask/internal/modules"
	"github.com/LDmitryLD/testtask/internal/storages"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Runner interface {
	Run() error
}

type App struct {
	conf   config.AppConf
	logger *zap.Logger
	grpc   server.Server
	Sig    chan os.Signal
}

func NewApp(conf config.AppConf, logger *zap.Logger) *App {
	return &App{
		conf:   conf,
		logger: logger,
		Sig:    make(chan os.Signal, 1),
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())

	errGroup, ctx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		sigInt := <-a.Sig
		a.logger.Info("signal interrupt recieved", zap.String("os_signal", sigInt.String()))
		cancel()
		return nil
	})

	errGroup.Go(func() error {
		err := a.grpc.Serve(ctx)
		if err != nil {
			a.logger.Error("app: server error", zap.Error(err))
			return err
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return err
	}

	return nil
}

func (a *App) Bootstrap() Runner {
	_, sqlAdapter, err := db.NewSqlDB(a.conf.DB, a.logger)
	if err != nil {
		a.logger.Fatal("error init db:", zap.Error(err))
	}

	newStorages := storages.NewStorages(sqlAdapter)

	services := modules.NewServices(newStorages, a.logger)

	ratesServiceGRPC := rates.NewRatesService(services.Rates)

	a.grpc = server.NewGRPCServer(a.conf.GRPCServer, ratesServiceGRPC, a.logger)

	return a
}
