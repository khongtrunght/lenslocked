package models

import (
	"fmt"

	"github.com/wneessen/go-mail"
)

const (
	DefaultSender = "support@khongtrunght.com"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Email struct {
	From      string
	To        string
	Subject   string
	PlainText string
	HTML      string
}

func NewEmailService(config SMTPConfig) (*EmailService, error) {
	client, err := mail.NewClient(
		config.Host,
		mail.WithPort(config.Port),
		mail.WithUsername(config.Username),
		mail.WithPassword(config.Password),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
	)
	if err != nil {
		return nil, fmt.Errorf("new email service: %w", err)
	}

	es := &EmailService{
		// TODO: Set the default sender
		client: client,
	}
	return es, nil
}

type EmailService struct {
	// DefaultSender is used as the default sender when one isn't provided for an
	// email. This is also used in functions where the email is a predetermined,
	// like the forgot password email.
	DefaultSender string

	client *mail.Client
}

func (es *EmailService) ForgotPassword(to, resetURL string) error {
	email := Email{
		To:        to,
		Subject:   "Reset your password",
		HTML:      fmt.Sprintf(`<p>To reset your password, click <a href="%s">%s</a>.</p>`, resetURL, resetURL),
		PlainText: "To reset your password, click here: " + resetURL,
	}

	if err := es.Send(email); err != nil {
		return fmt.Errorf("forgot password: %w", err)
	}
	return nil
}

func (es *EmailService) Send(email Email) error {
	message := mail.NewMsg()
	es.setFrom(message, email)
	message.To(email.To)
	message.Subject(email.Subject)

	switch {
	case email.PlainText != "" && email.HTML != "":
		message.SetBodyString(mail.TypeTextPlain, email.PlainText)
		message.AddAlternativeString(mail.TypeTextHTML, email.HTML)
	case email.PlainText != "":
		message.SetBodyString(mail.TypeTextPlain, email.PlainText)
	case email.HTML != "":
		message.SetBodyString(mail.TypeTextHTML, email.HTML)
	}

	if err := es.client.DialAndSend(message); err != nil {
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func (es *EmailService) setFrom(message *mail.Msg, email Email) {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case es.DefaultSender != "":
		from = es.DefaultSender
	default:
		from = DefaultSender
	}
	message.From(from)
}
