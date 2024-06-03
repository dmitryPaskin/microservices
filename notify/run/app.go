package run

import (
	"bytes"
	"context"
	"encoding/json"
	"microservices/notify/config"
	"microservices/notify/internal/infrastructure/errors"
	notifyservice "microservices/notify/internal/infrastructure/service"
	"microservices/notify/internal/modules/mq"
	"microservices/notify/internal/modules/user/service"
	"microservices/notify/internal/provider"

	"go.uber.org/zap"
)

var (
	limitExceededTitle   = "Rate Limit Exceeded"
	limitExceededMessage = []byte("The geoservice's rate limit of requests per minute has been exceeded. Please try again later.")
)

type App struct {
	conf   config.AppConf
	logger *zap.Logger
}

func NewApp(conf config.AppConf, logger *zap.Logger) *App {
	return &App{conf: conf, logger: logger}
}

func (a *App) Run() int {

	email := provider.NewEmail(a.conf.Provider.Email, a.logger)

	phone := provider.NewPhone(a.logger)

	notifier := notifyservice.NewNotify(a.conf.Provider.Email, email, phone, a.logger)

	// url := fmt.Sprintf("amqp://guest:guest@%s:%s/", a.conf.MQ.Host, a.conf.MQ.Port)
	// conn, err := amqp.Dial(url)
	// if err != nil {
	// 	a.logger.Fatal("connec rabbitMQ err", zap.Error(err))
	// }
	// defer conn.Close()

	mq, err := mq.GetMessageQueue(a.conf.MQ, a.logger)
	if err != nil {
		a.logger.Fatal("create rmq err", zap.Error(err))
	}
	defer mq.Close()

	messages, err := mq.Subscribe("rate_limit")
	if err != nil {
		a.logger.Fatal("subscribe rate_limit err", zap.Error(err))
	}

	userClient, err := service.NewUserGRPCClient(a.conf.UserRPC)
	if err != nil {
		a.logger.Fatal("init user grpc client err", zap.Error(err))
	}

	// forever := make(chan struct{})
	// a.logger.Info("Started")
	// var errCode int
	// go func() {
	// 	for msg := range messages {

	// 		m, err := readMessage(msg.Data)
	// 		if err != nil {
	// 			a.logger.Error("decode data from message err", zap.Error(err))
	// 			continue
	// 		}

	// 		out := notifier.Push(service.PushIn{
	// 			Identifier: m.Email,
	// 			Phone:      m.Phone,
	// 			Type:       service.PushEmailAndPhone,
	// 			Title:      limitExceededTitle,
	// 			Data:       limitExceededMessage,
	// 		})

	// 		if out.ErrorCode != 0 {
	// 			errCode = out.ErrorCode
	// 			close(forever)
	// 			return
	// 		}
	// 	}
	// }()

	// <-forever
	// return errCode

	errChan := make(chan int, 1)
	a.logger.Info("Started")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context) {
		for {
			select {
			case msg, ok := <-messages:
				if !ok {
					a.logger.Error("mq chan clossed")
					errChan <- errors.RabbitMqClosedChan
					return
				}
				m, err := readMessage(msg.Data)
				if err != nil {
					a.logger.Error("decode data from message err", zap.Error(err))
					continue
				}

				user, err := userClient.GetByID(ctx, m.ID)
				if err != nil {
					a.logger.Error("get user by id err", zap.Error(err))
					continue
				}

				out := notifier.Push(notifyservice.PushIn{
					Identifier: user.Email,
					Phone:      user.Phone,
					Type:       notifyservice.PushEmailAndPhone,
					Title:      limitExceededTitle,
					Data:       limitExceededMessage,
				})
				if out.ErrorCode != 0 {
					errChan <- out.ErrorCode
					return
				}
				//
				if err := mq.Ack(&msg); err != nil {
					a.logger.Error("Failed to acknowledge message", zap.Error(err))
				}
				//
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	select {
	case errCode := <-errChan:
		return errCode
	case <-ctx.Done():
		return 0
	}
}

type RateLimitMsg struct {
	ID int
}

func readMessage(data []byte) (RateLimitMsg, error) {
	var msg RateLimitMsg
	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&msg); err != nil {
		return RateLimitMsg{}, err
	}

	return msg, nil
}
