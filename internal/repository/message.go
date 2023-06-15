package repository

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/datasource/databaseds"
	"github.com/daniarmas/chat/internal/entity"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message entity.Message) (*entity.Message, error)
	GetMessagesByChatId(ctx context.Context, chatId string, createTimeCursor time.Time) ([]*entity.Message, error)
}

type messageRepository struct {
	messageDbDatasource databaseds.MessageDbDatasource
}

func NewMessage(messageDbDatasource databaseds.MessageDbDatasource) MessageRepository {
	return &messageRepository{
		messageDbDatasource: messageDbDatasource,
	}
}

func (repo messageRepository) GetMessagesByChatId(ctx context.Context, chatId string, createTimeCursor time.Time) ([]*entity.Message, error) {
	messages, err := repo.messageDbDatasource.GetMessagesByChatId(ctx, chatId, createTimeCursor)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (repo messageRepository) CreateMessage(ctx context.Context, message entity.Message) (*entity.Message, error) {
	res, err := repo.messageDbDatasource.CreateMessage(ctx, message)
	if err != nil {
		return nil, err
	}
	return res, nil
}
