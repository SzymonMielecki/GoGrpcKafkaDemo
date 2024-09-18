package logic

import (
	"context"

	"github.com/SzymonMielecki/chatApp/types"
	"github.com/SzymonMielecki/chatApp/usersServer/persistance"
	pb "github.com/SzymonMielecki/chatApp/usersService"
)

type server struct {
	pb.UnimplementedUsersServiceServer
	db *persistance.DB
}

func NewServer(db *persistance.DB) *server {
	return &server{db: db}
}

func (s *server) LoginUser(ctx context.Context, in *pb.LoginUserRequest) *pb.LoginUserResponse {
	var found types.User
	s.db.Find(&types.User{Username: in.UsernameOrEmail}).Or(&types.User{Email: in.UsernameOrEmail}).First(&found)
	if found.ID == "" {
		return &pb.LoginUserResponse{
			Success: false,
			Id:      "",
			Message: "User not found",
		}
	}
	if found.Password != in.Password {
		return &pb.LoginUserResponse{
			Success: false,
			Id:      "",
			Message: "Incorrect password",
		}
	}
	return &pb.LoginUserResponse{
		Success: true,
		Id:      found.ID,
		Message: "Logged in",
	}
}

func (s *server) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) *pb.RegisterUserResponse {
	var found_username types.User
	s.db.Find(&types.User{Username: in.Username}).First(&found_username)
	if found_username.ID != "" {
		return &pb.RegisterUserResponse{
			Success: false,
			Id:      "",
			Message: "Username already exists",
		}
	}
	var found_email types.User
	s.db.Find(&types.User{Email: in.Email}).First(&found_email)
	if found_email.ID != "" {
		return &pb.RegisterUserResponse{
			Success: false,
			Id:      "",
			Message: "Email already exists",
		}
	}
	user := types.User{
		Username: in.Username,
		Email:    in.Email,
		Password: in.Password,
	}
	err := s.db.CreateUser(&user)
	if err != nil {
		return &pb.RegisterUserResponse{
			Success: false,
			Id:      "",
			Message: "Failed to register",
		}
	}
	return &pb.RegisterUserResponse{
		Success: true,
		Id:      user.ID,
		Message: "Registered",
	}
}
