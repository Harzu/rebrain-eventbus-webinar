package config

import "time"

type RMQConsumer struct {
	ConnectionURL    string        `envconfig:"RMQ_CONSUMER_CONNECTION_URL"`
	ReconnectTimeout time.Duration `envconfig:"RMQ_CONSUMER_CONNECTION_TIMEOUT,default=1m"`
	Queues           *RMQQueues
}

type RMQQueues struct {
	OrderCreated queue
}

func (q *RMQQueues) Slice() []queue {
	return []queue{
		q.OrderCreated,
	}
}

type queue struct {
	Name       string
	RoutingKey string `envconfig:"optional"`
	Exchange   string
}
