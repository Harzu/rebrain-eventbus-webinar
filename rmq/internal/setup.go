package internal

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/Harzu/rebrain-eventbus-webinar/rmq/internal/config"
	rmqHandler "github.com/Harzu/rebrain-eventbus-webinar/rmq/internal/handlers/rabbitmq"
	"github.com/Harzu/rebrain-eventbus-webinar/rmq/internal/services"
	"github.com/Harzu/rebrain-eventbus-webinar/rmq/internal/system/eventbus/rabbitmq"
	"github.com/Harzu/rebrain-eventbus-webinar/rmq/internal/system/metrics"
)

type Application struct {
	config        *config.Config
	logger        *zap.Logger
	httpServer    *http.Server
	metricsClient *metrics.Client
	rmqConsumer   *rabbitmq.Consumer
}

func NewApplication(cfg *config.Config, logger *zap.Logger) (*Application, error) {
	serviceContainer, err := services.New(cfg, logger)
	if err != nil {
		return nil, err
	}

	metricsClient := metrics.New()
	handler := rmqHandler.NewHandler(cfg.RMQConsumer.Queues, serviceContainer, metricsClient)
	consumer, err := rabbitmq.NewConsumer(cfg.RMQConsumer, handler, logger)
	if err != nil {
		return nil, err
	}

	httpServer := &http.Server{Addr: fmt.Sprintf(":%s", cfg.Port)}

	return &Application{
		config:        cfg,
		logger:        logger,
		httpServer:    httpServer,
		metricsClient: metricsClient,
		rmqConsumer:   consumer,
	}, nil
}

func (a *Application) Run() error {
	a.rmqConsumer.Consume()

	go func() {
		http.Handle("/metrics", a.metricsClient.Handler())
		a.logger.Info("run http serve", zap.String("http_port", a.config.Port))
		if err := a.httpServer.ListenAndServe(); err != nil {
			a.logger.Error("http serve error", zap.Error(err))
		}
	}()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	return a.shutdown()
}

func (a *Application) shutdown() error {
	return multierr.Combine(
		a.httpServer.Close(),
		a.rmqConsumer.Close(),
	)
}
