package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/Harzu/rebrain-eventbus-webinar/rmq/internal/config"
	"github.com/Harzu/rebrain-eventbus-webinar/rmq/internal/handlers/rabbitmq"
	"github.com/Harzu/rebrain-eventbus-webinar/rmq/internal/system/contextlogger"
)

type Consumer struct {
	ctx     context.Context
	cancel  context.CancelFunc
	logger  *zap.Logger
	conn    *amqp.Connection
	ch      *amqp.Channel
	config  *config.RMQConsumer
	handler rabbitmq.Handler
	queues  []config.RMQQueue
}

func NewConsumer(cfg *config.RMQConsumer, handler rabbitmq.Handler, logger *zap.Logger) (*Consumer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = contextlogger.Enrich(ctx, logger)

	return &Consumer{
		ctx:     ctx,
		cancel:  cancel,
		logger:  logger,
		config:  cfg,
		handler: handler,
		queues: []config.RMQQueue{
			cfg.OrderCreated,
		},
	}, nil
}

func (c *Consumer) Consume() {
	go func() {
		for c.ctx.Err() == nil {
			if err := c.setupConnection(); err != nil {
				c.logger.Error("failed to connect to RabbitMQ", zap.Error(err))
				// waiting one second to reconnect
				time.Sleep(c.config.ReconnectTimeout)
				continue
			}
			c.logger.Info("successfully connected to RabbitMQ")

			wg := &sync.WaitGroup{}
			for _, q := range c.queues {
				if q.Name == "" {
					continue
				}

				err := c.ch.ExchangeDeclare(q.Exchange, amqp.ExchangeFanout, true, false, false, false, nil)
				if err != nil {
					c.logger.Error("failed to declare a exchange", zap.Error(err))
					continue
				}

				queue, err := c.ch.QueueDeclare(q.Name, true, false, false, false, nil)
				if err != nil {
					c.logger.Error("failed to declare a queue", zap.Error(err))
					continue
				}

				if err := c.ch.QueueBind(q.Name, q.RoutingKey, q.Exchange, false, nil); err != nil {
					c.logger.Error("failed to bind a queue", zap.Error(err))
					continue
				}

				if err := c.ch.Qos(1, 0, false); err != nil {
					c.logger.Error("failed to configure Qos", zap.Error(err))
					continue
				}

				messages, err := c.ch.Consume(queue.Name, "", false, false, false, false, nil)
				if err != nil {
					c.logger.Error("failed to register a consumer", zap.Error(err))
					continue
				}

				wg.Add(1)
				go func() {
					for message := range messages {
						c.handler.Handle(c.ctx, queue.Name, message)
						if err := message.Ack(false); err != nil {
							c.logger.Error("failed to acknowledge a message", zap.Error(err))
						}
					}
					wg.Done()
				}()
			}

			wg.Wait()
			c.logger.Info("RabbitMQ connection was closed")
		}
	}()
}

func (c *Consumer) setupConnection() error {
	conn, err := amqp.Dial(c.config.ConnectionURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a RabbitMQ channel: %w", err)
	}

	c.conn = conn
	c.ch = ch

	return nil
}

func (c *Consumer) Close() error {
	c.cancel()

	return multierr.Append(
		c.ch.Close(),
		c.conn.Close(),
	)
}
