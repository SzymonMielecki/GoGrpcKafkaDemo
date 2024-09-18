package persistance

import (
	"github.com/SzymonMielecki/chatApp/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB(host, user, password, dbname, port string) (*DB, error) {
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable TimeZone=Europe/Warsaw"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&types.User{})
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) CreateUser(user *types.User) error {
	return db.DB.Create(user).Error
}

func (db *DB) GetUserById(id string) (*types.User, error) {
	var user types.User

	return &user, db.DB.Where("id = ?", id).First(&user).Error
}

func (db *DB) GetUserByEmail(email string) (*types.User, error) {
	var user types.User
	return &user, db.DB.Where("email = ?", email).First(&user).Error
}

func (db *DB) GetUserByUsername(username string) (*types.User, error) {
	var user types.User
	return &user, db.DB.Where("username = ?", username).First(&user).Error
}
