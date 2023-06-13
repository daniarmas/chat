package repository

import (
	"context"
	"fmt"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
	myerror "github.com/daniarmas/chat/pkg/my_error"
	"github.com/daniarmas/chat/pkg/sqldatabase"
)

type ChatRepository interface {
	CreateChat(ctx context.Context, chat entity.Chat) (*entity.Chat, error)
	GetChat(ctx context.Context, firstUserId string, secondUserId string) (*entity.Chat, error)
}

type chatRepository struct {
	database *sqldatabase.Sql
}

func NewChatRepository(database *sqldatabase.Sql) ChatRepository {
	return &chatRepository{
		database: database,
	}
}

func (repo chatRepository) CreateChat(ctx context.Context, chat entity.Chat) (*entity.Chat, error) {
	chatModel := models.ChatOrm{}
	chat.Channel = fmt.Sprintf("%s:%s", chat.FirstUserId, chat.SecondUserId)
	chatModel.MapToChatGorm(&chat)
	result := repo.database.Gorm.Create(&chatModel)
	if result.Error != nil {
		return nil, result.Error
	}

	res := chatModel.MapFromChatGorm()
	return res, nil
}

func (repo chatRepository) GetChat(ctx context.Context, firstUserId string, secondUserId string) (*entity.Chat, error) {
	var chat *models.ChatOrm
	var channel1, channel2 string
	channel1 = fmt.Sprintf("%s:%s", firstUserId, secondUserId)
	channel2 = fmt.Sprintf("%s:%s", secondUserId, firstUserId)
	result := repo.database.Gorm.Where("channel = ?", channel1).Or("channel = ?", channel2).Take(&chat)
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
