package services

import (
	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/services/rabbitmq"
	rmqclient "github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/system/eventbus/rabbitmq"
)

type Container struct {
	RMQProducer rabbitmq.Producer
}

func New(rmqClient *rmqclient.Client) *Container {
	return &Container{
		RMQProducer: rabbitmq.NewProducer(rmqClient),
	}
}
