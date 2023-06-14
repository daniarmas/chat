package repository

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/datasource/cache"
	"github.com/daniarmas/chat/internal/datasource/dbdatasource"
	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/pkg/sqldatabase"
)

type ChatRepository interface {
	CreateChat(ctx context.Context, chat *entity.Chat) (*entity.Chat, error)
	GetChatById(ctx context.Context, chatId string) (*entity.Chat, error)
	GetChat(ctx context.Context, userId string, otherUserId string) (*entity.Chat, error)
	GetChats(ctx context.Context, userId string, updateTimeCursor time.Time) ([]*entity.Chat, error)
}

type chatRepository struct {
	database         *sqldatabase.Sql
	chatDbDatasource dbdatasource.ChatDbDatasource
	chatCache        cache.ChatCacheDatasource
}

func NewChatRepository(database *sqldatabase.Sql, chatCache cache.ChatCacheDatasource, chatDbDatasource dbdatasource.ChatDbDatasource) ChatRepository {
	return &chatRepository{
		database:         database,
		chatDbDatasource: chatDbDatasource,
		chatCache:        chatCache,
	}
}

func (repo chatRepository) CreateChat(ctx context.Context, chat *entity.Chat) (*entity.Chat, error) {
	chat, err := repo.chatDbDatasource.CreateChat(ctx, chat)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (repo chatRepository) GetChat(ctx context.Context, userId string, otherUserId string) (*entity.Chat, error) {
	chat, err := repo.chatDbDatasource.GetChat(ctx, userId, otherUserId)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (repo chatRepository) GetChatById(ctx context.Context, chatId string) (*entity.Chat, error) {
	chat, err := repo.chatDbDatasource.GetChatById(ctx, chatId)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (repo chatRepository) GetChats(ctx context.Context, userId string, updateTimeCursor time.Time) ([]*entity.Chat, error) {
	chat, err := repo.chatDbDatasource.GetChats(ctx, userId, updateTimeCursor)
	if err != nil {
		return nil, err
	}
	return chat, nil
}
