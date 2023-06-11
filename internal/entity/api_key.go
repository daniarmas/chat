package entity

import (
	"time"

	"github.com/google/uuid"
)

type ApiKey struct {
	ID             *uuid.UUID `json:"id"`
	AppVersion     string     `json:"app_version"`
	Revoked        bool       `json:"revoked"`
	Jwt            string     `json:"jwt"`
	ExpirationTime time.Time  `json:"expiration_time"`
	CreateTime     time.Time  `json:"create_time"`
}
