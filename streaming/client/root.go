package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/SzymonMielecki/GoGrpcKafkaDemo/types"
	"github.com/twmb/franz-go/pkg/kgo"
)

type StreamingClient struct {
	client    *kgo.Client
	partition int
}

func NewStreamingClient(ctx context.Context, topic string, partition int, brokers []string) (*StreamingClient, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ConsumeTopics(topic),
		kgo.ConsumerGroup("chat-app-consumer-group"),
	}
	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create client in streaming/root.go: \n%v", err)
	}
	return &StreamingClient{client: client}, nil
}

func (s *StreamingClient) Close() {
	s.client.Close()
}
func (s *StreamingClient) ReceiveMessages(ctx context.Context, ch chan<- *types.Message, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
			fetches := s.client.PollFetches(ctx)
			fetches.EachError(func(t string, p int32, err error) {
				log.Printf("fetch err topic %s partition %d: %v", t, p, err)
			})
			fetches.EachRecord(func(r *kgo.Record) {
				msg := &types.Message{}
				err := json.Unmarshal(r.Value, msg)
				if err != nil {
					log.Printf("error unmarshalling message in streaming/root.go: %v", err)
					return
				}
				ch <- msg
			})
		}
	}
}
