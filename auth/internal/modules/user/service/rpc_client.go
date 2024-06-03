package service

import (
	"context"
	"fmt"
	"log"
	"microservices/auth/config"
	"microservices/auth/internal/models"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "microservices/auth/usergrpc/user"
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
		log.Printf("auth GetClientRPC host %s port %s", conf.Host, conf.Port)
		conn, err := grpc.Dial(fmt.Sprintf("%s:%s", conf.Host, conf.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Println("grpc server connect error:", err)
			return nil, err
		}
		log.Println("auth: GetlientRPC: connected ")
		client := pb.NewUsererClient(conn)
		return NewUserGRPCClient(client), nil
	default:
		return nil, fmt.Errorf("invalid protocol %s", protocol)
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
	var out PrfileOut
	err := u.client.Call("UserServiceRPC.Profile", in, &out)
	if err != nil {
		return models.User{}, err
	}

	return models.User{Email: out.Email, Password: out.Password}, nil
}

func (u *UserRPCClient) Create(ctx context.Context, user models.User) error {
	in := CreateIn{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Phone:    user.Phone,
	}
	var out CreateOut

	err := u.client.Call("UserServiceRPC.Create", in, &out)
	if err != nil {
		return err
	}

	return nil
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
		ID:       int(res.Id),
		Name:     res.Name,
		Email:    res.Email,
		Password: res.Passwrd,
		Phone:    res.Phone,
	}

	return user, nil
}

func (u *UserGRPCClient) Create(ctx context.Context, user models.User) error {
	_, err := u.client.Create(ctx, &pb.CreateRequest{Name: user.Name, Email: user.Email, Password: user.Password, Phone: user.Phone})
	if err != nil {
		return err
	}

	return nil
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
