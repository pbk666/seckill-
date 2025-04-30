package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

var writer *kafka.Writer

func InitProducer(brokers []string, topic string) {
	writer = &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
	}
}

func SendMessage(key string, value []byte) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
	}
	err := writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Kafka 写入失败: %v\n", err)
	}
	return err
}
