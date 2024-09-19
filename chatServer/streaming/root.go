package streaming

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/SzymonMielecki/chatApp/types"
	"github.com/segmentio/kafka-go"
)

type Streaming struct {
	conn *kafka.Conn
}

func NewStreaming(topic string, partition int) *Streaming {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	return &Streaming{conn: conn}
}
func (s *Streaming) Close() {
	s.conn.Close()
}

func (s *Streaming) SendMessage(message *types.Message) error {
	s.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// Serialize the message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = s.conn.WriteMessages(
		kafka.Message{
			Value: messageBytes,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
