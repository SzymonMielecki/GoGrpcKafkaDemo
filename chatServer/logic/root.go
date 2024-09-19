package logic

import (
	"context"

	"github.com/SzymonMielecki/chatApp/chatServer/persistance"
	"github.com/SzymonMielecki/chatApp/chatServer/streaming"
	pb "github.com/SzymonMielecki/chatApp/chatService"
	"github.com/SzymonMielecki/chatApp/types"
)

type Server struct {
	pb.UnimplementedChatServiceServer
	db        *persistance.DB
	streaming *streaming.Streaming
}

func NewServer(db *persistance.DB, streaming *streaming.Streaming) *Server {
	return &Server{db: db, streaming: streaming}
}

func (s *Server) Close() {
	s.streaming.Close()
}

func (s *Server) SendMessage(ctx context.Context, in *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	msg, err := s.db.CreateMessage(
		&types.Message{
			SenderID: uint(in.SenderId),
			Content:  in.Content,
		},
	)
	if err != nil {
		return &pb.SendMessageResponse{
			Success:      false,
			MessageId:    0,
			ErrorMessage: "Sending message failed",
		}, err
	}

	s.streaming.SendMessage(msg)

	return &pb.SendMessageResponse{
		Success:      true,
		MessageId:    uint32(msg.ID),
		ErrorMessage: "",
	}, nil
}
