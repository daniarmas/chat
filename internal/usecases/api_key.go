package usecases

import (
	"context"

	"github.com/daniarmas/chat/internal/datasource/jwtds"
	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/inputs"
	"github.com/daniarmas/chat/internal/repository"
	"github.com/rs/zerolog/log"
)

type ApiKeyUsecase interface {
	CreateApiKey(ctx context.Context, input inputs.CreateApiKeyInput) (*entity.ApiKey, error)
}

type apiKeyUsecase struct {
	apiKeyRepository repository.ApiKeyRepository
	jwtDatasource    jwtds.JwtDatasource
}

func NewApiKey(apiKeyRepo repository.ApiKeyRepository, jwtDatasource jwtds.JwtDatasource) ApiKeyUsecase {
	return &apiKeyUsecase{
		apiKeyRepository: apiKeyRepo,
		jwtDatasource:    jwtDatasource,
	}
}

func (u apiKeyUsecase) CreateApiKey(ctx context.Context, input inputs.CreateApiKeyInput) (*entity.ApiKey, error) {
	apiKey, err := u.apiKeyRepository.CreateApiKey(ctx, &entity.ApiKey{AppVersion: input.AppVersion, Revoked: input.Revoked, ExpirationTime: input.ExpirationTime})
	if err != nil {
		go log.Error().Msgf(err.Error())
		return nil, err
	}
	apiKeyJwt, _ := u.jwtDatasource.CreateApiKey(apiKey)
	apiKey.Jwt = apiKeyJwt
	return apiKey, nil
}
