package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JwtCustomClaims struct {
	ID     string `json:"id"`
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

type JwtCustomRefreshClaims struct {
	ID     uuid.UUID `json:"id"`
	UserId uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

type ApiKeyJwtCustomClaims struct {
	ID         uuid.UUID `json:"id"`
	AppVersion string    `json:"app_version"`
	jwt.StandardClaims
}
