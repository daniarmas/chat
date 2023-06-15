package usecases

import (
	"context"
	"encoding/json"
	"time"

	"github.com/daniarmas/chat/internal/config"
	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/inputs"
	"github.com/daniarmas/chat/internal/repository"
	"github.com/daniarmas/chat/pkg/response"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type MessageUsecase interface {
	ReceiveMessagesByChat(ctx context.Context, input inputs.ReceiveMessagesInput) (<-chan *entity.Message, error)
	ReceiveMessages(ctx context.Context, userId string) (<-chan *entity.Message, error)
	SendMessage(ctx context.Context, input inputs.SendMessage, userId uuid.UUID) (*entity.Message, error)
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

func (usecase *messageUsecase) ReceiveMessages(ctx context.Context, userId string) (<-chan *entity.Message, error) {
	ch := make(chan *entity.Message)

	go func() {
		// There is no error because go-redis automatically reconnects on error.
		pubsub := usecase.redis.Subscribe(ctx, userId)

		// Close the subscription when we are done.
		defer pubsub.Close()

		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				panic(err)
			}

			// parse the message payload into a Message object
			var messageObj entity.Message
			err = json.Unmarshal([]byte(msg.Payload), &messageObj)
			if err != nil {
				panic(err) // handle the error appropriately
			}

			// The channel may have gotten closed due to the client disconnecting.
			// To not have our Goroutine block or panic, we do the send in a select block.
			// This will jump to the default case if the channel is closed.
			select {
			case ch <- &messageObj: // This is the actual send.
				go log.Info().Msgf("Msg Sended: %s", messageObj.Content)
				// Our message went through, do nothing
			default: // This is run when our send does not work.
				go log.Info().Msgf("Channel usecase closed")
				// You can handle any deregistration of the channel here.
				return // We'll just return ending the routine.
			}
		}
	}()

	return ch, nil
}

func (usecase *messageUsecase) ReceiveMessagesByChat(ctx context.Context, input inputs.ReceiveMessagesInput) (<-chan *entity.Message, error) {
	ch := make(chan *entity.Message)

	go func() {
		// There is no error because go-redis automatically reconnects on error.
		pubsub := usecase.redis.Subscribe(ctx, input.ChatId)

		// Close the subscription when we are done.
		defer pubsub.Close()

		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				panic(err)
			}

			// parse the message payload into a Message object
			var messageObj entity.Message
			err = json.Unmarshal([]byte(msg.Payload), &messageObj)
			if err != nil {
				panic(err) // handle the error appropriately
			}

			// The channel may have gotten closed due to the client disconnecting.
			// To not have our Goroutine block or panic, we do the send in a select block.
			// This will jump to the default case if the channel is closed.
			select {
			case ch <- &messageObj: // This is the actual send.
				go log.Info().Msgf("Msg Sended: %s", messageObj.Content)
				// Our message went through, do nothing
			default: // This is run when our send does not work.
				go log.Info().Msgf("Channel usecase closed")
				// You can handle any deregistration of the channel here.
				return // We'll just return ending the routine.
			}
		}
	}()

	return ch, nil
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

func (usecase *messageUsecase) SendMessage(ctx context.Context, input inputs.SendMessage, userId uuid.UUID) (*entity.Message, error) {
	message, err := usecase.messageRepository.CreateMessage(ctx, entity.Message{Content: input.Content, ChatId: input.ChatID, UserId: &userId})
	if err != nil {
		go log.Error().Msgf(err.Error())
		return nil, err
	}
	chat, err := usecase.chatRepository.GetChatById(ctx, message.ChatId.String())
	if err != nil {
		go log.Error().Msgf(err.Error())
		return nil, err
	}
	var otherUserId uuid.UUID
	if chat.FirstUserId.String() != userId.String() {
		otherUserId = *chat.FirstUserId
	} else {
		otherUserId = *chat.SecondUserId
	}
	// Publish the message on the redis channel corresponding to the chat
	err = usecase.redis.Publish(ctx, input.ChatID.String(), message).Err()
	if err != nil {
		panic(err)
	}
	// Publish the message on the redis channel corresponding to the user
	err = usecase.redis.Publish(ctx, otherUserId.String(), message).Err()
	if err != nil {
		panic(err)
	}
	return message, nil
}
