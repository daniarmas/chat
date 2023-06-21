package entity

import (
	"time"
)

type AccessToken struct {
	ID             string        `json:"id"`
	User           *User         `json:"user"`
	UserId         string        `json:"user_id"`
	RefreshToken   *RefreshToken `json:"refresh_token"`
	RefreshTokenId string        `json:"refresh_token_id"`
	Jwt            string        `json:"jwt"`
	ExpirationTime time.Time     `json:"expiration_time"`
	CreateTime     time.Time     `json:"create_time"`
}
