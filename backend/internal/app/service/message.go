package service

import (
	"github.com/sacurio/jb-challenge/internal/app/model"
	"github.com/sacurio/jb-challenge/internal/app/repository"
)

type (
	Message interface {
		Register(userID uint, content string, timestamp int64) error
		GetLastMessages(count int) []model.Message
	}

	message struct {
		repository repository.Message
	}
)

func NewMessage(repository repository.Message) Message {
	return &message{
		repository: repository,
	}
}
func (m *message) Register(userID uint, content string, timestamp int64) error {
	msg := &model.Message{
		UserID:    userID,
		Content:   content,
		Timestamp: timestamp,
	}

	if err := m.repository.Create(msg); err != nil {
		return err
	}

	return nil
}

func (m *message) GetLastMessages(count int) []model.Message {
	return m.repository.FindLastMessages(count)
}
