package kafka

import (
	"context"
	"time"

	"github.com/Harzu/rebrain-eventbus-webinar/rmq/pkg/kafkaevents"
	"go.uber.org/zap"

	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/services/rabbitmq"
	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/system/contextlogger"
)

type paymentSuccessHandler struct {
	producer rabbitmq.Producer
}

func newPaymentSuccessHandler(producer rabbitmq.Producer) topicHandler {
	return &paymentSuccessHandler{
		producer: producer,
	}
}

func (h *paymentSuccessHandler) Message() interface{} {
	return &kafkaevents.PaymentSuccessEvent{}
}

func (h *paymentSuccessHandler) Handle(ctx context.Context, msg interface{}) error {
	m, ok := msg.(*kafkaevents.PaymentSuccessEvent)
	if !ok {
		return InvalidMessageErr
	}

	contextlogger.Fetch(ctx).Info("handle message", zap.Any("message", m))
	time.Sleep(2 * time.Second)
	return h.producer.OrderCreatedProduce(1)
}
