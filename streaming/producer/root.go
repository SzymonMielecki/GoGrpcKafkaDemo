package producer

import (
	"context"
	"fmt"
	"sync"

	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type StreamingProducer struct {
	p     *kafka.Producer
	topic string
}

func NewStreamingProducer(ctx context.Context, topic string, broker string) (*StreamingProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"acks":              "all",
	})
	if err != nil {
		return nil, err
	}
	return &StreamingProducer{p: p, topic: topic}, nil
}

func (s *StreamingProducer) Init(ctx context.Context) {
	go func() {
		for e := range s.p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()
}

func (s *StreamingProducer) Close() {
	s.p.Close()
}

func (s *StreamingProducer) Produce(ctx context.Context, message *types.Message, wg *sync.WaitGroup) error {
	s.p.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &s.topic, Partition: kafka.PartitionAny},
			Key:            []byte(fmt.Sprintf("%d", message.SenderID)),
			Value:          []byte(message.Content),
		}, nil)
	return nil
}
