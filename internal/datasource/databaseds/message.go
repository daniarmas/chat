package databaseds

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
	"github.com/daniarmas/chat/pkg/sqldatabase"
)

type MessageDbDatasource interface {
	CreateMessage(ctx context.Context, message entity.Message) (*entity.Message, error)
	GetMessagesByChatId(ctx context.Context, chatId string, createTimeCursor time.Time) ([]*entity.Message, error)
}

type messageDbDatasource struct {
	database *sqldatabase.Sql
}

func NewMessage(database *sqldatabase.Sql) MessageDbDatasource {
	return &messageDbDatasource{
		database: database,
	}
}

func (repo messageDbDatasource) GetMessagesByChatId(ctx context.Context, chatId string, createTimeCursor time.Time) ([]*entity.Message, error) {
	var cursor time.Time
	if createTimeCursor.IsZero() {
		cursor = time.Now().UTC()
	} else {
		cursor = createTimeCursor
	}
	var messagesOrm []models.MessageOrm
	var messages []*entity.Message
	result := repo.database.Gorm.Where("chat_id = ?", chatId).Where("create_time < ?", cursor).Limit(11).Order("create_time DESC").Find(&messagesOrm)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, element := range messagesOrm {
		messages = append(messages, element.MapFromMessageGorm())
	}
	return messages, nil
}

func (repo messageDbDatasource) CreateMessage(ctx context.Context, message entity.Message) (*entity.Message, error) {
	messageModel := models.MessageOrm{}
	messageModel.MapToMessageGorm(&message)
	result := repo.database.Gorm.Create(&messageModel)
	if result.Error != nil {
		return nil, result.Error
	}

	res := messageModel.MapFromMessageGorm()
	return res, nil
}
