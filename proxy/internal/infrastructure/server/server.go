package server

import (
	"context"
	"net/http"
	"proxy/config"

	"go.uber.org/zap"
)

type Server interface {
	Serve(ctx context.Context) error
}

type HTTPServer struct {
	conf   config.Server
	srv    *http.Server
	logger *zap.Logger
}

func NewHTTPServer(conf config.Server, server *http.Server, logger *zap.Logger) Server {
	return &HTTPServer{conf: conf, srv: server, logger: logger}
}

func (s *HTTPServer) Serve(ctx context.Context) error {
	var err error

	chErr := make(chan error)

	go func() {
		s.logger.Info("server starter", zap.String("addr", s.srv.Addr))
		if err = s.srv.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Error("http listen and server error:", zap.Error(err))
			chErr <- err
		}
	}()

	select {
	case <-chErr:
		return err
	case <-ctx.Done():
	}

	ctxShutdown, cancel := context.WithTimeout(context.Background(), s.conf.ShutdoundTimeout)
	defer cancel()
	err = s.srv.Shutdown(ctxShutdown)

	return err
}
