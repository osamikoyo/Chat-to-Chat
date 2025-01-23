package data

import (
	"github.com/osamikoyo/chat-to-chat/internal/data/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Storage struct {
	data *gorm.DB
}

func New() (*Storage, error) {
	db, err := gorm.Open(sqlite.Open("storage/main.db"))

	db.AutoMigrate(&models.Message{})

	return &Storage{data: db}, err
}

func (s Storage) Save(message models.Message) error {
	return s.data.Create(&message).Error
}

func (s Storage) Get(count int, sender, reciewer string) ([]models.Message, error) {
	var messages []models.Message

	if err := s.data.Order("id desc").Limit(count).Where("reciewer = ? AND sender = ? OR reciewer = ? OR sender = ?",
		sender, reciewer, reciewer, sender).Error; err != nil{
		return nil, err
	}

	return messages, nil
}