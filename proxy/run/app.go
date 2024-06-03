package run

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"proxy/config"
	"proxy/internal/infrastructure/component"
	"proxy/internal/infrastructure/responder"
	"proxy/internal/infrastructure/router"
	"proxy/internal/infrastructure/server"
	"proxy/internal/modules"
	aservice "proxy/internal/modules/auth/service"
	geoservice "proxy/internal/modules/geo/service"
	userservice "proxy/internal/modules/user/service"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Runner interface {
	Run() error
}

type App struct {
	conf   config.AppConf
	logger *zap.Logger
	srv    server.Server
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
		err := a.srv.Serve(ctx)
		if err != nil && err != http.ErrServerClosed {
			a.logger.Info("app: server error:", zap.Error(err))
			return err
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return err
	}

	return nil
}

func (a *App) Boostrap(options ...interface{}) Runner {
	protocol := a.conf.RPCServer.Type

	responseManager := responder.NewResponder(a.logger)

	components := component.NewComponents(a.conf, responseManager, a.logger)

	geoClientRPC, err := geoservice.GetlientRPC(protocol, a.conf.GeoRPC)
	if err != nil {
		a.logger.Fatal("error init geo client", zap.Error(err))
	}

	userClientRPC, err := userservice.GetlientRPC(protocol, a.conf.UserRPC)
	if err != nil {
		a.logger.Fatal("error init user client", zap.Error(err))
	}

	authClientRPC, err := aservice.GetlientRPC(protocol, a.conf.AuthRPC)
	if err != nil {
		a.logger.Fatal("error init auth client", zap.Error(err))
	}

	controllers := modules.NewControllers(authClientRPC, geoClientRPC, userClientRPC, components)

	r := router.NewRouter(controllers)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.conf.Server.Port),
		Handler: r,
	}

	a.srv = server.NewHTTPServer(a.conf.Server, srv, a.logger)

	return a
}
