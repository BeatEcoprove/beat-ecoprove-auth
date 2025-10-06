package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/BeatEcoprove/identityService/config"
	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	ctx    context.Context
	writer *kafka.Writer
}

const (
	version      = 1
	currentTopic = interfaces.AuthEventTopic
)

var kafkaConnection *KafkaPublisher = nil

func GetKafkaPublisher() (*KafkaPublisher, error) {
	if kafkaConnection != nil {
		return kafkaConnection, nil
	}

	env := config.GetConfig()

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  setBrokers(env),
		Balancer: &kafka.LeastBytes{},
	})

	return &KafkaPublisher{
		writer: writer,
		ctx:    context.Background(),
	}, nil
}

func (kc *KafkaPublisher) Publish(payload interfaces.BrokerPayload, topic interfaces.BrokerScope) error {
	event, err := createBaseEvent(payload)

	if err != nil {
		return err
	}

	jsonPayload, err := json.Marshal(event)

	if err != nil {
		return err
	}

	return kc.writer.WriteMessages(kc.ctx, kafka.Message{
		Key:   []byte(event.Key),
		Topic: string(topic),
		Value: jsonPayload,
	})
}

func createBaseEvent(payload interfaces.BrokerPayload) (*interfaces.BrokerMessage, error) {
	eventPayload, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	return &interfaces.BrokerMessage{
		Key:     uuid.NewString(),
		Version: version,
		Metadata: interfaces.BrokerMetadata{
			Source: string(currentTopic),
		},
		Payload:    eventPayload,
		EventType:  payload.GetEventType(),
		OccurredAt: time.Now(),
	}, nil
}

func setBrokers(env *config.Config) []string {
	set := func(host string, port int) string {
		return fmt.Sprintf("%s:%d", host, port)
	}

	return []string{set(env.KAFKA_HOST, env.KAFKA_PORT)}
}

func (kc *KafkaPublisher) Close() error {
	return kc.writer.Close()
}
