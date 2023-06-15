package databaseds

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
	myerror "github.com/daniarmas/chat/pkg/my_error"
	"github.com/daniarmas/chat/pkg/sqldatabase"
)

type ChatDbDatasource interface {
	CreateChat(ctx context.Context, chat *entity.Chat) (*entity.Chat, error)
	GetChatById(ctx context.Context, chatId string) (*entity.Chat, error)
	GetChat(ctx context.Context, userId string, otherUserId string) (*entity.Chat, error)
	GetChats(ctx context.Context, userId string, updateTimeCursor time.Time) ([]*entity.Chat, error)
}

type chatPostgresDatasource struct {
	database *sqldatabase.Sql
}

func NewChat(database *sqldatabase.Sql) ChatDbDatasource {
	return &chatPostgresDatasource{
		database: database,
	}
}

func (data chatPostgresDatasource) CreateChat(ctx context.Context, chat *entity.Chat) (*entity.Chat, error) {
	chatModel := models.ChatOrm{}
	chatModel.MapToChatGorm(chat)
	result := data.database.Gorm.Create(&chatModel)
	if result.Error != nil {
		return nil, result.Error
	}

	res := chatModel.MapFromChatGorm()
	return res, nil
}

func (data chatPostgresDatasource) GetChat(ctx context.Context, userId string, otherUserId string) (*entity.Chat, error) {
	var chat *models.ChatOrm
	result := data.database.Gorm.Where(
		data.database.Gorm.Where("first_user_id = ?", userId).Or("second_user_id = ?", userId),
	).Where(
		data.database.Gorm.Where("first_user_id = ?", otherUserId).Or("second_user_id = ?", otherUserId),
	).Take(&chat)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, myerror.NotFoundError{}
		} else {
			return nil, myerror.InternalServerError{}
		}
	}
	res := chat.MapFromChatGorm()
	return res, nil
}

func (data chatPostgresDatasource) GetChatById(ctx context.Context, chatId string) (*entity.Chat, error) {
	var chat *models.ChatOrm
	result := data.database.Gorm.Where("id = ?", chatId).Take(&chat)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, myerror.NotFoundError{}
		} else {
			return nil, myerror.InternalServerError{}
		}
	}
	res := chat.MapFromChatGorm()
	return res, nil
}

func (data chatPostgresDatasource) GetChats(ctx context.Context, userId string, updateTimeCursor time.Time) ([]*entity.Chat, error) {
	var cursor time.Time
	if updateTimeCursor.IsZero() {
		cursor = time.Now().UTC()
	} else {
		cursor = updateTimeCursor
	}
	var chatsOrm []models.ChatOrm
	var chats []*entity.Chat

	result := data.database.Gorm.Where(
		data.database.Gorm.Where("first_user_id = ?", userId).Or("second_user_id = ?", userId),
	).Where("update_time < ?", cursor).Limit(11).Order("update_time DESC").Find(&chatsOrm)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, element := range chatsOrm {
		chats = append(chats, element.MapFromChatGorm())
	}

	return chats, nil
}
