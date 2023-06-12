package repository

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
	"github.com/daniarmas/chat/pkg/sqldatabase"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message entity.Message) (*entity.Message, error)
	GetMessagesByChat(ctx context.Context, firstUserId string, secondUserId string, createTimeCursor time.Time) ([]*entity.Message, error)
}

type messageRepository struct {
	database *sqldatabase.Sql
}

func NewMessageRepository(database *sqldatabase.Sql) MessageRepository {
	return &messageRepository{
		database: database,
	}
}

func (repo messageRepository) GetMessagesByChat(ctx context.Context, firstUserId string, secondUserId string, createTimeCursor time.Time) ([]*entity.Message, error) {
	var cursor time.Time
	if createTimeCursor.IsZero() {
		cursor = time.Now().UTC()
	} else {
		cursor = createTimeCursor
	}
	var messagesOrm []models.MessageOrm
	var messages []*entity.Message
	result := repo.database.Gorm.Where(
		repo.database.Gorm.Where("sender_id = ?", firstUserId).Or("receiver_id = ?", firstUserId),
	).Where(
		repo.database.Gorm.Where("sender_id = ?", secondUserId).Or("receiver_id = ?", secondUserId),
	).Where("create_time < ?", cursor).Limit(11).Order("create_time DESC").Find(&messagesOrm)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, element := range messagesOrm {
		messages = append(messages, element.MapFromMessageGorm())
	}
	return messages, nil
}

func (repo messageRepository) CreateMessage(ctx context.Context, message entity.Message) (*entity.Message, error) {
	messageModel := models.MessageOrm{}
	messageModel.MapToMessageGorm(&message)
	result := repo.database.Gorm.Create(&messageModel)
	if result.Error != nil {
		return nil, result.Error
	}

	res := messageModel.MapFromMessageGorm()
	return res, nil
}
