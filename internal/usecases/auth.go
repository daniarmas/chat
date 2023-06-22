package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/daniarmas/chat/internal/config"
	"github.com/daniarmas/chat/internal/datasource/hashds"
	"github.com/daniarmas/chat/internal/datasource/jwtds"
	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/inputs"
	"github.com/daniarmas/chat/internal/repository"
	myerror "github.com/daniarmas/chat/pkg/my_error"
	"github.com/daniarmas/chat/pkg/response"
)

type AuthUsecase interface {
	SignIn(ctx context.Context, input inputs.SignInInput) (*response.SignInResponse, error)
	SignOut(ctx context.Context, userId string) error
	Me(ctx context.Context, userId string) (*entity.User, error)
}

type authUsecase struct {
	userRepository    repository.UserRepository
	refreshRepository repository.RefreshTokenRepository
	accessRepository  repository.AccessTokenRepository
	jwtDatasource     jwtds.JwtDatasource
	hashDatasource    hashds.HashDatasource
	cfg               *config.Config
}

// NewAuth will create new an authUsecase object representation of usecases.AuthUsecase interface
func NewAuth(userRepo repository.UserRepository, refreshRepository repository.RefreshTokenRepository, accessRepository repository.AccessTokenRepository, jwtDatasource jwtds.JwtDatasource, hashDatasource hashds.HashDatasource, cfg *config.Config) AuthUsecase {
	return &authUsecase{
		userRepository:    userRepo,
		refreshRepository: refreshRepository,
		accessRepository:  accessRepository,
		cfg:               cfg,
		hashDatasource:    hashDatasource,
		jwtDatasource:     jwtDatasource,
	}
}

func (u *authUsecase) SignOut(ctx context.Context, userId string) error {
	err := u.accessRepository.DeleteAccessTokenByUserId(ctx, userId)
	if err != nil {
		switch err.(type) {
		case myerror.NotFoundError:
			// Do nothing
		default:
			return err
		}
	}
	err = u.refreshRepository.DeleteRefreshTokenByUserId(ctx, userId)
	if err != nil {
		switch err.(type) {
		case myerror.NotFoundError:
			// Do nothing
		default:
			return err
		}
	}
	return nil
}

func (u *authUsecase) Me(ctx context.Context, userId string) (*entity.User, error) {
	user, err := u.userRepository.GetUserById(ctx, userId)
	if err != nil {
		switch err.(type) {
		case myerror.NotFoundError:
			return nil, errors.New("the credentials are incorrect")
		default:
			return nil, err
		}
	}
	return user, nil
}

func (u *authUsecase) SignIn(ctx context.Context, in inputs.SignInInput) (*response.SignInResponse, error) {
	// Check if exists a user with the given email
	user, err := u.userRepository.GetUserByEmail(ctx, in.Email)
	if err != nil {
		switch err.(type) {
		case myerror.NotFoundError:
			return nil, errors.New("the credentials are incorrect")
		default:
			return nil, err
		}
	}
	// Check if the password is correct
	passIsCorrect := u.hashDatasource.CheckHash(in.Password, user.Password)
	if !passIsCorrect {
		return nil, errors.New("the credentials are incorrect")
	}
	// Check if the user is already loged in the system
	refreshTokenCheck, err := u.refreshRepository.GetRefreshTokenByUserId(ctx, user.ID)
	if err != nil {
		switch err.(type) {
		case myerror.NotFoundError:
			// Do nothing
		default:
			return nil, err
		}
	}
	if refreshTokenCheck != nil && in.Logout {
		err = u.accessRepository.DeleteAccessTokenByRefreshTokenId(ctx, refreshTokenCheck.ID)
		if err != nil {
			return nil, err
		}
		err = u.refreshRepository.DeleteRefreshToken(ctx, *refreshTokenCheck)
		if err != nil {
			return nil, err
		}
	} else if refreshTokenCheck != nil {
		return nil, errors.New("the user is already logged in")
	}
	// Create RefreshToken and AccessToken in the database for track sessions.
	refreshTokenExpireTime := time.Now().Add(time.Hour * time.Duration(u.cfg.RefreshTokenExpireHours)).UTC()
	refreshToken, err := u.refreshRepository.CreateRefreshToken(ctx, entity.RefreshToken{UserId: user.ID, ExpirationTime: refreshTokenExpireTime})
	if err != nil {
		return nil, err
	}
	accessTokenExpireTime := time.Now().Add(time.Hour * time.Duration(u.cfg.AccessTokenExpireHours)).UTC()
	accessToken, err := u.accessRepository.CreateAccessToken(ctx, entity.AccessToken{UserId: user.ID, ExpirationTime: accessTokenExpireTime, RefreshTokenId: refreshToken.ID})
	if err != nil {
		return nil, err
	}
	// Create the user jwt for the accessJwt and the refreshToken
	// accessJwt, _ := jwt_utils.CreateAccessToken(, u.cfg.JwtSecret, 1)
	refreshJwt, _ := u.jwtDatasource.CreateRefreshToken(refreshToken, refreshTokenExpireTime)
	accessJwt, _ := u.jwtDatasource.CreateAccessToken(accessToken, accessTokenExpireTime)

	return &response.SignInResponse{
		User:         user,
		AccessToken:  accessJwt,
		RefreshToken: refreshJwt,
	}, nil
}
