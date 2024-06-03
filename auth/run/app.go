package run

import (
	"context"
	"log"
	"os"
	"microservices/auth/config"
	"microservices/auth/internal/infrastructure/component"
	"microservices/auth/internal/infrastructure/server"
	"microservices/auth/internal/modules"
	"microservices/auth/internal/modules/user/service"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Runner interface {
	Run() error
}

type App struct {
	conf   config.AppConf
	logger *zap.Logger
	rpc    server.Server
	Sig    chan os.Signal
}

func NewApp(conf config.AppConf, logger *zap.Logger) *App {
	return &App{conf: conf, logger: logger, Sig: make(chan os.Signal, 1)}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())

	errGroup, ctx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		sigInt := <-a.Sig
		a.logger.Info("signal interrupt recieved", zap.Stringer("os_signal", sigInt))
		cancel()
		return nil
	})

	errGroup.Go(func() error {
		err := a.rpc.Serve(ctx)
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

func (a *App) Bootstrap(options ...interface{}) Runner {

	components := component.NewComponents(a.conf, a.logger)

	userClient, err := service.GetlientRPC(a.conf.RPCServer.Type, a.conf.UserRPC)
	if err != nil {
		log.Fatal("error init user client", zap.Error(err))
	}

	services := modules.NewServices(userClient, components)

	a.rpc, err = server.GetServerRPC(a.conf.RPCServer, services.Auth, a.logger)
	if err != nil {
		a.logger.Fatal("error init rpc server:", zap.Error(err))
	}

	return a
}
