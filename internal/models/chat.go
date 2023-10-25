package models

import (
	"time"

	"github.com/daniarmas/chat/internal/entity"
)

type Chat struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	CreateTime time.Time `json:"create_time"`
}

// This methods map to and from a ApiKeyModel for avoid using models in the usecases.
func (a *Chat) MapToChatModel(chat *entity.Chat) {
	if chat != nil {
		a.ID = chat.ID
		a.Name = chat.Name
		a.CreateTime = chat.CreateTime
	}
}

func (a Chat) MapFromChatModel() *entity.Chat {
	return &entity.Chat{
		ID:         a.ID,
		Name:       a.Name,
		CreateTime: a.CreateTime,
	}
}
