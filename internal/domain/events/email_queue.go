package events

import "time"

type EmailQueueEvent struct {
	Id        string            `json:"id"`
	Recipient string            `json:"recipient"`
	Template  string            `json:"template"`
	Variables map[string]string `json:"variables"`
	SendAt    time.Time         `json:"send_at"`
}

func (e *EmailQueueEvent) GetEventType() string {
	return "email_queue"
}
