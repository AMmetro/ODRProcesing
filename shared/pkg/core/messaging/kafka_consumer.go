package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

// KafkaConsumer интерфейс для получения сообщений из Kafka
type KafkaConsumer interface {
	Start(ctx context.Context) error
	Close() error
	SetMessageHandler(handler func(message []byte) error)
}

// kafkaConsumerImpl реализация Kafka consumer
type kafkaConsumerImpl struct {
	reader         *kafka.Reader
	messageHandler func(message []byte) error
	stopCh         chan struct{}
}

// NewKafkaConsumer создает новый Kafka consumer
func NewKafkaConsumer(brokers []string, groupID string, topic string) (KafkaConsumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &kafkaConsumerImpl{
		reader: reader,
		stopCh: make(chan struct{}),
	}, nil
}

// SetMessageHandler устанавливает обработчик сообщений
func (c *kafkaConsumerImpl) SetMessageHandler(handler func(message []byte) error) {
	c.messageHandler = handler
}

// Start запускает consumer и начинает слушать сообщения
func (c *kafkaConsumerImpl) Start(ctx context.Context) error {
	if c.messageHandler == nil {
		return fmt.Errorf("message handler not set")
	}

	go func() {
		for {
			select {
			case <-c.stopCh:
				return
			default:
				msg, err := c.reader.FetchMessage(ctx)
				if err != nil {
					log.Printf("error reading message: %v", err)
					continue
				}

				log.Println("MESSAGE RECEIVED")
				log.Println(string(msg.Value))

				// Обработка сообщения
				if err := c.messageHandler(msg.Value); err != nil {
					log.Printf("error handling message: %v", err)
					continue
				}

				// Коммитим смещение после успешной обработки
				if err := c.reader.CommitMessages(ctx, msg); err != nil {
					log.Printf("error committing message: %v", err)
				}
			}
		}
	}()

	return nil
}

// Close закрывает connection к Kafka
func (c *kafkaConsumerImpl) Close() error {
	close(c.stopCh)
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}

// UnmarshalMessage десериализует JSON сообщение в структуру
func UnmarshalMessage(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
