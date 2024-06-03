package provider

import (
	"bytes"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

const (
	sendSmsURL = "http://81.163.28.166:8080/api/sms/send"
)

type Phone struct {
	client *http.Client
	logger *zap.Logger
}

type SMS struct {
	Phones string `json:"phones"`
	Mes    string `json:"mes"`
}

func NewPhone(logger *zap.Logger) *Phone {

	return &Phone{
		client: &http.Client{},
		logger: logger,
	}
}

func (p *Phone) Send(in SendIn) error {
	sms := SMS{
		Phones: in.Phone,
		Mes:    string(in.Data),
	}
	jsonSMS, err := json.Marshal(sms)
	if err != nil {
		p.logger.Error("marshal sms err", zap.Error(err))
		return err
	}

	req, err := http.NewRequest("POST", sendSmsURL, bytes.NewBuffer(jsonSMS))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		//p.logger.Error("ошибка при отправке SMS", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	return nil
}
