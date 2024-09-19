package logic

import (
	"context"

	"github.com/SzymonMielecki/chatApp/types"
	"github.com/SzymonMielecki/chatApp/usersServer/persistance"
	pb "github.com/SzymonMielecki/chatApp/usersService"
)

type Server struct {
	pb.UnimplementedUsersServiceServer
	db *persistance.DB
}

func NewServer(db *persistance.DB) *Server {
	return &Server{db: db}
}

func (s *Server) LoginUser(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	found_username, username_err := s.db.GetUserByUsername(in.Username)
	found_email, email_err := s.db.GetUserByEmail(in.Email)
	if username_err != nil && email_err != nil {
		return &pb.LoginUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "User not found",
		}, nil
	}
	if found_username != nil {
		return &pb.LoginUserResponse{
			Success: true,
			User:    found_username.ToProto(),
			Message: "Logged in as " + found_username.Username,
		}, nil
	}
	if found_email != nil {
		return &pb.LoginUserResponse{
			Success: true,
			User:    found_email.ToProto(),
			Message: "Logged in as " + found_email.Username,
		}, nil
	}
	return &pb.LoginUserResponse{
		Success: false,
		User:    &pb.User{},
		Message: "User not found",
	}, nil
}

func (s *Server) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	user_from_username, _ := s.db.GetUserByUsername(in.Username)
	if user_from_username != nil {
		return &pb.RegisterUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "Username already exists",
		}, nil
	}
	user_from_email, _ := s.db.GetUserByEmail(in.Email)
	if user_from_email != nil {
		return &pb.RegisterUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "Email already exists",
		}, nil
	}
	user, err := s.db.CreateUser(&types.User{
		Username:     in.Username,
		Email:        in.Email,
		PasswordHash: in.PasswordHash,
	})
	if err != nil {
		return &pb.RegisterUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "Error creating user, " + err.Error(),
		}, nil
	}
	return &pb.RegisterUserResponse{
		Success: true,
		User:    user.ToProto(),
		Message: "Registered",
	}, nil
}

func (s *Server) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.db.GetUserById(uint(in.Id))
	if err != nil {
		return &pb.GetUserResponse{
			Success: false,
			User:    nil,
			Message: "User not found",
		}, err
	}
	return &pb.GetUserResponse{
		Success: true,
		User:    user.ToProto(),
		Message: "User found",
	}, nil
}

func (s *Server) CheckUser(ctx context.Context, in *pb.CheckUserRequest) (*pb.CheckUserResponse, error) {
	user, err := s.db.GetUserByUsernameAndEmail(in.Username, in.Email)
	if err != nil {
		return &pb.CheckUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "User not found, error: " + err.Error(),
		}, nil
	}
	if user.PasswordHash != in.PasswordHash {
		return &pb.CheckUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "Incorrect password",
		}, nil
	}
	return &pb.CheckUserResponse{
		Success: true,
		User:    user.ToProto(),
		Message: "User found",
	}, nil
}
