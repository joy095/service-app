package kafka

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

func NewKafkaConfig(brokers []string, topic string) *KafkaConfig {
	return &KafkaConfig{
		Brokers: brokers,
		Topic:   topic,
	}
}
