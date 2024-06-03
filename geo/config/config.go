package config

import (
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
)

const (
	AppName = "APP_NAME"
)

type AppConf struct {
	AppName   string
	DB        DB
	Cache     Cache
	RPCServer RPCServer
	Logger    Logger
	MQ        MQ
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

type MQ struct {
	Host    string
	Port    string
	Type    string
	GroupID string
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

type Cache struct {
	Host string
	Port string
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
		MQ: MQ{
			Host:    os.Getenv("MQ_HOST"),
			Port:    os.Getenv("MQ_PORT"),
			Type:    os.Getenv("MQ_TYPE"),
			GroupID: os.Getenv("MQ_GROUP"),
		},
		Cache: Cache{
			Host: os.Getenv("REDIS_HOST"),
			Port: os.Getenv("REDIS_PORT"),
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
