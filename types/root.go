package types

import (
	pb "github.com/SzymonMielecki/GoGrpcKafkaDemo/usersService"
)

type User struct {
	ID           uint
	Username     string
	Email        string
	PasswordHash string
}

func (u *User) ToProto() *pb.User {
	return &pb.User{
		Id:           uint32(u.ID),
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	}
}

type Message struct {
	ID       uint
	Content  string
	SenderID uint
}
