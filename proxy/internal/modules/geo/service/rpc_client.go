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

	"github.com/go-chi/jwtauth/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	pb "proxy/geogrpc"
)

const (
	grpcProtocol    = "grpc"
	rpcProtocol     = "rpc"
	jsonrpcProtocol = "json-rpc"
)

func GetlientRPC(protocol string, conf config.GeoRPC) (Georer, error) {
	switch protocol {
	case rpcProtocol:
		client, err := newClient(conf, protocol)
		if err != nil {
			return nil, err
		}
		return NewGeoRPCClient(client), nil
	case jsonrpcProtocol:
		client, err := newClient(conf, protocol)
		if err != nil {
			return nil, err
		}
		return NewGeoRPCClient(client), nil
	case grpcProtocol:
		conn, err := grpc.Dial(fmt.Sprintf("%s:%s", conf.Host, conf.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Println("grpc server connect error:", err)
			return nil, err
		}

		client := pb.NewGeorerClient(conn)
		return NewGeoGRPCCLient(client), nil
	default:
		return nil, fmt.Errorf("invalid protocol")
	}
}

type GeoRPCClient struct {
	client *rpc.Client
}

func NewGeoRPCClient(client *rpc.Client) *GeoRPCClient {
	return &GeoRPCClient{
		client: client,
	}
}

func (g *GeoRPCClient) SearchAddresses(ctx context.Context, in SearchAddressesIn) SearchAddressesOut {
	var out SearchAddressesOut
	err := g.client.Call("GeoServiceRPC.SearchAddresses", in, &out)
	if err != nil {
		out.Err = err
	}

	return out
}

func (g *GeoRPCClient) GeoCode(in GeoCodeIn) GeoCodeOut {
	var out GeoCodeOut
	err := g.client.Call("GeoServiceRPC.GeoCode", in, &out)
	if err != nil {
		out.Err = err
	}

	return out
}

type GeoGRPCClient struct {
	client pb.GeorerClient
}

func NewGeoGRPCCLient(client pb.GeorerClient) *GeoGRPCClient {
	return &GeoGRPCClient{
		client: client,
	}
}

func (g *GeoGRPCClient) SearchAddresses(ctx context.Context, in SearchAddressesIn) SearchAddressesOut {
	_, claims, _ := jwtauth.FromContext(ctx)

	md := metadata.New(map[string]string{
		"id": claims["id"].(string),
	})

	newCtx := metadata.NewOutgoingContext(ctx, md)

	res, err := g.client.SearchAddresses(newCtx, &pb.SearchAddressesRequest{Query: in.Query})
	if err != nil {
		return SearchAddressesOut{Err: err}
	}
	address := models.Address{
		Lat: res.Address.Lat,
		Lon: res.Address.Lon,
	}

	out := SearchAddressesOut{
		Address: address,
		Err:     err,
	}

	return out
}

func (g *GeoGRPCClient) GeoCode(in GeoCodeIn) GeoCodeOut {
	res, err := g.client.GeoCode(context.Background(), &pb.GeoCodeRequest{Lat: in.Lat, Lng: in.Lng})

	out := GeoCodeOut{
		Lat: res.Lat,
		Lng: res.Lng,
		Err: err,
	}

	return out
}

func newClient(conf config.GeoRPC, protocol string) (*rpc.Client, error) {
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
