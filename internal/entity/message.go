package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID         *uuid.UUID `json:"id" redis:"id"` // Unique identifier for the message
	Chat       *Chat      `json:"chat" redis:"chat"`
	ChatId     *uuid.UUID `json:"chat_id" redis:"chat_id"`
	User       *User      `json:"user" redis:"user"`
	UserId     *uuid.UUID `json:"user_id" redis:"user_id"`
	Content    string     `json:"content" redis:"content"`         // Content of the message
	CreateTime time.Time  `json:"create_time" redis:"create_time"` // Timestamp of when the message was sent
}

func (m *Message) MarshalBinary() ([]byte, error) {
	// serialize the message to JSON
	serialized, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	// convert the serialized JSON to binary format
	binaryData := make([]byte, len(serialized))
	copy(binaryData, serialized)

	return binaryData, nil
}
