package consumer

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type StreamingConsumer struct {
	c     *kafka.Consumer
	topic string
}

func NewStreamingConsumer(ctx context.Context, topic string, broker string) (*StreamingConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest"})
	if err != nil {
		return nil, err
	}
	return &StreamingConsumer{c: c, topic: topic}, nil
}

func (s *StreamingConsumer) Close() {
	s.c.Close()
}
func (s *StreamingConsumer) ReceiveMessages(ctx context.Context, ch chan<- *types.Message, wg *sync.WaitGroup) {
	run := true

	for run {
		select {
		case <-ctx.Done():
			wg.Done()
			run = false
		default:
			ev, err := s.c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				continue
			}
			fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))

			senderId, err := strconv.ParseUint(string(ev.Key), 10, 64)
			if err != nil {
				fmt.Printf("Error parsing key: %v\n", err)
				continue
			}
			content := string(ev.Value)

			ch <- &types.Message{
				SenderID: uint(senderId),
				Content:  content,
			}
		}
	}
}
