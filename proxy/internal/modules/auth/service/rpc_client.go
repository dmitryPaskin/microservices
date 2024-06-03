package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"proxy/config"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	pb "proxy/authgrpc/auth"
)

const (
	grpcProtocol    = "grpc"
	rpcProtocol     = "rpc"
	jsonrpcProtocol = "json-rpc"
)

func GetlientRPC(protocol string, conf config.AuthRPC) (Auther, error) {
	switch protocol {
	case rpcProtocol:
		client, err := newClient(conf, protocol)
		if err != nil {
			return nil, err
		}
		return NewAuthRPCClient(client), nil
	case jsonrpcProtocol:
		client, err := newClient(conf, protocol)
		if err != nil {
			return nil, err
		}
		return NewAuthRPCClient(client), nil
	case grpcProtocol:
		conn, err := grpc.Dial(fmt.Sprintf("%s:%s", conf.Host, conf.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Println("grpc server connect error:", err)
			return nil, err
		}

		client := pb.NewAutherClient(conn)
		return NewAuthGRPCClient(client), nil
	default:
		return nil, fmt.Errorf("invalid protocol")
	}
}

type AuthRPCClient struct {
	client *rpc.Client
}

func NewAuthRPCClient(client *rpc.Client) *AuthRPCClient {
	return &AuthRPCClient{
		client: client,
	}
}

func (a *AuthRPCClient) Register(in RegisterIn) RegisterOut {
	var out RegisterOut
	err := a.client.Call("AuthServiceRPC.Register", in, &out)
	if err != nil {
		out.Message = err.Error()
	}

	return out
}

func (a *AuthRPCClient) Login(in LoginIn) LoginOut {
	var out LoginOut
	err := a.client.Call("AuthServiceRPC.Login", in, &out)
	if err != nil {
		out.Message = err.Error()
	}

	return out
}

type AuthGRPCClient struct {
	client pb.AutherClient
}

func NewAuthGRPCClient(client pb.AutherClient) *AuthGRPCClient {
	return &AuthGRPCClient{
		client: client,
	}
}

func (g *AuthGRPCClient) Register(in RegisterIn) RegisterOut {
	res, err := g.client.Register(context.Background(), &pb.RegisterRequest{Email: in.Email, Password: in.Password, Name: in.Name, Phone: in.Phone})
	if err != nil {
		switch status.Code(err) {
		case codes.AlreadyExists:
			return RegisterOut{
				Status:  http.StatusConflict,
				Message: err.Error(),
			}
		default:
			return RegisterOut{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			}
		}

	}

	return RegisterOut{
		Status:  int(res.GetStatus()),
		Message: res.GetMessage(),
	}
}

func (g *AuthGRPCClient) Login(in LoginIn) LoginOut {
	res, _ := g.client.Login(in.Ctx, &pb.LoginRequest{Email: in.Email, Password: in.Password})

	out := LoginOut{
		Success: res.GetSuccess(),
		Message: res.GetMessage(),
	}

	return out
}

func newClient(conf config.AuthRPC, protocol string) (*rpc.Client, error) {
	var (
		client *rpc.Client
		err    error
		host   = conf.Host
		port   = conf.Port
	)

	switch protocol {
	case rpcProtocol:
		client, err = rpc.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		if err != nil {
			return nil, err
		}
		log.Println("rpc client connected")
		return client, nil
	case jsonrpcProtocol:
		// без этого костыля сервер редко успевает запуститься и коннект проваливается
		time.Sleep(1 * time.Second)

		client, err = jsonrpc.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		if err != nil {
			return nil, err
		}
		log.Println("jsonrpc client connected")
		return client, nil

	default:
		return nil, fmt.Errorf("invalid protocol %s", protocol)
	}
}
