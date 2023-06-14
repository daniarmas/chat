package repository

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/datasource/cache"
	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
	myerror "github.com/daniarmas/chat/pkg/my_error"
	"github.com/daniarmas/chat/pkg/sqldatabase"
)

type ChatRepository interface {
	CreateChat(ctx context.Context, chat entity.Chat) (*entity.Chat, error)
	GetChatById(ctx context.Context, chatId string) (*entity.Chat, error)
	GetChat(ctx context.Context, userId string, otherUserId string) (*entity.Chat, error)
	GetChats(ctx context.Context, userId string, updateTimeCursor time.Time) ([]*entity.Chat, error)
}

type chatRepository struct {
	database  *sqldatabase.Sql
	chatCache cache.ChatCacheDatasource
}

func NewChatRepository(database *sqldatabase.Sql, chatCache cache.ChatCacheDatasource) ChatRepository {
	return &chatRepository{
		database:  database,
		chatCache: chatCache,
	}
}

func (repo chatRepository) CreateChat(ctx context.Context, chat entity.Chat) (*entity.Chat, error) {
	chatModel := models.ChatOrm{}
	chatModel.MapToChatGorm(&chat)
	result := repo.database.Gorm.Create(&chatModel)
	if result.Error != nil {
		return nil, result.Error
	}

	res := chatModel.MapFromChatGorm()
	return res, nil
}

func (repo chatRepository) GetChat(ctx context.Context, userId string, otherUserId string) (*entity.Chat, error) {
	var chat *models.ChatOrm
	result := repo.database.Gorm.Where(
		repo.database.Gorm.Where("first_user_id = ?", userId).Or("second_user_id = ?", userId),
	).Where(
		repo.database.Gorm.Where("first_user_id = ?", otherUserId).Or("second_user_id = ?", otherUserId),
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

func (repo chatRepository) GetChatById(ctx context.Context, chatId string) (*entity.Chat, error) {
	var chat *models.ChatOrm
	result := repo.database.Gorm.Where("id = ?", chatId).Take(&chat)
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

	// // Get chats from cache
	// chats, err := repo.chatCache.GetChats(ctx, userId, cursor)
	// if err != nil {
	// 	return nil, err
	// }

	// if chats != nil {
	// 	return chats, nil
	// }

	// Get chats from database
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
