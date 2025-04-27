package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/joy095/post-service/kafka"
	"github.com/joy095/post-service/models"
)

func ProduceMessage(c *gin.Context) {
	var message models.Message
	err := c.BindJSON(&message)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	config := kafka.NewKafkaConfig([]string{"localhost:9092"}, "my_topic")
	producer, err := kafka.NewKafkaProducer(config)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = producer.Produce(&message)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Sent to Kafka"})
}

func ConsumeMessages(c *gin.Context) {
	config := kafka.NewKafkaConfig([]string{"localhost:9092"}, "my_topic")
	consumer, err := kafka.NewKafkaConsumer(config)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	messages, err := consumer.Consume()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, messages)
}
