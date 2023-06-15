package response

import (
	"time"

	"github.com/daniarmas/chat/internal/entity"
)

type GetMessagesByChatResponse struct {
	Messages []*entity.Message `json:"messages"`
	Cursor   time.Time         `json:"cursor"`
}
