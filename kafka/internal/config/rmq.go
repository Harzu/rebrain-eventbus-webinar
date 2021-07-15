package config

import "time"

type RMQProducer struct {
	ConnectionURL    string        `envconfig:"RMQ_PRODUCER_CONNECTION_URL"`
	ReconnectTimeout time.Duration `envconfig:"RMQ_PRODUCER_CONNECTION_TIMEOUT,default=1m"`
	OrderCreated     RMQQueue
}

type RMQQueue struct {
	Name       string
	RoutingKey string `envconfig:"optional"`
	Exchange   string
}
