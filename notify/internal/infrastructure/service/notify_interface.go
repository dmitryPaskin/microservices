package service

import (
	"log"
	"microservices/notify/config"
	"microservices/notify/internal/infrastructure/errors"
	"microservices/notify/internal/provider"

	"go.uber.org/zap"
)

const (
	PushEmail = iota + 1
	PushPhone
	PushEmailAndPhone
)

type Notifier interface {
	Push(in PushIn) PushOut
}

type Notify struct {
	conf   config.Email
	email  provider.Sender
	phone  provider.Sender
	logger *zap.Logger
}

func NewNotify(conf config.Email, email provider.Sender, phone provider.Sender, logger *zap.Logger) Notifier {
	return &Notify{conf: conf, email: email, phone: phone, logger: logger}
}

func (n *Notify) Push(in PushIn) PushOut {
	var err error

	switch in.Type {
	case PushPhone:
		err = n.phone.Send(provider.SendIn{
			Phone: in.Phone,
			Data:  in.Data,
		})
		if err != nil {
			return PushOut{ErrorCode: errors.NotifyEmailSendErr}
		}
	case PushEmail:
		err = n.email.Send(provider.SendIn{
			To:    in.Identifier,
			From:  n.conf.From,
			Title: in.Title,
			Type:  provider.TextPlain,
			Data:  in.Data,
		})
		if err != nil {
			n.logger.Error("send email err", zap.Error(err))
			return PushOut{
				ErrorCode: errors.NotifyEmailSendErr,
			}
		}
	default:
		_ = n.phone.Send(provider.SendIn{
			Phone: in.Phone,
			Data:  in.Data,
		})
		err = n.email.Send(provider.SendIn{
			To:    in.Identifier,
			From:  n.conf.From,
			Title: in.Title,
			Type:  provider.TextPlain,
			Data:  in.Data,
		})
		if err != nil {
			n.logger.Error("send email err", zap.Error(err))
			return PushOut{
				ErrorCode: errors.NotifyEmailSendErr,
			}
		}
	}
	log.Printf("На почту [%s] должно быть отправлено сообщение от [%s]: [%s]", in.Identifier, n.conf.From, in.Data)
	return PushOut{}
}

type PushIn struct {
	Identifier string
	Phone      string
	Type       int
	Title      string
	Data       []byte
	Options    []interface{}
}

type PushOut struct {
	ErrorCode int
}
