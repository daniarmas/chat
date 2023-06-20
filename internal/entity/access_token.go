package entity

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	ID             *uuid.UUID    `json:"id"`
	User           *User         `json:"user"`
	UserId         *uuid.UUID    `json:"user_id"`
	RefreshToken   *RefreshToken `json:"refresh_token"`
	RefreshTokenId *uuid.UUID    `json:"refresh_token_id"`
	Jwt            string        `json:"jwt"`
	ExpirationTime time.Time     `json:"expiration_time"`
	CreateTime     time.Time     `json:"create_time"`
}
