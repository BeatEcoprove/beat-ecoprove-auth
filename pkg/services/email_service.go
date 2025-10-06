package services

import (
	"errors"
	"time"

	"github.com/BeatEcoprove/identityService/internal/domain/events"
	"github.com/BeatEcoprove/identityService/pkg/adapters"
	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
	"github.com/BeatEcoprove/identityService/pkg/shared"
	"github.com/google/uuid"
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
		broker     interfaces.Broker
		sentEmails []EmailInput
	}

	EmailInput struct {
		To       string `validate:"email"`
		Template *EmailTemplate
	}
)

var (
	ErrEmptyEmailList = errors.New("there isn't any email sent yet")
)

func NewConfirmEmailTemplate() *EmailTemplate {
	return &EmailTemplate{
		Id:        "confirm-account",
		Subject:   "Confirm Account",
		Paramters: make(map[string]string),
	}
}

func NewForgotEmailTemplate(code string) *EmailTemplate {
	return &EmailTemplate{
		Id:      "forgot-password",
		Subject: "Forgot Password",
		Paramters: map[string]string{
			"code": code,
		},
	}
}

func NewEmailService(rabbitmq interfaces.Broker) *EmailService {
	return &EmailService{
		broker: rabbitmq,
	}
}

func (es *EmailService) Last() (*EmailInput, error) {
	emailsLen := len(es.sentEmails)

	if emailsLen == 0 {
		return nil, ErrEmptyEmailList
	}

	return &es.sentEmails[emailsLen-1], nil
}

func (es *EmailService) Send(input EmailInput) error {
	payload, err := convertToPush(input)

	if err != nil {
		return err
	}

	if err := es.broker.Publish(payload, adapters.EmailEventTopic); err != nil {
		return err
	}

	es.sentEmails = append(es.sentEmails, input)
	return nil
}

func convertToPush(input EmailInput) (*events.EmailQueueEvent, error) {
	err := shared.Validate(input)

	if err != nil {
		return nil, err
	}

	return &events.EmailQueueEvent{
		Id:        uuid.NewString(),
		Recipient: input.To,
		Template:  input.Template.Id,
		Variables: input.Template.Paramters,
		SendAt:    time.Now(),
	}, err
}
