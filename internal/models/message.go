package models

import (
	"time"

	"github.com/daniarmas/chat/internal/entity"
)

type Message struct {
	ID         string    `json:"id"`
	ChatId     string    `json:"chat_id"`
	UserId     string    `json:"user_id"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"create_time"`
}

// This methods map to and from a UserGorm for avoid using gorm models in the usecases.
func (a *Message) MapToMessageModel(message *entity.Message) {
	chatOrm := Chat{}
	chatOrm.MapToChatModel(message.Chat)
	userOrm := User{}
	userOrm.MapToUserModel(message.User)
	a.ID = message.ID
	a.ChatId = message.ChatId
	a.UserId = message.UserId
	a.Content = message.Content
	a.CreateTime = message.CreateTime
}

func (a Message) MapFromMessageModel() *entity.Message {
	return &entity.Message{
		ID:         a.ID,
		ChatId:     a.ChatId,
		UserId:     a.UserId,
		Content:    a.Content,
		CreateTime: a.CreateTime,
	}
}
