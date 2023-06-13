package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
	myerror "github.com/daniarmas/chat/pkg/my_error"
	"github.com/daniarmas/chat/pkg/sqldatabase"
)

type ChatRepository interface {
	CreateChat(ctx context.Context, chat entity.Chat) (*entity.Chat, error)
	GetChat(ctx context.Context, firstUserId string, secondUserId string) (*entity.Chat, error)
	GetChats(ctx context.Context, userId string, updateTimeCursor time.Time) ([]*entity.Chat, error)
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

func (repo chatRepository) GetChats(ctx context.Context, userId string, updateTimeCursor time.Time) ([]*entity.Chat, error) {
	var cursor time.Time
	if updateTimeCursor.IsZero() {
		cursor = time.Now().UTC()
	} else {
		cursor = updateTimeCursor
	}
	var chatsOrm []models.ChatOrm
	var chats []*entity.Chat
	result := repo.database.Gorm.Where(
		repo.database.Gorm.Where("first_user_id = ?", userId).Or("second_user_id = ?", userId),
	).Where("update_time < ?", cursor).Limit(11).Order("update_time DESC").Find(&chatsOrm)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, element := range chatsOrm {
		chats = append(chats, element.MapFromChatGorm())
	}
	return chats, nil
}
