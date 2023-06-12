package entity

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	ID             *uuid.UUID    `json:"id"`
	User           *User         `json:"user"`
	RefreshToken   *RefreshToken `json:"refresh_token"`
	Jwt            string        `json:"jwt"`
	ExpirationTime time.Time     `json:"expiration_time"`
	CreateTime     time.Time     `json:"create_time"`
}
