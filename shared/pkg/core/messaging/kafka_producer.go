package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer интерфейс для отправки сообщений в Kafka
type KafkaProducer interface {
	SendMessage(ctx context.Context, topic string, key string, message interface{}) error
	Close() error
}

// kafkaProducerImpl реализация Kafka producer
type kafkaProducerImpl struct {
	writer *kafka.Writer
}

// NewKafkaProducer создает новый Kafka producer
func NewKafkaProducer(brokers []string) (KafkaProducer, error) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll, // ждать подтверждения от всех реплик
		MaxAttempts:  3,                // повторные попытки
		BatchSize:    100,              // размер батча
		BatchTimeout: 1 * time.Second,  // таймаут батча
		WriteTimeout: 10 * time.Second, // таймаут на запись
	}

	return &kafkaProducerImpl{
		writer: writer,
	}, nil
}

// SendMessage отправляет сообщение в Kafka
func (p *kafkaProducerImpl) SendMessage(ctx context.Context, topic string, key string, message interface{}) error {
	// Сериализуем сообщение в JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	// Отправляем сообщение в Kafka
	err = p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: messageBytes,
	})

	if err != nil {
		return fmt.Errorf("write message to kafka: %w", err)
	}

	return nil
}

// Close закрывает connection к Kafka
func (p *kafkaProducerImpl) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}
