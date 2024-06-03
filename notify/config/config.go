package config

import (
	"os"

	"go.uber.org/zap"
)

type AppConf struct {
	AppName  string
	Logger   Logger
	Provider Provider
	MQ       MQ
	UserRPC  UserRPC
}

type UserRPC struct {
	Host string
	Port string
}

type MQ struct {
	Host    string
	Port    string
	Type    string
	GroupID string
}

type Logger struct {
	Level string
}

type Provider struct {
	Email Email
}

type Email struct {
	From        string
	Port        string
	Credentials Credentials
}

type Credentials struct {
	Host     string `json:"-" yaml:"host"`
	Login    string `json:"-" yaml:"login"`
	Password string `json:"-" yaml:"password"`
}

func NewAppConf() AppConf {
	return AppConf{
		AppName: os.Getenv("APP_NAME"),
		Logger: Logger{
			Level: os.Getenv("LOGGER_LEVEL"),
		},
		MQ: MQ{
			Host:    os.Getenv("MQ_HOST"),
			Port:    os.Getenv("MQ_PORT"),
			Type:    os.Getenv("MQ_TYPE"),
			GroupID: os.Getenv("MQ_GROUP"),
		},
	}
}

func (a *AppConf) Init(logger *zap.Logger) {

	a.Provider.Email.From = os.Getenv("EMAIL_FROM")
	a.Provider.Email.Port = os.Getenv("EMAIL_PORT")
	a.Provider.Email.Credentials.Host = os.Getenv("EMAIL_HOST")
	a.Provider.Email.Credentials.Login = os.Getenv("EMAIL_LOGIN")
	a.Provider.Email.Credentials.Password = os.Getenv("EMAIL_PASSWORD")

	a.UserRPC.Host = os.Getenv("USER_RPC_HOST")
	a.UserRPC.Port = os.Getenv("USER_RPC_PORT")
}
