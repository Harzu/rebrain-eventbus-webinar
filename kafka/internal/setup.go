package internal

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/config"
	kafkahandler "github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/handlers/kafka"
	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/services"
	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/system/eventbus/kafka"
	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/system/eventbus/rabbitmq"
	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/system/metrics"
)

type Application struct {
	config        *config.Config
	logger        *zap.Logger
	httpServer    *http.Server
	metricsClient *metrics.Client
	consumer      *kafka.Consumer
	rmqClient     *rabbitmq.Client
}

func NewApplication(cfg *config.Config, logger *zap.Logger) (*Application, error) {
	rmqClient, err := rabbitmq.NewClient(cfg.RMQProducer, logger)
	if err != nil {
		return nil, err
	}

	serviceContainer := services.New(rmqClient)
	metricsClient := metrics.New()
	handler := kafkahandler.NewHandler(cfg.KafkaConsumer.Topics, serviceContainer, metricsClient)
	consumer, err := kafka.NewConsumer(cfg.KafkaConsumer, logger, handler, metricsClient)
	if err != nil {
		return nil, err
	}

	httpServer := &http.Server{Addr: fmt.Sprintf(":%s", cfg.Port)}

	return &Application{
		config:        cfg,
		logger:        logger,
		httpServer:    httpServer,
		metricsClient: metricsClient,
		consumer:      consumer,
		rmqClient:     rmqClient,
	}, nil
}

func (a *Application) Run() error {
	a.consumer.Consume()

	go func() {
		http.Handle("/metrics", a.metricsClient.Handler())
		a.logger.Info("run http serve", zap.String("http_port", a.config.Port))
		if err := a.httpServer.ListenAndServe(); err != nil {
			a.logger.Error("http serve error", zap.Error(err))
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	return a.shutdown()
}

func (a *Application) shutdown() error {
	return multierr.Combine(
		a.consumer.Close(),
		a.rmqClient.Close(),
		a.httpServer.Close(),
	)
}
