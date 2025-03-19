package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

// KafkaProducer wraps the Sarama producer
type KafkaProducer struct {
	Producer sarama.SyncProducer
	Topic    string
}

// KafkaConsumer wraps the Sarama consumer
type KafkaConsumer struct {
	Consumer sarama.ConsumerGroup
	Topic    string
	Ready    chan bool
}

// ConsumerGroupHandler implements sarama.ConsumerGroupHandler
type ConsumerGroupHandler struct {
	ready chan bool
	wg    sync.WaitGroup
}

// Setup is run at the beginning of a new session
func (h *ConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

// Cleanup is run at the end of a session
func (h *ConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes messages from a partition
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	h.wg.Add(1)
	defer h.wg.Done()

	for message := range claim.Messages() {
		log.Printf("Message topic:%s partition:%d offset:%d value:%s\n",
			message.Topic, message.Partition, message.Offset, string(message.Value))
		session.MarkMessage(message, "")
	}

	return nil
}

// NewKafkaProducer creates a new Kafka producer using Sarama
func NewKafkaProducer(brokers, topic string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Retry.Backoff = 500 * time.Millisecond
	config.Producer.Return.Successes = true

	// Wait for Kafka to be ready
	for i := 0; i < 10; i++ {
		producer, err := sarama.NewSyncProducer(strings.Split(brokers, ","), config)
		if err == nil {
			return &KafkaProducer{
				Producer: producer,
				Topic:    topic,
			}, nil
		}
		log.Printf("Failed to create producer, retrying in 5 seconds: %v", err)
		time.Sleep(5 * time.Second)
	}

	// Try one last time and return any error
	producer, err := sarama.NewSyncProducer(strings.Split(brokers, ","), config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		Producer: producer,
		Topic:    topic,
	}, nil
}

// NewKafkaConsumer creates a new Kafka consumer using Sarama
func NewKafkaConsumer(brokers, topic, group string) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Version = sarama.V2_4_0_0 // Use version compatible with your Kafka server

	// Wait for Kafka to be ready
	for i := 0; i < 10; i++ {
		consumer, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, config)
		if err == nil {
			return &KafkaConsumer{
				Consumer: consumer,
				Topic:    topic,
				Ready:    make(chan bool),
			}, nil
		}
		log.Printf("Failed to create consumer, retrying in 5 seconds: %v", err)
		time.Sleep(5 * time.Second)
	}

	// Try one last time and return any error
	consumer, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		Consumer: consumer,
		Topic:    topic,
		Ready:    make(chan bool),
	}, nil
}

// Produce sends a message to Kafka
func (kp *KafkaProducer) Produce(message string) error {
	msg := &sarama.ProducerMessage{
		Topic: kp.Topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err := kp.Producer.SendMessage(msg)
	return err
}

// Close closes the producer
func (kp *KafkaProducer) Close() error {
	return kp.Producer.Close()
}

// Consume starts consuming messages
func (kc *KafkaConsumer) Consume(ctx context.Context) {
	handler := &ConsumerGroupHandler{
		ready: kc.Ready,
	}

	go func() {
		for {
			if err := kc.Consumer.Consume(ctx, []string{kc.Topic}, handler); err != nil {
				log.Printf("Error from consumer: %v", err)
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-kc.Ready // Wait until the consumer is ready
}

// Close closes the consumer
func (kc *KafkaConsumer) Close() error {
	return kc.Consumer.Close()
}

func main() {
	// Kafka configuration
	brokers := getEnv("KAFKA_BROKERS", "kafka:9092")
	topic := getEnv("KAFKA_TOPIC", "gin-kafka-topic")
	group := getEnv("KAFKA_GROUP", "gin-consumer-group")

	log.Printf("Connecting to Kafka at %s", brokers)

	// Create Kafka producer
	producer, err := NewKafkaProducer(brokers, topic)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()
	log.Println("Kafka producer created successfully")

	// Create Kafka consumer
	consumer, err := NewKafkaConsumer(brokers, topic, group)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Start consumer in a goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer.Consume(ctx)
	log.Println("Kafka consumer ready")

	// Set up Gin router
	router := gin.Default()

	// Define API endpoints
	router.POST("/messages", func(c *gin.Context) {
		var message struct {
			Content string `json:"content" binding:"required"`
		}

		if err := c.ShouldBindJSON(&message); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Produce message to Kafka
		if err := producer.Produce(message.Content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to produce message"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Message sent to Kafka"})
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Service is running"})
	})

	// Start the Gin server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Shutdown the server with a timeout
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
