package jwt_utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/daniarmas/chat/config"
	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func CreateAccessToken(acToken *entity.AccessToken, secret string, expirationTime time.Time) (accessToken string, err error) {
	claims := &models.JwtCustomClaims{
		ID:     *acToken.ID,
		UserId: *acToken.User.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, err
}

func CreateRefreshToken(rfToken *entity.RefreshToken, secret string, expirationTime time.Time) (refreshToken string, err error) {
	claimsRefresh := &models.JwtCustomRefreshClaims{
		ID:     *rfToken.ID,
		UserId: *rfToken.User.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return rt, err
}

func CreateApiKey(apiKey *entity.ApiKey) (apikey string, err error) {
	cfg := config.NewConfig()
	claimsRefresh := &models.ApiKeyJwtCustomClaims{
		ID:         *apiKey.ID,
		AppVersion: apiKey.AppVersion,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: apiKey.ExpirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(cfg.JwtSecret))
	if err != nil {
		return "", err
	}
	return rt, err
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractTokenClaim(requestToken string, secret string) (*models.JwtCustomClaims, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return nil, errors.New("invalid token")
	}

	// return claims["id"].(string), nil
	return &models.JwtCustomClaims{
		ID:     uuid.MustParse(claims["id"].(string)),
		UserId: uuid.MustParse(claims["user_id"].(string)),
	}, nil
}

func ExtractApiKeyDataFromToken(requestToken string, secret string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return "", fmt.Errorf("invalid Token")
	}

	return claims["id"].(string), nil
}
