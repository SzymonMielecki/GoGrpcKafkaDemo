package logic

import (
	"context"
	"strconv"

	"github.com/SzymonMielecki/GoGrpcKafkaDemo/types"
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/usersServer/persistance"
	pb "github.com/SzymonMielecki/GoGrpcKafkaDemo/usersService"
	"github.com/go-redis/cache/v9"
)

type Server struct {
	pb.UnimplementedUsersServiceServer
	db *persistance.DB
	c  *cache.Cache
}

func NewServer(db *persistance.DB, c *cache.Cache) *Server {
	return &Server{db: db, c: c}
}

func (s *Server) LoginUser(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	var found_username *types.User
	username_err := s.c.Once(&cache.Item{
		Key:   "username:" + in.Username,
		Value: &found_username,
		Do: func(i *cache.Item) (interface{}, error) {
			u, err := s.db.GetUserByUsername(in.Username)
			if err != nil {
				return nil, err
			}
			return u, nil
		},
	})
	var found_email *types.User
	email_err := s.c.Once(&cache.Item{
		Key:   "email:" + in.Email,
		Value: &found_email,
		Do: func(i *cache.Item) (interface{}, error) {
			u, err := s.db.GetUserByEmail(in.Email)
			if err != nil {
				return nil, err
			}
			return u, nil
		},
	})
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
	var user_from_username *types.User
	_ = s.c.Once(&cache.Item{
		Key:   "username:" + in.Username,
		Value: &user_from_username,
		Do: func(i *cache.Item) (interface{}, error) {
			u, err := s.db.GetUserByUsername(in.Username)
			if err != nil {
				return nil, err
			}
			return u, nil
		},
	})
	if user_from_username != nil {
		return &pb.RegisterUserResponse{
			Success: false,
			User:    &pb.User{},
			Message: "Username already exists",
		}, nil
	}
	var user_from_email *types.User
	_ = s.c.Once(&cache.Item{
		Key:   "email:" + in.Email,
		Value: &user_from_email,
		Do: func(i *cache.Item) (interface{}, error) {
			u, err := s.db.GetUserByEmail(in.Email)
			if err != nil {
				return nil, err
			}
			return u, nil
		},
	})
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
	var user *types.User
	err := s.c.Once(&cache.Item{
		Key:   "id:" + strconv.Itoa(int(in.Id)),
		Value: &user,
		Do: func(i *cache.Item) (interface{}, error) {
			u, err := s.db.GetUserById(uint(in.Id))
			if err != nil {
				return nil, err
			}
			return u, nil
		},
	})
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
	var user *types.User
	err := s.c.Once(&cache.Item{
		Key:   "ID:" + strconv.Itoa(int(in.Id)),
		Value: &user,
		Do: func(i *cache.Item) (interface{}, error) {
			u, err := s.db.GetUserById(uint(in.Id))
			if err != nil {
				return nil, err
			}
			return u, nil
		},
	})
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
