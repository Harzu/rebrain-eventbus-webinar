package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"runtime/debug"
	"time"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"

	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/config"
	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/services"
	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/system/contextlogger"
	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/system/metrics"
)

const kafkaLabelPrefix = "kfk_"

var (
	HandlerNotRegisteredErr = errors.New("handler not registered")
	InvalidMessageErr       = errors.New("invalid message")
)

type topicHandler interface {
	Handle(ctx context.Context, message interface{}) error
	Message() interface{}
}

type consumerGroupHandler struct {
	handlers map[string]topicHandler
	metrics  *metrics.Client
}

func NewHandler(
	topics *config.KafkaTopics,
	services *services.Container,
	metrics *metrics.Client,
) sarama.ConsumerGroupHandler {
	return &consumerGroupHandler{
		metrics: metrics,
		handlers: map[string]topicHandler{
			topics.PaymentSuccess: newPaymentSuccessHandler(services.RMQProducer),
		},
	}
}

func (h consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerGroupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for message := range claim.Messages() {
		session.MarkMessage(message, "")
		h.consumeMessage(session.Context(), message)
	}
	return nil
}

func (h consumerGroupHandler) consumeMessage(ctx context.Context, consumerMessage *sarama.ConsumerMessage) {
	logger := contextlogger.Fetch(ctx).
		With(zap.String("decoded_value", string(consumerMessage.Value)))

	startTime := time.Now()
	defer func() {
		if r := recover(); r != nil {
			logger.Error(
				"panic handle message",
				zap.Error(r.(error)),
				zap.Stack(string(debug.Stack())),
			)
		}
		h.metrics.EventbusMessagesProcessingTime.Add(makeMessageType(consumerMessage.Topic), time.Since(startTime).Seconds())
	}()

	logger.Info("message claimed")
	h.metrics.EventbusMessageCount.AddClaimed(makeMessageType(consumerMessage.Topic))
	h.metrics.KafkaMessageCurrentOffsetNumber.Set(consumerMessage.Topic, consumerMessage.Offset)

	handler, err := h.findHandler(consumerMessage)
	if errors.Is(err, HandlerNotRegisteredErr) {
		logger.Info("skip message", zap.Error(err))
		h.metrics.EventbusMessageCount.AddSkipped(makeMessageType(consumerMessage.Topic))
		return
	} else if err != nil {
		logger.Error("find message handler error", zap.Error(err))
		h.metrics.EventbusMessageCount.AddFailed(makeMessageType(consumerMessage.Topic))
		return
	}

	handlerMsg := handler.Message()
	if err := json.Unmarshal(consumerMessage.Value, &handlerMsg); err != nil {
		logger.Error("parse handler message error", zap.Error(err))
		h.metrics.EventbusMessageCount.AddFailed(makeMessageType(consumerMessage.Topic))
		return
	}

	if err := handler.Handle(ctx, handlerMsg); err != nil {
		logger.Error("message handle error", zap.Error(err))
		h.metrics.EventbusMessageCount.AddFailed(makeMessageType(consumerMessage.Topic))
		return
	}

	logger.Info("message successfully handled")
	h.metrics.EventbusMessageCount.AddSuccess(makeMessageType(consumerMessage.Topic))
}

func (h consumerGroupHandler) findHandler(consumerMessage *sarama.ConsumerMessage) (topicHandler, error) {
	handler, ok := h.handlers[consumerMessage.Topic]
	if !ok {
		return nil, HandlerNotRegisteredErr
	}
	return handler, nil
}

func makeMessageType(eventName string) string {
	return kafkaLabelPrefix + eventName
}
