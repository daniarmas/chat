package models

import (
	"time"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApiKeyOrm struct {
	ID             *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	AppVersion     string     `gorm:"not null" json:"app_version"`
	Revoked        bool       `gorm:"default:false" json:"revoked"`
	ExpirationTime time.Time  `json:"expiration_time"`
	CreateTime     time.Time  `json:"create_time"`
}

func (ApiKeyOrm) TableName() string {
	return "api_key"
}

func (i *ApiKeyOrm) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ExpirationTime.IsZero() {
		// Add 3 month to time.Now for set 3 month of expiration time.
		expirationTime := time.Now().AddDate(0, 3, 0).UTC()
		i.ExpirationTime = expirationTime
	}
	i.CreateTime = time.Now().UTC()
	return
}

// This methods map to and from a ApiKeyGorm for avoid using gorm models in the usecases.
func (a *ApiKeyOrm) MapToApiKeyGorm(apiKey *entity.ApiKey) {
	a.ID = apiKey.ID
	a.AppVersion = apiKey.AppVersion
	a.Revoked = apiKey.Revoked
	a.ExpirationTime = apiKey.ExpirationTime
	a.CreateTime = apiKey.CreateTime
}

func (a ApiKeyOrm) MapFromApiKeyGorm() *entity.ApiKey {
	return &entity.ApiKey{
		ID:             a.ID,
		AppVersion:     a.AppVersion,
		Revoked:        a.Revoked,
		ExpirationTime: a.ExpirationTime,
		CreateTime:     a.CreateTime,
	}
}
