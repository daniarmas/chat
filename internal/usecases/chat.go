package usecases

import (
	"context"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/inputs"
	"github.com/daniarmas/chat/internal/repository"
	myerror "github.com/daniarmas/chat/pkg/my_error"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type ChatUsecase interface {
	GetOrCreateChat(ctx context.Context, input inputs.GetOrCreateChatInput, userId string) (*entity.Chat, error)
}

type chatUsecase struct {
	chatRepository repository.ChatRepository
}

func NewChatUsecase(chatRepo repository.ChatRepository) ChatUsecase {
	return &chatUsecase{
		chatRepository: chatRepo,
	}
}

func (u chatUsecase) GetOrCreateChat(ctx context.Context, input inputs.GetOrCreateChatInput, userId string) (*entity.Chat, error) {
	userIdUUID := uuid.MustParse(userId)
	chat, err := u.chatRepository.GetChat(ctx, userId, input.OtherUserId.String())
	switch err.(type) {
	case nil:
		// Do nothing
	case myerror.NotFoundError:
		chat, err = u.chatRepository.CreateChat(ctx, entity.Chat{FirstUserId: &userIdUUID, SecondUserId: input.OtherUserId})
		if err != nil {
			log.Error().Msgf(err.Error())
			return nil, err
		}
	default:
		log.Error().Msgf(err.Error())
		return nil, err
	}
	return chat, nil
}
