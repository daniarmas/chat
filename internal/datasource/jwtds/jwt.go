package jwtds

import (
	"errors"
	"fmt"
	"time"

	"github.com/daniarmas/chat/internal/config"
	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JwtDatasource interface {
	IsAuthorized(requestToken string) (bool, error)
	CreateApiKey(apiKey *entity.ApiKey) (apikey string, err error)
	CreateRefreshToken(rfToken *entity.RefreshToken, expirationTime time.Time) (refreshToken string, err error)
	CreateAccessToken(acToken *entity.AccessToken, expirationTime time.Time) (accessToken string, err error)
	ExtractTokenClaim(requestToken string) (*models.JwtCustomClaims, error)
	ExtractApiKeyDataFromToken(requestToken string) (string, error)
}

type jwtDatasource struct {
	cfg *config.Config
}

func NewJwtDatasource(cfg *config.Config) JwtDatasource {
	return &jwtDatasource{
		cfg: cfg,
	}
}

func (ds jwtDatasource) CreateAccessToken(acToken *entity.AccessToken, expirationTime time.Time) (accessToken string, err error) {
	claims := &models.JwtCustomClaims{
		ID:     *acToken.ID,
		UserId: *acToken.UserId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(ds.cfg.JwtSecret))
	if err != nil {
		return "", err
	}
	return t, err
}

func (ds jwtDatasource) CreateRefreshToken(rfToken *entity.RefreshToken, expirationTime time.Time) (refreshToken string, err error) {
	claimsRefresh := &models.JwtCustomRefreshClaims{
		ID:     *rfToken.ID,
		UserId: *rfToken.UserId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(ds.cfg.JwtSecret))
	if err != nil {
		return "", err
	}
	return rt, err
}

func (ds jwtDatasource) CreateApiKey(apiKey *entity.ApiKey) (apikey string, err error) {
	claimsRefresh := &models.ApiKeyJwtCustomClaims{
		ID:         *apiKey.ID,
		AppVersion: apiKey.AppVersion,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: apiKey.ExpirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(ds.cfg.JwtSecret))
	if err != nil {
		return "", err
	}
	return rt, err
}

func (ds jwtDatasource) IsAuthorized(requestToken string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ds.cfg.JwtSecret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (ds jwtDatasource) ExtractTokenClaim(requestToken string) (*models.JwtCustomClaims, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ds.cfg.JwtSecret), nil
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

func (ds jwtDatasource) ExtractApiKeyDataFromToken(requestToken string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ds.cfg.JwtSecret), nil
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
