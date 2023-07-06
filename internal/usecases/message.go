package usecases

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/config"
	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/inputs"
	"github.com/daniarmas/chat/internal/repository"
	"github.com/daniarmas/chat/pkg/response"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type MessageUsecase interface {
	ReceiveMessagesByChat(ctx context.Context, input inputs.ReceiveMessagesInput) (chan *entity.Message, error)
	ReceiveMessages(ctx context.Context, userId string) (chan *entity.Message, error)
	SendMessage(ctx context.Context, input inputs.SendMessage, userId string) (*entity.Message, error)
	GetMessageByChat(ctx context.Context, input inputs.GetMessagesByChatId, userId string, createTimeCursor time.Time) (*response.GetMessagesByChatResponse, error)
}

type messageUsecase struct {
	userRepository    repository.UserRepository
	messageRepository repository.MessageRepository
	chatRepository    repository.ChatRepository
	cfg               *config.Config
	redis             *redis.Client
}

// NewAuthUsecase will create new an authUsecase object representation of usecases.AuthUsecase interface
func NewMessage(userRepo repository.UserRepository, messageRepository repository.MessageRepository, chatRepository repository.ChatRepository, cfg *config.Config, redis *redis.Client) MessageUsecase {
	return &messageUsecase{
		userRepository:    userRepo,
		messageRepository: messageRepository,
		chatRepository:    chatRepository,
		cfg:               cfg,
		redis:             redis,
	}
}

func (usecase *messageUsecase) ReceiveMessages(ctx context.Context, userId string) (chan *entity.Message, error) {
	res, err := usecase.messageRepository.ReceiveMessageByUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (usecase *messageUsecase) ReceiveMessagesByChat(ctx context.Context, input inputs.ReceiveMessagesInput) (chan *entity.Message, error) {
	res, err := usecase.messageRepository.ReceiveMessageByChat(ctx, input.ChatId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *messageUsecase) GetMessageByChat(ctx context.Context, input inputs.GetMessagesByChatId, userId string, createTimeCursor time.Time) (*response.GetMessagesByChatResponse, error) {
	var res response.GetMessagesByChatResponse
	messages, err := m.messageRepository.GetMessagesByChatId(ctx, input.ChatId, createTimeCursor)
	if err != nil {
		go log.Error().Msgf(err.Error())
		return nil, err
	}
	res.Messages = messages
	return &res, nil
}

func (usecase *messageUsecase) SendMessage(ctx context.Context, input inputs.SendMessage, userId string) (*entity.Message, error) {
	message, err := usecase.messageRepository.CreateMessage(ctx, entity.Message{Content: input.Content, ChatId: input.ChatID, UserId: userId})
	if err != nil {
		go log.Error().Msgf(err.Error())
		return nil, err
	}
	chat, err := usecase.chatRepository.GetChatById(ctx, message.ChatId)
	if err != nil {
		go log.Error().Msgf(err.Error())
		return nil, err
	}
	var otherUserId string
	if chat.FirstUserId != userId {
		otherUserId = chat.FirstUserId
	} else {
		otherUserId = chat.SecondUserId
	}
	err = usecase.messageRepository.PublishMessage(ctx, *message, otherUserId)
	if err != nil {
		panic(err)
	}
	return message, nil
}
