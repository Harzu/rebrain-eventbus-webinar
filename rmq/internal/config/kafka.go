package config

type KafkaProducer struct {
	Brokers []string `envconfig:"KAFKA_BROKERS"`
}
