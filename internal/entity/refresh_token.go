package entity

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID             *uuid.UUID `json:"id"`
	UserId         *uuid.UUID `json:"user_id"`
	ExpirationTime time.Time  `json:"expiration_time"`
	CreateTime     time.Time  `json:"create_time"`
}
