package kafka

import (
	"context"
	"decleration/model"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

type OrderHandlerFunc func(ctx context.Context, msg model.OrderMessage) error

func ConsumeOrderMessages(handler OrderHandlerFunc) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "seckill-topic",
		GroupID: "order_group",
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka读取失败:", err)
			continue
		}

		// 假设你使用 JSON 传输
		var msg model.OrderMessage
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Println("消息反序列化失败:", err)
			continue
		}

		if err := handler(context.Background(), msg); err != nil {
			log.Println("处理订单失败:", err)
		}
	}
}
