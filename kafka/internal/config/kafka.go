package config

type KafkaConsumer struct {
	Brokers       []string `envconfig:"KAFKA_BROKERS"`
	ConsumerGroup string   `envconfig:"KAFKA_CONSUMER_GROUP"`
	SASLEnabled   bool     `envconfig:"KAFKA_SASLE_ENABLE,default=false"`
	User          string   `envconfig:"KAFKA_USER,optional"`
	Password      string   `envconfig:"KAFKA_PASSWORD,optional"`
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
