package mq

import (
	"fmt"
	"microservices/geo/config"

	"github.com/streadway/amqp"
	"gitlab.com/ptflp/gopubsub/kafkamq"
	"gitlab.com/ptflp/gopubsub/queue"
	"gitlab.com/ptflp/gopubsub/rabbitmq"
	"go.uber.org/zap"
)

// type MessageQueue struct {
// 	mq queue.MessageQueuer
// }

func GetMessageQueue(conf config.MQ, logger *zap.Logger) (queue.MessageQueuer, error) {
	switch conf.Type {
	case "rabbit":
		return newRabbitMQ(conf, logger)
	case "kafka":
		return newKafkaMQ(conf, logger)
	default:
		return nil, fmt.Errorf("invalid mq type %s", conf.Type)
	}
}

func newKafkaMQ(conf config.MQ, logger *zap.Logger) (queue.MessageQueuer, error) {
	broker := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	return kafkamq.NewKafkaMQ(broker, conf.GroupID)

}

func newRabbitMQ(conf config.MQ, logger *zap.Logger) (queue.MessageQueuer, error) {
	url := fmt.Sprintf("amqp://guest:guest@%s:%s/", conf.Host, conf.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Error("mq connection err", zap.Error(err))
		return nil, err
	}

	rmq, err := rabbitmq.NewRabbitMQ(conn)
	if err != nil {
		logger.Error("new rabbbitMQ err", zap.Error(err))
		return nil, err
	}

	if err := rabbitmq.CreateExchange(conn, "rate_limit", "direct"); err != nil {
		logger.Error("create exchange err", zap.Error(err))
		return nil, err
	}

	return rmq, nil
}
