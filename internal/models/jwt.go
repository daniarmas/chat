package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type JwtCustomClaims struct {
	ID     string `json:"id"`
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

type JwtCustomRefreshClaims struct {
	ID     string `json:"id"`
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

type ApiKeyJwtCustomClaims struct {
	ID         string `json:"id"`
	AppVersion string `json:"app_version"`
	jwt.StandardClaims
}
