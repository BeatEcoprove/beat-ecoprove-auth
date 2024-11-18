package adapters

import (
	"time"
)

type (
	PushMetadata struct {
		Timestamp string `json:"timestamp"`
		RequestId int    `json:"request_id"`
		ServiceId string `json:"service_id"`
	}

	PushMessage struct {
		EventType string       `json:"event_type"`
		Payload   interface{}  `json:"payload"`
		Metadata  PushMetadata `json:"metadata"`
	}

	EmailPayload struct {
		To        string            `json:"to"`
		Subject   string            `json:"subject"`
		Paramters map[string]string `json:"paramters"`
	}

	RabbitMq interface {
		PublishMessage(payload *PushMessage) error
		Close() error
	}
)

var request_id int = 0

func PushEmail(payload EmailPayload) *PushMessage {
	request_id++

	return &PushMessage{
		EventType: "send_email",
		Payload:   payload,
		Metadata: PushMetadata{
			Timestamp: time.Now().GoString(),
			RequestId: request_id,
			ServiceId: "identity",
		},
	}
}
