package service

import (
	"context"
	"fmt"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"proxy/config"
	"proxy/internal/models"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "proxy/usergrpc/user"
)

const (
	grpcProtocol    = "grpc"
	rpcProtocol     = "rpc"
	jsonrpcProtocol = "json-rpc"
)

func GetlientRPC(protocol string, conf config.UserRPC) (Userer, error) {
	switch protocol {
	case rpcProtocol:
		client, err := newClient(conf, protocol)
		if err != nil {
			return nil, err
		}
		return NewUserRPCClient(client), nil
	case jsonrpcProtocol:
		client, err := newClient(conf, protocol)
		if err != nil {
			return nil, err
		}
		return NewUserRPCClient(client), nil
	case grpcProtocol:
		conn, err := grpc.Dial(fmt.Sprintf("%s:%s", conf.Host, conf.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Println("grpc server connect error:", err)
			return nil, err
		}

		client := pb.NewUsererClient(conn)
		return NewUserGRPCClient(client), nil
	default:
		return nil, fmt.Errorf("invalid protocol")
	}
}

type UserRPCClient struct {
	client *rpc.Client
}

func NewUserRPCClient(client *rpc.Client) *UserRPCClient {
	return &UserRPCClient{
		client: client,
	}
}

func (u *UserRPCClient) Profile(ctx context.Context, email string) (models.User, error) {
	in := ProfileIn{Email: email}
	var out ProfileOut
	err := u.client.Call("UserServiceRPC.Profile", in, &out)
	if err != nil {
		return models.User{}, err
	}

	return models.User{Name: out.Name, Email: out.Email}, nil
}

func (u *UserRPCClient) List(ctx context.Context) ([]models.User, error) {
	var (
		in  ListIn
		out ListOut
	)
	err := u.client.Call("UserServiceRPC.List", in, &out)
	if err != nil {
		return nil, err
	}

	return out.Users, nil
}

type UserGRPCClient struct {
	client pb.UsererClient
}

func NewUserGRPCClient(client pb.UsererClient) *UserGRPCClient {
	return &UserGRPCClient{
		client: client,
	}
}

func (u *UserGRPCClient) Profile(ctx context.Context, email string) (models.User, error) {
	res, err := u.client.Profile(ctx, &pb.ProfileRequest{Email: email})
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Name:  res.Name,
		Email: res.Email,
		Phone: res.Phone,
	}

	return user, nil
}

func (u *UserGRPCClient) List(ctx context.Context) ([]models.User, error) {
	res, err := u.client.List(ctx, &pb.ListRequest{})
	if err != nil {
		return nil, err
	}
	users := make([]models.User, len(res.Users))
	for i, user := range res.Users {
		u := models.User{
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Phone,
		}
		users[i] = u
	}

	return users, nil
}

func newClient(conf config.UserRPC, protocol string) (*rpc.Client, error) {
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
