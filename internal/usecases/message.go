package usecases

import (
	"context"
	"time"

	"github.com/daniarmas/chat/config"
	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/inputs"
	"github.com/daniarmas/chat/internal/repository"
	"github.com/daniarmas/chat/pkg/response"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type MessageUsecase interface {
	SendMessage(ctx context.Context, input inputs.SendMessage, userId string) (*entity.Message, error)
	GetMessageByChat(ctx context.Context, input inputs.GetMessagesByChat, userId string, createTimeCursor time.Time) (*response.GetMessagesByChatResponse, error)
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

func (m *messageUsecase) GetMessageByChat(ctx context.Context, input inputs.GetMessagesByChat, userId string, createTimeCursor time.Time) (*response.GetMessagesByChatResponse, error) {
	var res response.GetMessagesByChatResponse
	messages, err := m.messageRepository.GetMessagesByChat(ctx, userId, input.ChatUserId, createTimeCursor)
	if err != nil {
		log.Fatal().Msgf(err.Error())
		return nil, err
	}
	res.Messages = messages
	return &res, nil
}

func (m *messageUsecase) SendMessage(ctx context.Context, input inputs.SendMessage, userId string) (*entity.Message, error) {
	userIdUUID := uuid.MustParse(userId)
	message, err := m.messageRepository.CreateMessage(ctx, entity.Message{ReceiverID: input.ReceiverID, Content: input.Content, SenderID: &userIdUUID})
	if err != nil {
		log.Fatal().Msgf(err.Error())
		return nil, err
	}
	return message, nil
}
