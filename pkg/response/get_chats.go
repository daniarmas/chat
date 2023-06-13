package response

import (
	"time"

	"github.com/daniarmas/chat/internal/entity"
)

type GetChatsResponse struct {
	Chats  []*entity.Chat `json:"chats"`
	Cursor time.Time      `json:"cursor"`
}
