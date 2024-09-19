package logic

import (
	"context"
	"sync"

	"github.com/SzymonMielecki/chatApp/chatServer/persistance"
	"github.com/SzymonMielecki/chatApp/streaming"
)

type Server struct {
	db        *persistance.DB
	streaming *streaming.Streaming
}

func NewServer(db *persistance.DB, streaming *streaming.Streaming) *Server {
	return &Server{db: db, streaming: streaming}
}

func (s *Server) Close() {
	s.streaming.Close()
}

func (s *Server) UploadMessages(ctx context.Context, wg *sync.WaitGroup) {
	s.streaming.UploadMessages(ctx, s.db, wg)
}
