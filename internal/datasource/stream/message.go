package stream

import (
	"context"
	"encoding/json"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type MessageStreamDatasource interface {
	PublishMessage(ctx context.Context, message *entity.Message, userId string) error
	SubscribeByChat(ctx context.Context, chatId string) (chan *entity.Message, error)
	SubscribeByUser(ctx context.Context, userId string) (chan *entity.Message, error)
}

type messageStreamDatasource struct {
	redis *redis.Client
}

func NewMessageStreamDatasource(redis *redis.Client) MessageStreamDatasource {
	return &messageStreamDatasource{
		redis: redis,
	}
}

func (ds *messageStreamDatasource) PublishMessage(ctx context.Context, message *entity.Message, userId string) error {
	// Publish the message on the redis channel corresponding to the chat
	err := ds.redis.Publish(ctx, message.ChatId, message).Err()
	if err != nil {
		panic(err)
	}
	// Publish the message on the redis channel corresponding to the user
	err = ds.redis.Publish(ctx, userId, message).Err()
	if err != nil {
		panic(err)
	}
	return nil
}

func (ds *messageStreamDatasource) SubscribeByChat(ctx context.Context, chatId string) (chan *entity.Message, error) {
	ch := make(chan *entity.Message)

	go func() {
		// There is no error because go-redis automatically reconnects on error.
		pubsub := ds.redis.Subscribe(ctx, chatId)

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

func (ds *messageStreamDatasource) SubscribeByUser(ctx context.Context, userId string) (chan *entity.Message, error) {
	ch := make(chan *entity.Message)

	go func() {
		// There is no error because go-redis automatically reconnects on error.
		pubsub := ds.redis.Subscribe(ctx, userId)

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
