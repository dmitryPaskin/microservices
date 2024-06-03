package config

import (
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
)

const (
	AppName = "APP_NAME"

	serverPort         = "SERVER_PORT"
	envShutdownTimeout = "SHUTDOWN_TIMEOUT"

	parseShutdownTimeoutError    = "config: parse server shutdown timeout error"
	parseRpcShutdownTimeoutError = "config: parse rpc server shutdown timeout error"
)

type AppConf struct {
	AppName   string
	Server    Server
	DB        DB
	Cache     Cache
	RPCServer RPCServer
	GeoRPC    GeoRPC
	AuthRPC   AuthRPC
	UserRPC   UserRPC
	Logger    Logger
}

type DB struct {
	Driver   string
	Name     string
	User     string
	Password string
	Host     string
	Port     string
	MaxConn  int
	Timeout  int
}

type Server struct {
	Port             string
	ShutdoundTimeout time.Duration
}

type RPCServer struct {
	Port          string
	ShutdoundTime time.Duration
	Type          string
}

type GeoRPC struct {
	Host string
	Port string
}

type AuthRPC struct {
	Host string
	Port string
}

type UserRPC struct {
	Host string
	Port string
}

type Cache struct {
	Host string
	Port string
}

type Logger struct {
	Level string
}

func NewAppConf() AppConf {
	port := os.Getenv(serverPort)

	return AppConf{
		AppName: os.Getenv("APP_NAME"),
		Server: Server{
			Port: port,
		},
		DB: DB{
			Driver:   os.Getenv("DB_DRIVER"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
		Cache: Cache{
			Host: os.Getenv("REDIS_HOST"),
			Port: os.Getenv("REDIS_PORT"),
		},
	}

}

func (a *AppConf) Init(logger *zap.Logger) {
	shutDownTimeOut, err := strconv.Atoi(os.Getenv(envShutdownTimeout))
	if err != nil {
		logger.Fatal(parseShutdownTimeoutError)
	}
	shutDownTimeout := time.Duration(shutDownTimeOut) * time.Second

	dbTimeout, err := strconv.Atoi(os.Getenv("DB_TIMEOUT"))
	if err != nil {
		logger.Fatal("config: parse db timeout err", zap.Error(err))
	}
	dbMaxConn, err := strconv.Atoi(os.Getenv("MAX_CONN"))
	if err != nil {
		logger.Fatal("config: parse db max connection err", zap.Error(err))
	}
	a.DB.Timeout = dbTimeout
	a.DB.MaxConn = dbMaxConn

	a.Server.ShutdoundTimeout = shutDownTimeout

	a.RPCServer.Port = os.Getenv("RPC_PORT")
	a.RPCServer.Type = os.Getenv("RPC_PROTOCOL")
	a.GeoRPC.Host = os.Getenv("GEO_RPC_HOST")
	a.GeoRPC.Port = os.Getenv("GEO_RPC_PORT")

	a.AuthRPC.Port = os.Getenv("AUTH_RPC_PORT")
	a.AuthRPC.Host = os.Getenv("AUTH_RPC_HOST")

	a.UserRPC.Port = os.Getenv("USER_RPC_PORT")
	a.UserRPC.Host = os.Getenv("USER_RPC_HOST")
}
