package config

import (
	"fmt"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	Port          string `envconfig:"PORT,default=:9000"`
	LogLevel      string `envconfig:"LOG_LEVEL,default=INFO"`
	RMQProducer   *RMQProducer
	KafkaConsumer *KafkaConsumer
}

func Init() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Init(cfg); err != nil {
		return nil, fmt.Errorf("failed to init config: %w", err)
	}

	return cfg, nil
}
