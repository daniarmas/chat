package repository

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/datasource/databaseds"
	"github.com/daniarmas/chat/internal/datasource/stream"
	"github.com/daniarmas/chat/internal/entity"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message entity.Message) (*entity.Message, error)
	PublishMessage(ctx context.Context, message entity.Message, userId string) error
	GetMessagesByChatId(ctx context.Context, chatId string, createTimeCursor time.Time) ([]*entity.Message, error)
	ReceiveMessageByChat(ctx context.Context, chatId string) (chan *entity.Message, error)
	ReceiveMessageByUser(ctx context.Context, userId string) (chan *entity.Message, error)
}

type messageRepository struct {
	messageDbDatasource     databaseds.MessageDbDatasource
	messageStreamDatasource stream.MessageStreamDatasource
}

func NewMessage(messageDbDatasource databaseds.MessageDbDatasource, messageStreamDatasource stream.MessageStreamDatasource) MessageRepository {
	return &messageRepository{
		messageDbDatasource:     messageDbDatasource,
		messageStreamDatasource: messageStreamDatasource,
	}
}

func (repo messageRepository) PublishMessage(ctx context.Context, message entity.Message, userId string) error {
	err := repo.messageStreamDatasource.PublishMessage(ctx, &message, userId)
	if err != nil {
		return err
	}
	return nil
}

func (repo messageRepository) GetMessagesByChatId(ctx context.Context, chatId string, createTimeCursor time.Time) ([]*entity.Message, error) {
	res, err := repo.messageDbDatasource.GetMessagesByChatId(ctx, chatId, createTimeCursor)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo messageRepository) CreateMessage(ctx context.Context, message entity.Message) (*entity.Message, error) {
	res, err := repo.messageDbDatasource.CreateMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo messageRepository) ReceiveMessageByChat(ctx context.Context, chatId string) (chan *entity.Message, error) {
	res, err := repo.messageStreamDatasource.SubscribeByChat(ctx, chatId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo messageRepository) ReceiveMessageByUser(ctx context.Context, userId string) (chan *entity.Message, error) {
	res, err := repo.messageStreamDatasource.SubscribeByUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}
