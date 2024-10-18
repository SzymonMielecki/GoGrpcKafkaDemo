package persistance

import (
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/types"
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

func (db *DB) CreateUser(user *types.User) (*types.User, error) {
	err := db.DB.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DB) GetUserById(id uint) (*types.User, error) {
	var user types.User

	err := db.DB.Model(&types.User{
		Model: gorm.Model{
			ID: id,
		},
	}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) GetUserByEmail(email string) (*types.User, error) {
	var user types.User
	err := db.DB.Model(&types.User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) GetUserByUsername(username string) (*types.User, error) {
	var user types.User
	err := db.DB.Model(&types.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) GetUserByUsernameAndEmail(username, email string) (*types.User, error) {
	var user types.User
	err := db.DB.Model(&types.User{}).Where("username = ? AND email = ?", username, email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) UsernameExists(username string) bool {
	var user types.User
	err := db.DB.Model(&types.User{}).Where("username = ?", username).First(&user).Error
	return err == nil
}

func (db *DB) EmailExists(email string) bool {
	var user types.User
	err := db.DB.Model(&types.User{}).Where("email = ?", email).First(&user).Error
	return err == nil
}
