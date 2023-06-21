package models

import (
	"time"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenOrm struct {
	ID             *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	User           UserOrm    `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
	UserId         *uuid.UUID `json:"user_id"`
	ExpirationTime time.Time  `json:"expiration_time"`
	CreateTime     time.Time  `json:"create_time"`
}

func (RefreshTokenOrm) TableName() string {
	return "refresh_token"
}

func (i *RefreshTokenOrm) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC()
	return
}

// This methods map to and from a UserGorm for avoid using gorm models in the usecases.
func (a *RefreshTokenOrm) MapToRefreshTokenGorm(refreshToken *entity.RefreshToken) {
	userOrm := UserOrm{}
	userId := uuid.MustParse(userOrm.ID)
	// userOrm.MapToUserGorm(refreshToken.User)
	a.ID = refreshToken.ID
	a.User = userOrm
	a.UserId = &userId
	a.ExpirationTime = refreshToken.ExpirationTime
	a.CreateTime = refreshToken.CreateTime
}

func (a RefreshTokenOrm) MapFromRefreshTokenGorm() *entity.RefreshToken {
	return &entity.RefreshToken{
		ID: a.ID,
		// User:           a.User.MapFromUserGorm(),
		ExpirationTime: a.ExpirationTime,
		CreateTime:     a.CreateTime,
	}
}
