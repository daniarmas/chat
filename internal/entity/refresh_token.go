package entity

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID             *uuid.UUID `json:"id"`
	User           *User      `json:"user"`
	ExpirationTime time.Time  `json:"expiration_time"`
	CreateTime     time.Time  `json:"create_time"`
}
