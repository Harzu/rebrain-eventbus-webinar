package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"

	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/system/eventbus/rabbitmq"
	"github.com/Harzu/rebrain-eventbus-webinar/kafka/pkg/rmqevents"
)

const defaultRoutingKey = ""

type Producer interface {
	OrderCreatedProduce(orderId uint64) error
}

type rmqProducer struct {
	client *rabbitmq.Client
}

func NewProducer(client *rabbitmq.Client) Producer {
	return &rmqProducer{
		client: client,
	}
}

func (p *rmqProducer) OrderCreatedProduce(orderId uint64) error {
	event := rmqevents.NewOrderCreatedEvent(orderId)
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to prepare event to publish: %w", err)
	}

	err = p.client.Channel().ExchangeDeclare(
		rmqevents.OrderCreatedEventExchange,
		amqp.ExchangeFanout,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	err = p.client.Channel().Publish(
		rmqevents.OrderCreatedEventExchange,
		defaultRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}
