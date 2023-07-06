package stream

import (
	"context"

	"github.com/daniarmas/chat/internal/entity"
)

type MessageStreamDatasource interface {
	PublishMessage(ctx context.Context, message *entity.Message, userId string) error
	SubscribeByChat(ctx context.Context, chatId string) (chan *entity.Message,  error)
	SubscribeByUser(ctx context.Context, userId string) (chan *entity.Message, error)
}
