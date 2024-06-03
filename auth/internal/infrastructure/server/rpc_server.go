package server

import (
	"context"
	"fmt"
	"microservices/auth/config"
	"microservices/auth/internal/modules/auth/service"
	"microservices/auth/rpc/auth"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	pb "microservices/auth/authgrpc/auth"
)

const (
	grpcProtocol    = "grpc"
	rpcProtocol     = "rpc"
	jsonrpcProtocol = "json-rpc"
)

type Server interface {
	Serve(ctx context.Context) error
}

func GetServerRPC(conf config.RPCServer, authService service.Auther, logger *zap.Logger) (Server, error) {

	switch conf.Type {
	case grpcProtocol:

		return NewServerGRPC(conf, auth.NewUserServiceGRPC(authService), logger), nil
	case rpcProtocol:
		authRPC := auth.NewAuthServiceRPC(authService)
		RPCServer := rpc.NewServer()
		err := RPCServer.Register(authRPC)
		if err != nil {
			return nil, err
		}
		return NewServerRPC(conf, RPCServer, logger)
	case jsonrpcProtocol:
		authRPC := auth.NewAuthServiceRPC(authService)
		RPCServer := rpc.NewServer()
		err := RPCServer.Register(authRPC)
		if err != nil {
			return nil, err
		}
		return NewServerRPC(conf, RPCServer, logger)
	default:
		return nil, fmt.Errorf("invalid protocol")
	}
}

type ServerRPC struct {
	conf   config.RPCServer
	srv    *rpc.Server
	logger *zap.Logger
}

func NewServerRPC(conf config.RPCServer, server *rpc.Server, logger *zap.Logger) (Server, error) {
	switch conf.Type {
	case rpcProtocol:
		return &ServerRPC{conf: conf, srv: server, logger: logger}, nil
	case jsonrpcProtocol:
		return &ServerJSONRPC{conf: conf, srv: server, logger: logger}, nil
	default:
		return nil, fmt.Errorf("invalid protocol")
	}
}

func (s *ServerRPC) Serve(ctx context.Context) error {
	var err error

	chErr := make(chan error)
	go func() {
		var l net.Listener
		l, err = net.Listen("tcp", fmt.Sprintf(":%s", s.conf.Port))
		if err != nil {
			//log.Println("rpc server register error")
			s.logger.Error("rpc server register error", zap.Error(err))
			chErr <- err
		}

		s.logger.Info("rpc server started", zap.String("addr", l.Addr().String()))
		var conn net.Conn
		for {
			select {
			case <-ctx.Done():
				s.logger.Error("rpc: stopping server")
				return
			default:

				conn, err = l.Accept()
				if err != nil {
					s.logger.Error("json rpc: net tcp accpet error:", zap.Error(err))
				}
				go s.srv.ServeConn(conn)
			}
		}
	}()

	select {
	case <-chErr:
		return err
	case <-ctx.Done():
	}

	return err
}

type ServerJSONRPC struct {
	conf   config.RPCServer
	srv    *rpc.Server
	logger *zap.Logger
}

func (s *ServerJSONRPC) Serve(ctx context.Context) error {
	var err error

	chErr := make(chan error)
	go func() {
		var l net.Listener
		l, err = net.Listen("tcp", fmt.Sprintf(":%s", s.conf.Port))
		if err != nil {
			s.logger.Error("json rpc server register error:", zap.Error(err))
			chErr <- err
		}

		s.logger.Info("json rpc server started", zap.String("addr", l.Addr().String()))
		var conn net.Conn
		for {
			select {
			case <-ctx.Done():
				s.logger.Error("json rpc: stopping server")
				return
			default:
				conn, err = l.Accept()
				if err != nil {
					s.logger.Error("json rpc: net tcp accept error:", zap.Error(err))
				}
				go s.srv.ServeCodec(jsonrpc.NewServerCodec(conn))
			}
		}
	}()

	select {
	case <-chErr:
		return err
	case <-ctx.Done():
	}

	return err
}

type ServerGRPC struct {
	conf   config.RPCServer
	srv    *grpc.Server
	auth   *auth.AuthServiceGRPC
	logger *zap.Logger
}

func NewServerGRPC(conf config.RPCServer, auth *auth.AuthServiceGRPC, logger *zap.Logger) Server {
	return &ServerGRPC{
		conf:   conf,
		srv:    grpc.NewServer(),
		auth:   auth,
		logger: logger,
	}
}

func (s *ServerGRPC) Serve(ctx context.Context) error {
	var err error

	chErr := make(chan error)
	go func() {
		var l net.Listener
		l, err = net.Listen("tcp", fmt.Sprintf(":%s", s.conf.Port))
		if err != nil {
			s.logger.Error("gRPC server register error:", zap.Error(err))
			chErr <- err
		}

		s.logger.Info("gRPC server started", zap.String("addr", l.Addr().String()))

		pb.RegisterAutherServer(s.srv, s.auth)

		if err = s.srv.Serve(l); err != nil {
			s.logger.Error("grpc server error: ", zap.Error(err))
			chErr <- err
		}

	}()

	select {
	case <-chErr:
		return err
	case <-ctx.Done():
		s.srv.GracefulStop()
	}
	return err
}
