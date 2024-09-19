package logic

import (
	"context"
	"fmt"

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
	var found types.User
	err := s.db.Find(&types.User{Username: in.UsernameOrEmail}).Or(&types.User{Email: in.UsernameOrEmail}).First(&found).Error
	if err != nil {
		return &pb.LoginUserResponse{}, err
	}
	if found.Model.ID == 0 {
		return &pb.LoginUserResponse{
			Success: false,
			Message: "User not found",
			User:    &pb.User{},
		}, nil
	}
	if found.PasswordHash != in.PasswordHash {
		return &pb.LoginUserResponse{
			Success: false,
			Message: "Incorrect password",
			User:    &pb.User{},
		}, nil
	}
	return &pb.LoginUserResponse{
		Success: true,
		Message: "Logged in",
		User:    found.ToProto(),
	}, nil
}

func (s *Server) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	username_exists := s.db.UsernameExists(in.Username)
	if username_exists {
		return &pb.RegisterUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "Username already exists",
		}, fmt.Errorf("username already exists")
	}
	email_exists := s.db.EmailExists(in.Email)
	if email_exists {
		return &pb.RegisterUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "Email already exists",
		}, fmt.Errorf("email already exists")
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
			Message: "Error creating user",
		}, fmt.Errorf("error creating user: %w", err)
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
	user, err := s.db.GetUserByUsername(in.Username)
	if err != nil {
		return &pb.CheckUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "User not found",
		}, err
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
