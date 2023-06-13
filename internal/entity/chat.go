package entity

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID           *uuid.UUID `json:"id"` // Unique identifier for the message
	Channel      string     `json:"channel"`
	FirstUserId  *uuid.UUID `json:"firstUserId"`
	SecondUserId *uuid.UUID `json:"secondUserId"`
	CreateTime   time.Time  `json:"timestamp"`
}
