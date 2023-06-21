package entity

import (
	"time"
)

type ApiKey struct {
	ID             string    `json:"id"`
	AppVersion     string    `json:"app_version"`
	Revoked        bool      `json:"revoked"`
	Jwt            string    `json:"jwt"`
	ExpirationTime time.Time `json:"expiration_time"`
	CreateTime     time.Time `json:"create_time"`
}
