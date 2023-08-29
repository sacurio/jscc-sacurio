package repository

import (
	"github.com/sacurio/jb-challenge/internal/app/model"
	"gorm.io/gorm"
)

type (
	Message interface {
		Create(message *model.Message) error
		FindLastMessages(count int) []model.Message
	}

	message struct {
		db *gorm.DB
	}
)

func NewMessage(db *gorm.DB) Message {
	return &message{
		db: db,
	}
}

func (m *message) Create(message *model.Message) error {
	result := m.db.Create(message)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m *message) FindLastMessages(count int) []model.Message {
	var messages []model.Message
	m.db.Preload("User").Order("id").Limit(50).Find(&messages)

	return messages
}
