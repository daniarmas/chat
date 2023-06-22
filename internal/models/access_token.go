package models

import (
	"time"

	"github.com/daniarmas/chat/internal/entity"
)

type AccessToken struct {
	ID             string    `json:"id" redis:"id"`
	UserId         string    `json:"user_id" redis:"user_id"`
	RefreshTokenId string    `json:"refresh_token_id" redis:"refresh_token_id"`
	ExpirationTime time.Time `json:"expiration_time" redis:"expiration_time"`
	CreateTime     time.Time `json:"create_time" redis:"create_time"`
}

// This methods map to and from a UserModel for avoid using models in the usecases.
func (a *AccessToken) MapToAccessTokenModel(accessToken *entity.AccessToken) {
	refreshTokenOrm := RefreshToken{}
	refreshTokenOrm.MapToRefreshTokenModel(accessToken.RefreshToken)
	userOrm := User{}
	userOrm.MapToUserModel(accessToken.User)
	a.ID = accessToken.ID
	a.UserId = userOrm.ID
	a.RefreshTokenId = accessToken.ID
	a.ExpirationTime = accessToken.ExpirationTime
	a.CreateTime = accessToken.CreateTime
}

func (a AccessToken) MapFromAccessTokenModel() *entity.AccessToken {
	return &entity.AccessToken{
		ID:             a.ID,
		UserId:         a.UserId,
		RefreshTokenId: a.RefreshTokenId,
		ExpirationTime: a.ExpirationTime,
		CreateTime:     a.CreateTime,
	}
}
