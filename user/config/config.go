package config

import (
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type AppConf struct {
	AppName   string
	DB        DB
	RPCServer RPCServer
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

type RPCServer struct {
	Port          string
	ShutdoundTime time.Duration
	Type          string
}

type Logger struct {
	Level string
}

func NewAppConf() AppConf {

	return AppConf{
		AppName: os.Getenv("APP_NAME"),
		DB: DB{
			Driver:   os.Getenv("DB_DRIVER"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
	}

}

func (a *AppConf) Init(logger *zap.Logger) {

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

	a.RPCServer.Port = os.Getenv("RPC_PORT")
	a.RPCServer.Type = os.Getenv("RPC_PROTOCOL")

}
