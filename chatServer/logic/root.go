package logic

import (
	"context"
	"sync"

	"github.com/SzymonMielecki/chatApp/chatServer/persistance"
	"github.com/SzymonMielecki/chatApp/streaming/consumer"
	"github.com/SzymonMielecki/chatApp/types"
)

type Server struct {
	db       *persistance.DB
	consumer *consumer.StreamingConsumer
}

func NewServer(db *persistance.DB, consumer *consumer.StreamingConsumer) *Server {
	return &Server{db: db, consumer: consumer}
}

func (s *Server) Close() {
	s.consumer.Close()
}

func (s *Server) UploadMessages(ctx context.Context) {
	ch := make(chan *types.Message)
	defer close(ch)
	var wg sync.WaitGroup
	wg.Add(1)
	go s.consumer.ReceiveMessages(ctx, ch, &wg)
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return
		case msg := <-ch:
			s.db.CreateMessage(msg)
		}
	}
}
