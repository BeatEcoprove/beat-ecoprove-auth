package services

import (
	"errors"

	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
	"github.com/BeatEcoprove/identityService/pkg/shared"
)

type (
	EmailTemplate struct {
		Id        string
		Subject   string
		Paramters map[string]string
	}

	IEmailService interface {
		Send(input EmailInput) error
		Last() (*EmailInput, error)
	}

	EmailService struct {
		rabbitMq interfaces.RabbitMq
	}

	EmailInput struct {
		To       string `validate:"email"`
		Template *EmailTemplate
	}
)

var (
	sentEmails = make([]EmailInput, 0)

	ErrEmptyEmailList = errors.New("there isn't any email sent yet")
)

func NewForgotEmailTemplate(code string) *EmailTemplate {
	return &EmailTemplate{
		Id:      "forgot-password",
		Subject: "Forgot Password",
		Paramters: map[string]string{
			"code": code,
		},
	}
}

func NewEmailService(rabbitmq interfaces.RabbitMq) *EmailService {
	return &EmailService{
		rabbitMq: rabbitmq,
	}
}

func (es *EmailService) Last() (*EmailInput, error) {
	emailsLen := len(sentEmails)

	if emailsLen == 0 {
		return nil, ErrEmptyEmailList
	}

	return &sentEmails[emailsLen-1], nil
}

func convertToPush(input EmailInput) (*interfaces.EmailPayload, error) {
	err := shared.Validate(input)

	if err != nil {
		return nil, err
	}

	return &interfaces.EmailPayload{
		To:        input.To,
		Subject:   input.Template.Subject,
		Paramters: input.Template.Paramters,
	}, nil
}

func (es *EmailService) Send(input EmailInput) error {
	payload, err := convertToPush(input)

	if err != nil {
		return err
	}

	if err := es.rabbitMq.PublishMessage(interfaces.PushEmail(*payload)); err != nil {
		return err
	}

	sentEmails = append(sentEmails, input)
	return nil
}
