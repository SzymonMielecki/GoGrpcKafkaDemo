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
	var found types.User
	s.db.Find(&types.User{Username: in.UsernameOrEmail}).Or(&types.User{Email: in.UsernameOrEmail}).First(&found)
	if found.ID == "" {
		return &pb.LoginUserResponse{
			Success: false,
			Id:      "",
			Message: "User not found",
		}, nil
	}
	if found.PasswordHash != in.PasswordHash {
		return &pb.LoginUserResponse{
			Success: false,
			Id:      "",
			Message: "Incorrect password",
		}, nil
	}
	return &pb.LoginUserResponse{
		Success: true,
		Id:      found.ID,
		Message: "Logged in",
	}, nil
}

func (s *Server) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	found_username, err := s.db.GetUserByUsername(in.Username)
	if err != nil {
		return &pb.RegisterUserResponse{
			Success: false,
			Id:      "",
			Message: "Username already exists",
		}, nil
	}
	if found_username.ID != "" {
		return &pb.RegisterUserResponse{
			Success: false,
			Id:      "",
			Message: "Username already exists",
		}, nil
	}
	found_email, err := s.db.GetUserByEmail(in.Email)
	if err != nil {
		return &pb.RegisterUserResponse{
			Success: false,
			Id:      "",
			Message: "Email already exists",
		}, nil
	}
	if found_email.ID != "" {
		return &pb.RegisterUserResponse{
			Success: false,
			Id:      "",
			Message: "Email already exists",
		}, nil
	}
	user := types.User{
		Username:     in.Username,
		Email:        in.Email,
		PasswordHash: in.PasswordHash,
	}
	err = s.db.CreateUser(&user)
	if err != nil {
		return &pb.RegisterUserResponse{
			Success: false,
			Id:      "",
			Message: "Failed to register",
		}, err
	}
	return &pb.RegisterUserResponse{
		Success: true,
		Id:      user.ID,
		Message: "Registered",
	}, nil
}

func (s *Server) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.db.GetUserById(in.Id)
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
			Id:      "",
			Message: "User not found",
		}, err
	}
	if user.PasswordHash != in.PasswordHash {
		return &pb.CheckUserResponse{
			Success: false,
			Id:      "",
			Message: "Incorrect password",
		}, nil
	}
	return &pb.CheckUserResponse{
		Success: true,
		Id:      user.ID,
		Message: "User found",
	}, nil
}
