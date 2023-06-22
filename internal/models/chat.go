package models

import (
	"time"

	"github.com/daniarmas/chat/internal/entity"
)

type Chat struct {
	ID           string    `json:"id"`
	FirstUserId  string    `json:"firstUserId"`
	SecondUserId string    `json:"secondUserId"`
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"`
}

// This methods map to and from a ApiKeyModel for avoid using models in the usecases.
func (a *Chat) MapToChatModel(chat *entity.Chat) {
	if chat != nil {
		a.ID = chat.ID
		a.FirstUserId = chat.FirstUserId
		a.SecondUserId = chat.SecondUserId
		a.CreateTime = chat.CreateTime
		a.UpdateTime = chat.UpdateTime
	}
}

func (a Chat) MapFromChatModel() *entity.Chat {
	return &entity.Chat{
		ID:           a.ID,
		FirstUserId:  a.FirstUserId,
		SecondUserId: a.SecondUserId,
		CreateTime:   a.CreateTime,
		UpdateTime:   a.UpdateTime,
	}
}
