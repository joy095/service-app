package kafka

import (
	"github.com/IBM/sarama"
	"github.com/joy095/post-service/models"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(config *KafkaConfig) (*KafkaProducer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{producer: producer, topic: config.Topic}, nil
}

func (p *KafkaProducer) Produce(message *models.Message) error {
	msg := &sarama.ProducerMessage{
		Topic: p.topic, // access the topic name from the struct field
		Value: sarama.StringEncoder(message.Value),
	}
	_, _, err := p.producer.SendMessage(msg)
	return err
}
