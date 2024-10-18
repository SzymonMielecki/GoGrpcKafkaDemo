package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/SzymonMielecki/GoGrpcKafkaDemo/types"
	"github.com/twmb/franz-go/pkg/kgo"
)

type StreamingProducer struct {
	client *kgo.Client
}

func NewStreamingProducer(ctx context.Context, topic string, partition int, brokers []string) (*StreamingProducer, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.DefaultProduceTopic(topic),
	}
	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer in streaming/root.go: \n%v", err)
	}
	return &StreamingProducer{client: client}, nil
}

func (s *StreamingProducer) Close() {
	s.client.Close()
}

func (s *StreamingProducer) SendMessage(ctx context.Context, message *types.Message, wg *sync.WaitGroup) error {
	json, err := json.Marshal(message)
	if err != nil {
		return err
	}
	s.client.Produce(ctx, kgo.SliceRecord(json), func(r *kgo.Record, err error) {
		if err != nil {
			log.Printf("error producing record: %v", err)
		}
		wg.Done()
	})
	return err
}
