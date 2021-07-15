package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"

	"github.com/Harzu/rebrain-eventbus-webinar/rmq/internal/config"
	"github.com/Harzu/rebrain-eventbus-webinar/rmq/internal/system/contextlogger"
	"github.com/Harzu/rebrain-eventbus-webinar/rmq/pkg/kafkaevents"
)

type Client interface {
	PaymentSuccessProduce(ctx context.Context, userId, amount uint64) error
}

type Producer struct {
	saramaProducer sarama.SyncProducer
}

func NewProducer(cfg *config.KafkaProducer, logger *zap.Logger) (*Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForLocal
	saramaConfig.Producer.Retry.Max = 10
	saramaConfig.Producer.Return.Successes = true

	saramaProducer, err := sarama.NewSyncProducer(cfg.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama producer: %w", err)
	}

	logger.Info("kafka producer created")

	return &Producer{
		saramaProducer: saramaProducer,
	}, nil
}

func (p *Producer) PaymentSuccessProduce(ctx context.Context, userId, amount uint64) error {
	payload := kafkaevents.NewPaymentSuccessEvent(userId, amount)
	bytePayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to create message payload")
	}

	event := &sarama.ProducerMessage{
		Topic: kafkaevents.PaymentSuccessTopic,
		Value: sarama.ByteEncoder(bytePayload),
	}

	partition, offset, err := p.saramaProducer.SendMessage(event)
	if err != nil {
		return fmt.Errorf("failed to produce payment success event: %w", err)
	}

	contextlogger.Fetch(ctx).Info(
		"message produced",
		zap.String("topic", "topic"),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
	)

	return nil
}

func (p *Producer) Close() error {
	return p.saramaProducer.Close()
}
