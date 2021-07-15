package kafka

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/config"
	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/system/contextlogger"
	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/system/metrics"
)

type Consumer struct {
	ctx           context.Context
	cancel        context.CancelFunc
	client        sarama.Client
	consumerGroup sarama.ConsumerGroup
	topics        *config.KafkaTopics
	metrics       *metrics.Client
	handler       sarama.ConsumerGroupHandler
}

func NewConsumer(
	cfg *config.KafkaConsumer,
	logger *zap.Logger,
	handler sarama.ConsumerGroupHandler,
	metrics *metrics.Client,
) (*Consumer, error) {
	version, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		return nil, fmt.Errorf("parsing Kafka version error: %w", err)
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = version

	if cfg.SASLEnabled {
		saramaConfig.Net.SASL.User = cfg.User
		saramaConfig.Net.SASL.Password = cfg.Password
		saramaConfig.Net.SASL.Enable = true
	}

	logger.Info("Connecting to Kafka...")
	client, err := sarama.NewClient(cfg.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("make Kafka client error: %w", err)
	}

	logger.Info("Successfully connected to Kafka")
	consumerGroup, err := sarama.NewConsumerGroupFromClient(cfg.ConsumerGroup, client)
	if err != nil {
		return nil, fmt.Errorf("make Kafka consumer group error: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctx = contextlogger.Enrich(ctx, logger)

	return &Consumer{
		ctx:           ctx,
		cancel:        cancel,
		client:        client,
		consumerGroup: consumerGroup,
		topics:        cfg.Topics,
		metrics:       metrics,
		handler:       handler,
	}, nil
}

func (c *Consumer) Consume() {
	logger := contextlogger.Fetch(c.ctx)
	go func() {
		for c.ctx.Err() == nil {
			if err := c.consumerGroup.Consume(c.ctx, c.topics.Slice(), c.handler); err != nil {
				logger.Error("kafka consumer error", zap.Error(err))
			}
		}
	}()
}

func (c *Consumer) Close() error {
	c.cancel()

	return multierr.Combine(
		c.consumerGroup.Close(),
		c.client.Close(),
	)
}
