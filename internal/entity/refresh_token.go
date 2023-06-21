package entity

import (
	"time"
)

type RefreshToken struct {
	ID             string    `json:"id"`
	UserId         string    `json:"user_id"`
	ExpirationTime time.Time `json:"expiration_time"`
	CreateTime     time.Time `json:"create_time"`
}
