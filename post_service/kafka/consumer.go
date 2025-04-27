package kafka

import (
	"github.com/IBM/sarama"
	"github.com/joy095/post-service/models"
)

type KafkaConsumer struct {
	consumer sarama.Consumer
	topic    string
}

func NewKafkaConsumer(config *KafkaConfig) (*KafkaConsumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_2_1_0
	consumer, err := sarama.NewConsumer(config.Brokers, saramaConfig)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{consumer: consumer, topic: config.Topic}, nil
}

func (c *KafkaConsumer) Consume() ([]*models.Message, error) {
	partitionConsumer, err := c.consumer.ConsumePartition(c.topic, 0, sarama.OffsetOldest)
	if err != nil {
		return nil, err
	}
	defer partitionConsumer.Close()

	messages := make([]*models.Message, 0)
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			messages = append(messages, &models.Message{Value: string(msg.Value)})
		case err := <-partitionConsumer.Errors():
			return nil, err
		}
	}
}
