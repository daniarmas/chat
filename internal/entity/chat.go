package entity

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID           *uuid.UUID `json:"id"` // Unique identifier for the message
	FirstUserId  *uuid.UUID `json:"firstUserId"`
	SecondUserId *uuid.UUID `json:"secondUserId"`
	CreateTime   time.Time  `json:"create_time"`
	UpdateTime   time.Time  `json:"update_time"`
}
