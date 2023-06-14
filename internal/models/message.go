package models

import (
	"time"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageOrm struct {
	ID         *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Chat       ChatOrm    `gorm:"foreignKey:ChatId"`
	ChatId     *uuid.UUID `json:"chat_id"`
	Content    string     `gorm:"not null" json:"content"`
	CreateTime time.Time  `json:"create_time"`
}

func (MessageOrm) TableName() string {
	return "message"
}

func (i *MessageOrm) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	return
}

// This methods map to and from a UserGorm for avoid using gorm models in the usecases.
func (a *MessageOrm) MapToMessageGorm(message *entity.Message) {
	chatOrm := ChatOrm{}
	chatOrm.MapToChatGorm(message.Chat)
	a.ID = message.ID
	a.Chat = chatOrm
	a.ChatId = message.ChatId
	a.Content = message.Content
	a.CreateTime = message.CreateTime
}

func (a MessageOrm) MapFromMessageGorm() *entity.Message {
	return &entity.Message{
		ID:         a.ID,
		Chat:       a.Chat.MapFromChatGorm(),
		ChatId:     a.ChatId,
		Content:    a.Content,
		CreateTime: a.CreateTime,
	}
}
