package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"

	"github.com/Harzu/rebrain-eventbus-webinar/rmq/internal/services"
	"github.com/Harzu/rebrain-eventbus-webinar/rmq/internal/system/contextlogger"
	"github.com/Harzu/rebrain-eventbus-webinar/rmq/internal/system/metrics"
)

const rmqLabelPrefix = "rmq_"

var HandlerNotRegisterErr = errors.New("handler not registered")

type queueHandler interface {
	Handle(ctx context.Context, message interface{}) error
	Message() interface{}
}

type Handler interface {
	Handle(ctx context.Context, queueName string, message amqp.Delivery)
}

type handler struct {
	handlers map[string]queueHandler
	metrics  *metrics.Client
}

func NewHandler(_ *services.Container, _ *zap.Logger, metrics *metrics.Client) Handler {
	h := &handler{
		metrics:  metrics,
		handlers: map[string]queueHandler{},
	}

	return h
}

func (h handler) Handle(ctx context.Context, queueName string, message amqp.Delivery) {
	logger := contextlogger.Fetch(ctx)
	msgLogger := logger.
		With(
			zap.String("queue_name", queueName),
			zap.String("decoded_value", string(message.Body)),
		)

	startTime := time.Now()
	defer func() {
		h.metrics.EventbusMessagesProcessingTime.Add(makeMessageType(queueName), time.Since(startTime).Seconds())
	}()

	h.metrics.EventbusMessageCount.AddClaimed(makeMessageType(queueName))

	handler, err := h.findHandler(queueName)
	if errors.Is(err, HandlerNotRegisterErr) {
		msgLogger.Info("skip message", zap.Error(err))
		h.metrics.EventbusMessageCount.AddSkipped(makeMessageType(queueName))
		return
	} else if err != nil {
		msgLogger.Error("Find message handler error", zap.Error(err))
		h.metrics.EventbusMessageCount.AddFailed(makeMessageType(queueName))
		return
	}

	handlerMsg := handler.Message()
	if err := json.Unmarshal(message.Body, &handlerMsg); err != nil {
		msgLogger.Error("parse handler message error", zap.Error(err))
		h.metrics.EventbusMessageCount.AddFailed(makeMessageType(queueName))
		return
	}

	handlerLogger := logger.With(zap.String("queue_name", queueName))
	if err := handler.Handle(contextlogger.Enrich(ctx, handlerLogger), handlerMsg); err != nil {
		handlerLogger.Error("message handle error", zap.Error(err))
		h.metrics.EventbusMessageCount.AddFailed(makeMessageType(queueName))
		return
	}

	h.metrics.EventbusMessageCount.AddSuccess(makeMessageType(queueName))
}

func (h handler) findHandler(queueName string) (queueHandler, error) {
	handler, ok := h.handlers[queueName]
	if ok {
		return handler, nil
	}

	return handler, nil
}

func makeMessageType(eventName string) string {
	return fmt.Sprintf("%s%s", rmqLabelPrefix, eventName)
}
