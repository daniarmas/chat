package entity

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID         *uuid.UUID `json:"id"` // Unique identifier for the message
	Chat       *Chat      `json:"chat"`
	ChatId     *uuid.UUID `json:"chat_id"`
	Content    string     `json:"content"`   // Content of the message
	CreateTime time.Time  `json:"timestamp"` // Timestamp of when the message was sent
}
