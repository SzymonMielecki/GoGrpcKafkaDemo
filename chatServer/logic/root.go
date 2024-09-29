package logic

import (
	"context"
	"sync"

	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/chatServer/persistance"
	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/streaming/client"
	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/types"
)

type Server struct {
	db     *persistance.DB
	client *client.StreamingClient
}

func NewServer(db *persistance.DB, client *client.StreamingClient) *Server {
	return &Server{db: db, client: client}
}

func (s *Server) Close() {
	s.client.Close()
}

func (s *Server) UploadMessages(ctx context.Context) {
	ch := make(chan *types.Message)
	defer close(ch)
	var wg sync.WaitGroup
	wg.Add(1)
	go s.client.ReceiveMessages(ctx, ch, &wg)
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return
		case msg := <-ch:
			s.db.CreateMessage(&types.Message{
				Content:  msg.Content,
				SenderID: msg.SenderID,
			})
		}
	}
}
