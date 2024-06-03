package provider

import (
	"fmt"
	"microservices/notify/config"
	"net/smtp"

	"go.uber.org/zap"
)

const (
	textPlain = "text/plain"
	textHtml  = "text/html; charset=UTF-8"

	TextPlain = iota + 1
	TextHtml
)

type Email struct {
	conf   config.Email
	client smtp.Auth
	addr   string
	logger *zap.Logger
}

func NewEmail(conf config.Email, logger *zap.Logger) *Email {
	emailAuth := smtp.PlainAuth("", conf.Credentials.Login, conf.Credentials.Password, conf.Credentials.Host)

	return &Email{conf: conf, client: emailAuth, addr: fmt.Sprintf("%s:%s", conf.Credentials.Host, conf.Port), logger: logger}
}

func (e *Email) Send(in SendIn) error {
	emailBody := string(in.Data)
	var contentType string

	switch in.Type {
	case TextPlain:
		contentType = textPlain
	case TextHtml:
		contentType = textHtml
		emailBody = `
		<!DOCTYPE html>
		<html>
		<head>
			<title>` + in.Title + `</title>
			<style>
				body { font-family: 'Arial', sans-serif; }
				.container { background-color: #f0f0f0; padding: 20px; margin: 10px auto; width: 600px; }
				.header { background: #007bff; color: #ffffff; padding: 10px; text-align: center; }
				.content { padding: 20px; }
				.footer { background: #333; color: #ffffff; padding: 10px; text-align: center; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h2>` + in.Title + `</h2>
				</div>
				<div class="content">` + string(in.Data) + `</div>
				<div class="footer">
					<p>С уважением, ваша команда.</p>
				</div>
			</div>
		</body>
		</html>
		`
	default:
		contentType = textPlain
	}

	mime := "MIME-version: 1.0;\nContent-Type: " + contentType + "; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + in.Title + "\n"
	msg := []byte(subject + mime + "\n" + emailBody)

	if err := smtp.SendMail(e.addr, e.client, e.conf.From, []string{in.To}, msg); err != nil {
		e.logger.Error("email: sent msg err", zap.Error(err))
		return err
	}

	return nil
}

// var (
// 	messageHTML = `
// 	<!DOCTYPE html>
// 	<html>
// 	<head>
// 		<title>` + in.Title + `</title>
// 		<style>
// 			body { font-family: 'Arial', sans-serif; }
// 			.container { background-color: #f0f0f0; padding: 20px; margin: 10px auto; width: 600px; }
// 			.header { background: #007bff; color: #ffffff; padding: 10px; text-align: center; }
// 			.content { padding: 20px; }
// 			.footer { background: #333; color: #ffffff; padding: 10px; text-align: center; }
// 		</style>
// 	</head>
// 	<body>
// 		<div class="container">
// 			<div class="header">
// 				<h2>` + in.Title + `</h2>
// 			</div>
// 			<div class="content">` + string(in.Data) + `</div>
// 			<div class="footer">
// 				<p>С уважением, ваша команда.</p>
// 			</div>
// 		</div>
// 	</body>
// 	</html>
// 	`
// )
