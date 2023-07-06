package stream

import (
	"context"
	"encoding/json"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type messageStreamNatsDatasource struct {
	nc *nats.Conn
}

func NewMessageStreamNatsDatasource(nc *nats.Conn) MessageStreamDatasource {
	return &messageStreamNatsDatasource{
		nc: nc,
	}
}

func (ds *messageStreamNatsDatasource) PublishMessage(ctx context.Context, message *entity.Message, userId string) error {
	// Publish the message on the redis channel corresponding to the chat
	// Serialize the struct to JSON
	data, _ := json.Marshal(message)
	err := ds.nc.Publish(message.ChatId, data)
	if err != nil {
		go log.Panic().Msg(err.Error())
	}
	// Publish the message on the redis channel corresponding to the user
	err = ds.nc.Publish(userId, data)
	if err != nil {
		go log.Panic().Msg(err.Error())
	}
	return nil
}

func (ds *messageStreamNatsDatasource) SubscribeByChat(ctx context.Context, chatId string) (chan *entity.Message, error) {
	ch := make(chan *entity.Message)

	go func() {
		ds.nc.Subscribe(chatId, func(m *nats.Msg) {
			var msg entity.Message
			err := json.Unmarshal(m.Data, &msg)
			if err != nil {
				go log.Error().Msgf("Error parsing message data: %s", err)
				return
			}
			ch <- &msg
		})

		ds.nc.Flush()
		if err := ds.nc.LastError(); err != nil {
			go log.Error().Msgf("Error subscribing to NATS subject: %s", err)
			return
		}

		// Listen for values on the channel and a close signal
		// for {
		// 	select {
		// 	case _, ok := <-ch:
		// 		if !ok {
		// 			sub.Unsubscribe()
		// 			return
		// 		}
		// 	}
		// }

	}()

	return ch, nil
}

func (ds *messageStreamNatsDatasource) SubscribeByUser(ctx context.Context, userId string) (chan *entity.Message, error) {
	ch := make(chan *entity.Message)

	go func() {
		ds.nc.Subscribe(userId, func(m *nats.Msg) {
			var msg entity.Message
			err := json.Unmarshal(m.Data, &msg)
			if err != nil {
				go log.Error().Msgf("Error parsing message data: %s", err)
				return
			}
			ch <- &msg
		})

		ds.nc.Flush()
		if err := ds.nc.LastError(); err != nil {
			go log.Error().Msgf("Error subscribing to NATS subject: %s", err)
			return
		}

		// Listen for values on the channel and a close signal
		// for {
		// 	select {
		// 	case _, ok := <-ch:
		// 		if !ok {
		// 			sub.Unsubscribe()
		// 			return
		// 		}
		// 	}
		// }

	}()

	return ch, nil
}
