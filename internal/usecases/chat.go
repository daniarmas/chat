package usecases

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/inputs"
	"github.com/daniarmas/chat/internal/repository"
	myerror "github.com/daniarmas/chat/pkg/my_error"
	"github.com/daniarmas/chat/pkg/response"
	"github.com/rs/zerolog/log"
)

type ChatUsecase interface {
	GetOrCreateChat(ctx context.Context, input inputs.GetOrCreateChatInput, userId string) (*entity.Chat, error)
	GetChats(ctx context.Context, userId string, updateTimeCursor time.Time) (*response.GetChatsResponse, error)
}

type chatUsecase struct {
	chatRepository repository.ChatRepository
}

func NewChat(chatRepo repository.ChatRepository) ChatUsecase {
	return &chatUsecase{
		chatRepository: chatRepo,
	}
}

func (u chatUsecase) GetOrCreateChat(ctx context.Context, input inputs.GetOrCreateChatInput, userId string) (*entity.Chat, error) {
	chat, err := u.chatRepository.GetChat(ctx, userId, input.ReceiverId)
	switch err.(type) {
	case nil:
		// Do nothing
	case myerror.NotFoundError:
		chat, err = u.chatRepository.CreateChat(ctx, &entity.Chat{})
		if err != nil {
			go log.Error().Msgf(err.Error())
			return nil, err
		}
	default:
		go log.Error().Msgf(err.Error())
		return nil, err
	}
	return chat, nil
}

func (u chatUsecase) GetChats(ctx context.Context, userId string, updateTimeCursor time.Time) (*response.GetChatsResponse, error) {
	var res response.GetChatsResponse
	chats, err := u.chatRepository.GetChats(ctx, userId, updateTimeCursor)
	if err != nil {
		go log.Error().Msgf(err.Error())
		return nil, err
	}
	res.Chats = chats
	// if len(chats) != 0 {
	// 	res.Cursor = chats[len(chats)-1].UpdateTime
	// }
	return &res, nil

}
