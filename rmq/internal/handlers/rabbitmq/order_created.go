package rabbitmq

import (
	"context"
	"time"

	"github.com/Harzu/rebrain-eventbus-webinar/kafka/pkg/rmqevents"
	"go.uber.org/zap"

	"github.com/Harzu/rebrain-eventbus-webinar/rmq/internal/services/kafka"
	"github.com/Harzu/rebrain-eventbus-webinar/rmq/internal/system/contextlogger"
)

type orderCreatedHandler struct {
	producer kafka.Client
}

func newOrderCreatedHandler(producer kafka.Client) queueHandler {
	return &orderCreatedHandler{
		producer: producer,
	}
}

func (h *orderCreatedHandler) Message() interface{} {
	return &rmqevents.OrderCreatedEvent{}
}

func (h *orderCreatedHandler) Handle(ctx context.Context, msg interface{}) error {
	m, ok := msg.(*rmqevents.OrderCreatedEvent)
	if !ok {
		return InvalidMessageErr
	}

	contextlogger.Fetch(ctx).Info("handle message", zap.Any("message", m))
	time.Sleep(2 * time.Second)
	return h.producer.PaymentSuccessProduce(ctx, 1, 300)
}
