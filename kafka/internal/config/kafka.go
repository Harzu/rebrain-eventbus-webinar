package config

type KafkaConsumer struct {
	Brokers       []string `envconfig:"KAFKA_BROKERS"`
	ConsumerGroup string   `envconfig:"KAFKA_CONSUMER_GROUP"`
	Version       string   `envconfig:"KAFKA_VERSION"`
	Topics        *KafkaTopics
}

type KafkaTopics struct {
	PaymentSuccess string `envconfig:"KAFKA_TOPIC_PAYMENT_SUCCESS"`
}

func (t KafkaTopics) Slice() []string {
	return []string{
		t.PaymentSuccess,
	}
}
