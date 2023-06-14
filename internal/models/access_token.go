package models

import (
	"time"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccessTokenOrm struct {
	ID             *uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	User           UserOrm         `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
	UserId         *uuid.UUID      `json:"user_id"`
	RefreshTokenId *uuid.UUID      `json:"refresh_token_id"`
	RefreshToken   RefreshTokenOrm `gorm:"foreignKey:RefreshTokenId;constraint:OnDelete:CASCADE;"`
	ExpirationTime time.Time       `json:"expiration_time"`
	CreateTime     time.Time       `json:"create_time"`
}

func (AccessTokenOrm) TableName() string {
	return "access_token"
}

func (i *AccessTokenOrm) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	return
}

// This methods map to and from a UserGorm for avoid using gorm models in the usecases.
func (a *AccessTokenOrm) MapToAccessTokenGorm(accessToken *entity.AccessToken) {
	refreshTokenOrm := RefreshTokenOrm{}
	refreshTokenOrm.MapToRefreshTokenGorm(accessToken.RefreshToken)
	userOrm := UserOrm{}
	userOrm.MapToUserGorm(accessToken.User)
	a.ID = accessToken.ID
	a.User = userOrm
	a.UserId = userOrm.ID
	a.RefreshTokenId = accessToken.ID
	a.RefreshToken = refreshTokenOrm
	a.ExpirationTime = accessToken.ExpirationTime
	a.CreateTime = accessToken.CreateTime
}

func (a AccessTokenOrm) MapFromAccessTokenGorm() *entity.AccessToken {
	return &entity.AccessToken{
		ID:             a.ID,
		User:           a.User.MapFromUserGorm(),
		RefreshToken:   a.RefreshToken.MapFromRefreshTokenGorm(),
		ExpirationTime: a.ExpirationTime,
		CreateTime:     a.CreateTime,
	}
}
