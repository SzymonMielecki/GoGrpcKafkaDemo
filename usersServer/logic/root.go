package logic

import (
	"context"
	"strconv"

	"github.com/SzymonMielecki/GoGrpcKafkaDemo/usersServer/cache"
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/usersServer/persistance"

	"github.com/SzymonMielecki/GoGrpcKafkaDemo/types"
	pb "github.com/SzymonMielecki/GoGrpcKafkaDemo/usersService"
)

type Server struct {
	pb.UnimplementedUsersServiceServer
	db persistance.IDB[types.User]
	c  cache.ICache[types.User]
}

func NewServer(db persistance.IDB[types.User], c cache.ICache[types.User]) *Server {
	return &Server{db: db, c: c}
}

func (s *Server) Close() error {
	return s.db.Close()
}

func (s *Server) LoginUser(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	foundUsername, usernameErr := s.c.Handle(
		ctx,
		"username:"+in.Username,
		func() (*types.User, error) {
			return s.db.GetUserByUsername(in.Username)
		})
	foundEmail, emailErr := s.c.Handle(
		ctx,
		"email:"+in.Email,
		func() (*types.User, error) {
			u, err := s.db.GetUserByEmail(in.Email)
			if err != nil {
				return nil, err
			}
			return u, nil
		},
	)
	if usernameErr != nil && emailErr != nil {
		return &pb.LoginUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "User not found",
		}, nil
	}
	if foundUsername != nil {
		return &pb.LoginUserResponse{
			Success: true,
			User:    foundUsername.ToProto(),
			Message: "Logged in as " + foundUsername.Username,
		}, nil
	}
	if foundEmail != nil {
		return &pb.LoginUserResponse{
			Success: true,
			User:    foundEmail.ToProto(),
			Message: "Logged in as " + foundEmail.Username,
		}, nil
	}
	return &pb.LoginUserResponse{
		Success: false,
		User:    &pb.User{},
		Message: "User not found",
	}, nil
}

func (s *Server) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	userFromUsername, _ := s.c.Handle(
		ctx,
		"username:"+in.Username,
		func() (*types.User, error) {
			u, err := s.db.GetUserByUsername(in.Username)
			if err != nil {
				return nil, err
			}
			return u, nil
		},
	)
	if userFromUsername != nil {
		return &pb.RegisterUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "Username already exists",
		}, nil
	}
	userFromEmail, _ := s.c.Handle(
		ctx,
		"email:"+in.Email,
		func() (*types.User, error) {
			u, err := s.db.GetUserByEmail(in.Email)
			if err != nil {
				return nil, err
			}
			return u, nil
		},
	)
	if userFromEmail != nil {
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
			Message: "Error creating user",
		}, nil
	}
	return &pb.RegisterUserResponse{
		Success: true,
		User:    user.ToProto(),
		Message: "Registered",
	}, nil
}

func (s *Server) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.c.Handle(
		ctx,
		"id:"+strconv.Itoa(int(in.Id)),
		func() (*types.User, error) {
			u, err := s.db.GetUserById(uint(in.Id))
			if err != nil {
				return nil, err
			}
			return u, nil
		},
	)
	if err != nil {
		return &pb.GetUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "User not found",
		}, nil
	}
	return &pb.GetUserResponse{
		Success: true,
		User:    user.ToProto(),
		Message: "User found",
	}, nil
}

func (s *Server) CheckUser(ctx context.Context, in *pb.CheckUserRequest) (*pb.CheckUserResponse, error) {
	user, err := s.c.Handle(
		ctx,
		"ID:"+strconv.Itoa(int(in.Id)),
		func() (*types.User, error) {
			u, err := s.db.GetUserById(uint(in.Id))
			if err != nil {
				return nil, err
			}
			return u, nil
		},
	)
	if err != nil {
		return &pb.CheckUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "User not found",
		}, nil
	}
	if user == nil {
		return &pb.CheckUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "User not found",
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
