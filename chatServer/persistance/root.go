package persistance

import (
	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/types"
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

func (db *DB) CreateMessage(message *types.Message) (*types.Message, error) {
	if err := db.DB.Create(&message).Error; err != nil {
		return nil, err
	}
	return message, nil
}
