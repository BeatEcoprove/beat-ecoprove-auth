package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/BeatEcoprove/identityService/config"
	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	ctx    context.Context
	cancel context.CancelFunc
	reader *kafka.Reader
}

type HandlerBind struct {
	handler interfaces.Handler
	event   interfaces.BrokerPayload
}

const groupId = "auth_consumer"

var kafkaConsumer *KafkaConsumer
var eventHandlers map[string]HandlerBind = make(map[string]HandlerBind)
var consumeTopics = []string{
	string(interfaces.MessagingEventTopic),
}

func GetKafkaConsumer() (*KafkaConsumer, error) {
	if kafkaConsumer != nil {
		return kafkaConsumer, nil
	}

	env := config.GetConfig()
	ctx, cancel := context.WithCancel(context.Background())

	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     setBrokers(env),
			GroupID:     groupId,
			GroupTopics: consumeTopics,
			MinBytes:    1e3,
			MaxBytes:    10e6,
		}),
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (kc *KafkaConsumer) Register(handler interfaces.Handler, event interfaces.BrokerPayload) error {
	eventType, err := getEventType(event)

	if err != nil {
		return err
	}

	eventHandlers[eventType] = HandlerBind{handler: handler, event: event}
	return nil
}

func (kc *KafkaConsumer) Consume() {
	log.Println("🎧 Kafka consumer started for multiple topics")

	for {
		message, err := kc.reader.ReadMessage(kc.ctx)

		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		var event interfaces.BrokerMessage
		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("Failed to unmarshal message from %s: %v", message.Topic, err)
			continue
		}

		kc.handleEvent(event)
	}
}

func getEventType(payload any) (string, error) {
	t := reflect.TypeOf(payload)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	eventType := t.Name()

	if !strings.Contains(eventType, "Event") {
		return "", errors.New("provide a valid event")
	}

	eventType = strings.TrimSuffix(eventType, "Event")
	return toSnakeCase(eventType), nil
}

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

func (kc *KafkaConsumer) handleEvent(event interfaces.BrokerMessage) {
	bind, ok := eventHandlers[event.EventType]

	if !ok {
		log.Printf("event type isn't register yet, %+v", event)
		return
	}

	if err := json.Unmarshal(event.Payload, bind.event); err != nil {
		log.Printf("Failed to unmarshal event")
		return
	}

	if err := bind.handler.Call(bind.event); err != nil {
		log.Printf(err.Error())
	}
}

func (kc *KafkaConsumer) Close() error {
	kc.cancel()

	if err := kc.reader.Close(); err != nil {
		return err
	}

	log.Println("Kafka consumer stopped.")
	return nil
}
