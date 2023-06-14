package usecases

import (
	"context"
	"time"

	"github.com/daniarmas/chat/config"
	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/inputs"
	"github.com/daniarmas/chat/internal/repository"
	"github.com/daniarmas/chat/pkg/response"
	"github.com/rs/zerolog/log"
)

type MessageUsecase interface {
	SendMessage(ctx context.Context, input inputs.SendMessage, userId string) (*entity.Message, error)
	GetMessageByChat(ctx context.Context, input inputs.GetMessagesByChatId, userId string, createTimeCursor time.Time) (*response.GetMessagesByChatResponse, error)
}

type messageUsecase struct {
	userRepository    repository.UserRepository
	messageRepository repository.MessageRepository
	cfg               *config.Config
}

// NewAuthUsecase will create new an authUsecase object representation of usecases.AuthUsecase interface
func NewMessageUsecase(userRepo repository.UserRepository, messageRepository repository.MessageRepository, cfg *config.Config) MessageUsecase {
	return &messageUsecase{
		userRepository:    userRepo,
		messageRepository: messageRepository,
		cfg:               cfg,
	}
}

func (m *messageUsecase) GetMessageByChat(ctx context.Context, input inputs.GetMessagesByChatId, userId string, createTimeCursor time.Time) (*response.GetMessagesByChatResponse, error) {
	var res response.GetMessagesByChatResponse
	messages, err := m.messageRepository.GetMessagesByChatId(ctx, input.ChatId, createTimeCursor)
	if err != nil {
		log.Fatal().Msgf(err.Error())
		return nil, err
	}
	res.Messages = messages
	return &res, nil
}

func (m *messageUsecase) SendMessage(ctx context.Context, input inputs.SendMessage, userId string) (*entity.Message, error) {
	message, err := m.messageRepository.CreateMessage(ctx, entity.Message{Content: input.Content, ChatId: input.ChatID})
	if err != nil {
		log.Fatal().Msgf(err.Error())
		return nil, err
	}
	return message, nil
}
