package mq

import (
	"fmt"
	"microservices/notify/config"

	"github.com/streadway/amqp"
	"gitlab.com/ptflp/gopubsub/kafkamq"
	"gitlab.com/ptflp/gopubsub/queue"
	"gitlab.com/ptflp/gopubsub/rabbitmq"
	"go.uber.org/zap"
)

func GetMessageQueue(conf config.MQ, logger *zap.Logger) (queue.MessageQueuer, error) {
	switch conf.Type {
	case "rabbit":
		return newRabbitMQ(conf, logger)
	case "kafka":
		return newKafkaMQ(conf)
	default:
		return nil, fmt.Errorf("invalid MQ type")
	}
}

func newRabbitMQ(conf config.MQ, logger *zap.Logger) (queue.MessageQueuer, error) {
	url := fmt.Sprintf("amqp://guest:guest@%s:%s/", conf.Host, conf.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Error("connect rabbitMQ err", zap.Error(err))
		return nil, err
	}

	return rabbitmq.NewRabbitMQ(conn)
}

func newKafkaMQ(conf config.MQ) (queue.MessageQueuer, error) {
	broker := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	return kafkamq.NewKafkaMQ(broker, conf.GroupID)
}
