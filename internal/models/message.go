package models

import (
	"time"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageOrm struct {
	ID         *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	SenderID   *uuid.UUID `gorm:"not null" json:"sender_id"`
	ReceiverID *uuid.UUID `gorm:"not null" json:"receiver_id"`
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
	a.ID = message.ID
	a.Content = message.Content
	a.ReceiverID = message.ReceiverID
	a.SenderID = message.SenderID
	a.CreateTime = message.CreateTime
}

func (a MessageOrm) MapFromMessageGorm() *entity.Message {
	return &entity.Message{
		ID:         a.ID,
		Content:    a.Content,
		SenderID:   a.SenderID,
		ReceiverID: a.ReceiverID,
		CreateTime: a.CreateTime,
	}
}
