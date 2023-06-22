package models

import (
	"time"

	"github.com/daniarmas/chat/internal/entity"
)

type ApiKey struct {
	ID             string    `json:"id"`
	AppVersion     string    `json:"app_version"`
	Revoked        bool      `json:"revoked"`
	ExpirationTime time.Time `json:"expiration_time"`
	CreateTime     time.Time `json:"create_time"`
}

// This methods map to and from a ApiKeyModel for avoid using models in the usecases.
func (a *ApiKey) MapToApiKeyModel(apiKey *entity.ApiKey) {
	a.ID = apiKey.ID
	a.AppVersion = apiKey.AppVersion
	a.Revoked = apiKey.Revoked
	a.ExpirationTime = apiKey.ExpirationTime
	a.CreateTime = apiKey.CreateTime
}

func (a ApiKey) MapFromApiKeyModel() *entity.ApiKey {
	return &entity.ApiKey{
		ID:             a.ID,
		AppVersion:     a.AppVersion,
		Revoked:        a.Revoked,
		ExpirationTime: a.ExpirationTime,
		CreateTime:     a.CreateTime,
	}
}
