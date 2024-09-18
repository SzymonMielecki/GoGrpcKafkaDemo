package types

import (
	pb "github.com/SzymonMielecki/chatApp/usersService"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           string
	Username     string
	Email        string
	PasswordHash string
}

func (u *User) ToProto() *pb.User {
	return &pb.User{
		Id:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}
