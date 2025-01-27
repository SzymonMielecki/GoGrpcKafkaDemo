package persistance

import (
	"errors"
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/types"
)

type MockDB struct {
	db []*types.User
}

func NewMockDB() *MockDB {
	return &MockDB{}
}

func NewMockDBFromData(db []*types.User) *MockDB {
	return &MockDB{db}
}

func (m MockDB) Close() error {
	return nil
}

func (m MockDB) CreateUser(user *types.User) (*types.User, error) {
	if user == nil {
		return nil, errors.New("User is nil")
	}
	if m.UsernameExists(user.Username) {
		return nil, errors.New("Username already exists")
	}
	if m.EmailExists(user.Email) {
		return nil, errors.New("Email already exists")
	}
	user.ID = uint(len(m.db) + 1)
	m.db = append(m.db, user)
	return user, nil
}

func (m MockDB) GetUserById(userId uint) (*types.User, error) {
	for _, user := range m.db {
		if user.ID == userId {
			return user, nil
		}
	}
	return nil, errors.New("User not found")
}

func (m MockDB) GetUserByUsername(username string) (*types.User, error) {
	for _, user := range m.db {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("User not found")
}

func (m MockDB) GetUserByEmail(email string) (*types.User, error) {
	for _, user := range m.db {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("User not found")
}

func (m MockDB) UsernameExists(username string) bool {
	for _, user := range m.db {
		if user.Username == username {
			return true
		}
	}
	return false
}

func (m MockDB) EmailExists(email string) bool {
	for _, user := range m.db {
		if user.Email == email {
			return true
		}
	}
	return false
}
