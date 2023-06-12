package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	ID         *uuid.UUID `json:"id"`         // Unique identifier for the message
	SenderID   *uuid.UUID `json:"senderId"`   // ID of the user who sent the message
	ReceiverID *uuid.UUID `json:"receiverId"` // ID of the user who receive the message
	Content    string     `json:"content"`    // Content of the message
	CreateTime time.Time  `json:"timestamp"`  // Timestamp of when the message was sent
}
