package component

import (
	"microservices/geo/config"

	"gitlab.com/ptflp/gopubsub/queue"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
)

type Components struct {
	Conf      config.AppConf
	RateLimit ratelimit.Limiter
	Logger    *zap.Logger
	MQ        queue.MessageQueuer
}

func NewComponents(conf config.AppConf, logger *zap.Logger, rateLimit ratelimit.Limiter, mq queue.MessageQueuer) *Components {
	return &Components{
		Conf:      conf,
		Logger:    logger,
		RateLimit: rateLimit,
		MQ:        mq,
	}
}
