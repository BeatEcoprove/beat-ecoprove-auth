package adapters

import (
	"context"
	"encoding/json"

	"github.com/BeatEcoprove/identityService/config"
	interfaces "github.com/BeatEcoprove/identityService/pkg/adapters"
	amqp "github.com/rabbitmq/amqp091-go"
)

const PayloadContentType = "application/json"

var rabbitmq *RabbitMQConnection = nil

type RabbitMQConnection struct {
	rabbitmq *amqp.Connection
	ctx      context.Context
}

func buildURL() string {
	env := config.GetCofig()

	return "amqp://" + env.RABBITMQ_DEFAULT_USER + ":" + env.RABBITMQ_DEFAULT_PASS + "@" +
		env.RABBIT_MQ_HOST + ":" + env.RABBIT_MQ_PORT + "/" + env.RABBITMQ_DEFAULT_VHOST
}

func GetRabbitMqConnection() (*RabbitMQConnection, error) {
	if rabbitmq != nil {
		return rabbitmq, nil
	}

	connection, err := amqp.Dial(buildURL())

	if err != nil {
		return nil, err
	}

	return &RabbitMQConnection{
		rabbitmq: connection,
		ctx:      context.Background(),
	}, nil
}

func CreateAndBindQueue(ch *amqp.Channel) error {
	env := config.GetCofig()

	_, err := ch.QueueDeclare(
		env.RABBIT_MQ_QUEUE_MAIL, // name of the queue
		true,                     // durable (survives server restart)
		false,                    // delete when unused
		false,                    // exclusive (can only be used by the connection)
		false,                    // no-wait
		nil,                      // arguments
	)

	if err != nil {
		return err
	}

	// Bind the queue to the exchange
	err = ch.QueueBind(
		env.RABBIT_MQ_QUEUE_MAIL,  // queue name
		env.RABBIT_MQ_ROUTING_KEY, // routing key
		env.RABBIT_MQ_EXCHANGE,    // exchange name
		false,                     // no-wait
		nil,                       // arguments
	)

	if err != nil {
		return err
	}

	return nil
}

func (rc *RabbitMQConnection) PublishMessage(payload *interfaces.PushMessage) error {
	env := config.GetCofig()
	channel, err := rc.rabbitmq.Channel()

	if err != nil {
		return err
	}
	defer channel.Close()

	err = channel.ExchangeDeclare(
		env.RABBIT_MQ_EXCHANGE, // Exchange name
		"direct",               // Exchange type
		true,                   // Durable
		false,                  // Auto-deleted
		false,                  // Internal
		false,                  // No-wait
		nil,                    // Arguments
	)
	if err != nil {
		return err
	}

	CreateAndBindQueue(channel)

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	err = channel.PublishWithContext(
		rc.ctx,
		env.RABBIT_MQ_EXCHANGE,    // Exchange
		env.RABBIT_MQ_ROUTING_KEY, // Routing key
		false,                     // Mandatory
		false,                     // Immediate
		amqp.Publishing{
			ContentType: PayloadContentType,
			Body:        jsonPayload,
		},
	)

	return err
}

func (rc *RabbitMQConnection) Close() error {
	return rc.rabbitmq.Close()
}
