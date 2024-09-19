package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/SzymonMielecki/chatApp/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type StreamingProducer struct {
	producer       *kafka.Producer
	topicPartition *kafka.TopicPartition
}

func NewStreamingProducer(ctx context.Context, topic string, partition int, brokers []string) (*StreamingProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(brokers, ","),
		"group.id":          "chat",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create producer in streaming/root.go: \n%v", err)
	}
	topicPartition := kafka.TopicPartition{Topic: &topic, Partition: int32(partition)}
	return &StreamingProducer{producer: producer, topicPartition: &topicPartition}, nil
}

func (s *StreamingProducer) Close() {
	s.producer.Close()
}

func (s *StreamingProducer) SendMessage(ctx context.Context, message *types.Message) error {
	json, err := json.Marshal(message)
	if err != nil {
		return err
	}
	s.producer.Produce(&kafka.Message{
		TopicPartition: *s.topicPartition,
		Key:            []byte(strconv.FormatUint(uint64(message.ID), 10)),
		Value:          json,
	}, nil)
	return err
}
