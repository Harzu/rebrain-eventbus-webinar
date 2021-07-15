package services

import (
	"go.uber.org/zap"

	"github.com/Harzu/rebrain-eventbus-webinar/rmq/internal/config"
	"github.com/Harzu/rebrain-eventbus-webinar/rmq/internal/services/kafka"
)

type Container struct {
	KafkaProducer *kafka.Producer
}

func New(cfg *config.Config, logger *zap.Logger) (*Container, error) {
	kafkaProducer, err := kafka.NewProducer(cfg.KafkaProducer, logger)
	if err != nil {
		return nil, err
	}

	return &Container{
		KafkaProducer: kafkaProducer,
	}, nil
}

func (c *Container) Close() error {
	return c.KafkaProducer.Close()
}
