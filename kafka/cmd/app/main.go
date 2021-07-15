package main

import (
	"log"

	"go.uber.org/zap"

	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal"
	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/config"
	"github.com/Harzu/rebrain-eventbus-webinar/kafka/internal/system/logger"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalln("failed to init config", err)
	}

	appLogger, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatalln("failed to init logger", err)
	}

	app, err := internal.NewApplication(cfg, appLogger)
	if err != nil {
		appLogger.Fatal("failed to init app", zap.Error(err))
	}

	if err := app.Run(); err != nil {
		appLogger.Fatal("failed to run app", zap.Error(err))
	}
}
