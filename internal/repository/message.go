package repository

import (
	"context"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
	"github.com/daniarmas/chat/pkg/sqldatabase"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message entity.Message) (*entity.Message, error)
}

type message struct {
	database *sqldatabase.Sql
}

func NewMessageRepository(database *sqldatabase.Sql) MessageRepository {
	return &message{
		database: database,
	}
}

func (repo message) CreateMessage(ctx context.Context, message entity.Message) (*entity.Message, error) {
	messageModel := models.MessageOrm{}
	messageModel.MapToMessageGorm(&message)
	result := repo.database.Gorm.Create(&messageModel)
	if result.Error != nil {
		return nil, result.Error
	}

	res := messageModel.MapFromMessageGorm()
	return res, nil
}
