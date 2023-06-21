package models

import (
	"time"

	"github.com/daniarmas/chat/internal/entity"
)

type RefreshToken struct {
	ID             string    `json:"id"`
	UserId         string    `json:"user_id"`
	ExpirationTime time.Time `json:"expiration_time"`
	CreateTime     time.Time `json:"create_time"`
}

// This methods map to and from a UserGorm for avoid using gorm models in the usecases.
func (a *RefreshToken) MapToRefreshTokenModel(refreshToken *entity.RefreshToken) {
	userOrm := User{}
	// userOrm.MapToUserGorm(refreshToken.User)
	a.ID = refreshToken.ID
	a.UserId = userOrm.ID
	a.ExpirationTime = refreshToken.ExpirationTime
	a.CreateTime = refreshToken.CreateTime
}

func (a RefreshToken) MapFromRefreshTokenModel() *entity.RefreshToken {
	return &entity.RefreshToken{
		ID:             a.ID,
		ExpirationTime: a.ExpirationTime,
		CreateTime:     a.CreateTime,
	}
}
