package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/SzymonMielecki/chatApp/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type StreamingConsumer struct {
	consumer  *kafka.Consumer
	topic     string
	partition int
}

func NewStreamingConsumer(ctx context.Context, topic string, partition int, brokers []string) (*StreamingConsumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(brokers, ","),
		"group.id":          "chat",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create client in streaming/root.go: \n%v", err)
	}
	return &StreamingConsumer{consumer: consumer, topic: topic, partition: partition}, nil
}

func (s *StreamingConsumer) Close() error {
	return s.consumer.Close()
}
func (s *StreamingConsumer) ReceiveMessages(ctx context.Context, ch chan<- *types.Message, wg *sync.WaitGroup) {
	err := s.consumer.Subscribe(s.topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
	}

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
			msg, err := s.consumer.ReadMessage(time.Second)
			if err == nil {
				var message types.Message
				json.Unmarshal(msg.Value, &message)
				ch <- &message
			}
		}
	}
}
