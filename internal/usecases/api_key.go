package usecases

import (
	"context"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/inputs"
	"github.com/daniarmas/chat/internal/repository"
	"github.com/daniarmas/chat/pkg/jwt_utils"
	"github.com/rs/zerolog/log"
)

type ApiKeyUsecase interface {
	CreateApiKey(ctx context.Context, input inputs.CreateApiKeyInput) (*entity.ApiKey, error)
}

type apiKeyUsecase struct {
	apiKeyRepository repository.ApiKeyRepository
}

func NewApiKeyUsecase(apiKeyRepo repository.ApiKeyRepository) ApiKeyUsecase {
	return &apiKeyUsecase{
		apiKeyRepository: apiKeyRepo,
	}
}

func (u apiKeyUsecase) CreateApiKey(ctx context.Context, input inputs.CreateApiKeyInput) (*entity.ApiKey, error) {
	apiKey, err := u.apiKeyRepository.CreateApiKey(ctx, &entity.ApiKey{AppVersion: input.AppVersion, Revoked: input.Revoked, ExpirationTime: input.ExpirationTime})
	if err != nil {
		log.Fatal().Msgf(err.Error())
		return nil, err
	}
	apiKeyJwt, _ := jwt_utils.CreateApiKey(apiKey)
	apiKey.Jwt = apiKeyJwt
	return apiKey, nil
}
