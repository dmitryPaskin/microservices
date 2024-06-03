package component

import (
	"proxy/config"
	"proxy/internal/infrastructure/responder"

	"go.uber.org/zap"
)

type Components struct {
	Conf      config.AppConf
	Responder responder.Responder
	Logger    *zap.Logger
}

func NewComponents(conf config.AppConf, responder responder.Responder, logger *zap.Logger) *Components {
	return &Components{
		Conf:      conf,
		Responder: responder,
		Logger:    logger,
	}
}
