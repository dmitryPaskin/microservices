package config

import (
	"os"
	"time"

	"go.uber.org/zap"
)

const (
	AppName = "APP_NAME"

	rpcServerPort      = "RPC_SERVER_PORT"
	envShutdownTimeout = "SHUTDOWN_TIMEOUT"
)

type AppConf struct {
	AppName   string
	RPCServer RPCServer
	UserRPC   UserRPC
	Logger    Logger
}

type RPCServer struct {
	Port          string
	ShutdoundTime time.Duration
	Type          string
}

type UserRPC struct {
	Host string
	Port string
}

type Logger struct {
	Level string
}

func NewAppConf() AppConf {
	port := os.Getenv(rpcServerPort)

	return AppConf{
		AppName: os.Getenv("APP_NAME"),
		RPCServer: RPCServer{
			Port: port,
		},
	}

}

func (a *AppConf) Init(logger *zap.Logger) {

	a.RPCServer.Type = os.Getenv("RPC_PROTOCOL")
	a.UserRPC.Host = os.Getenv("USER_RPC_HOST")
	a.UserRPC.Port = os.Getenv("USER_RPC_PORT")
}
