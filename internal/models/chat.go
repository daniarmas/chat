package models

import (
	"time"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatOrm struct {
	ID           *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Channel      string     `gorm:"not null" json:"channel"`
	FirstUserId  *uuid.UUID `gorm:"not null" json:"firstUserId"`
	SecondUserId *uuid.UUID `gorm:"not null" json:"secondUserId"`
	CreateTime   time.Time  `json:"timestamp"`
}

func (ChatOrm) TableName() string {
	return "chat"
}

func (i *ChatOrm) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	return
}

// This methods map to and from a ApiKeyGorm for avoid using gorm models in the usecases.
func (a *ChatOrm) MapToChatGorm(chat *entity.Chat) {
	a.ID = chat.ID
	a.Channel = chat.Channel
	a.FirstUserId = chat.FirstUserId
	a.SecondUserId = chat.SecondUserId
	a.CreateTime = chat.CreateTime
}

func (a ChatOrm) MapFromChatGorm() *entity.Chat {
	return &entity.Chat{
		ID:           a.ID,
		Channel:      a.Channel,
		FirstUserId:  a.FirstUserId,
		SecondUserId: a.SecondUserId,
		CreateTime:   a.CreateTime,
	}
}
