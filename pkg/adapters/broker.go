package adapters

import (
	"encoding/json"
	"time"
)

const (
	AuthEventTopic      BrokerScope = "auth_events"
	MessagingEventTopic BrokerScope = "messaging_events"
	EmailEventTopic     BrokerScope = "messaging_events_email"
)

type (
	BrokerScope string

	BrokerPayload interface {
		GetEventType() string
	}

	Handler interface {
		Call(event any) error
	}

	Consumer interface {
		Consume()
		Register(handler Handler, event BrokerPayload) error
		Close() error
	}

	Broker interface {
		Publish(payload BrokerPayload, topic BrokerScope) error
		Close() error
	}

	BrokerMetadata struct {
		Source string `json:"source"`
	}

	BrokerMessage struct {
		Version    int             `json:"version"`
		Metadata   BrokerMetadata  `json:"metadata"`
		Key        string          `json:"key"`
		Payload    json.RawMessage `json:"payload"`
		EventType  string          `json:"event_type"`
		OccurredAt time.Time       `json:"occurred_at"`
	}
)

var request_id int = 0
