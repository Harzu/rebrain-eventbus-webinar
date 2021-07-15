package rabbitmq

import (
	"errors"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/config"
)

type Client struct {
	logger  *zap.Logger
	config  *config.RMQProducer
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewClient(cfg *config.RMQProducer, logger *zap.Logger) (*Client, error) {
	conn, channel, err := setupConnection(cfg.ConnectionURL)
	if err != nil {
		return nil, err
	}

	logger.Info("rabbitmq producer connect successful")

	c := &Client{
		logger:  logger,
		config:  cfg,
		conn:    conn,
		channel: channel,
	}

	c.reconnectHandler()
	return c, nil
}

func (c *Client) Channel() *amqp.Channel {
	return c.channel
}

func (c *Client) Close() error {
	err := multierr.Append(c.channel.Close(), c.conn.Close())
	if errors.Is(err, amqp.ErrClosed) {
		return nil
	}

	return err
}

func (c *Client) reconnectHandler() {
	go func() {
		closeErr := <-c.channel.NotifyClose(make(chan *amqp.Error))
		if closeErr == nil {
			c.logger.Info("graceful close")
			return
		}

		linearRetry := NewLinearBackoff(c.config.ReconnectTimeout)

		for {
			conn, channel, err := setupConnection(c.config.ConnectionURL)
			if err != nil {
				c.logger.Error("failed to reconnect")
				time.Sleep(linearRetry.NextBackoff())
				continue
			}

			c.conn = conn
			c.channel = channel

			c.logger.Info("connect successful")
			break
		}

		c.reconnectHandler()
	}()
}

func setupConnection(connectionURL string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(connectionURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial rabbitmq connect: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to init rabbitmq channel: %w", err)
	}

	return conn, channel, nil
}
