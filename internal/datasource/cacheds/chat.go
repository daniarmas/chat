package cacheds

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/redis/go-redis/v9"
)

type ChatCacheDatasource interface {
	GetChats(ctx context.Context, userId string, updateTimeCursor time.Time) ([]*entity.Chat, error)
}

type chatRedisDatasource struct {
	redis *redis.Client
}

func NewChatCacheDatasource(redis *redis.Client) ChatCacheDatasource {
	return &chatRedisDatasource{
		redis: redis,
	}
}

func (repo chatRedisDatasource) GetChats(ctx context.Context, userId string, updateTimeCursor time.Time) ([]*entity.Chat, error) {
	cacheKey := fmt.Sprintf("%s:%s", userId, updateTimeCursor)
	chatFields, err := repo.redis.HGetAll(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	// convert the field-value pairs to chat objects
	var chats []*entity.Chat

	if len(chatFields) != 0 {
		for _, chatBytes := range chatFields {
			var chat entity.Chat
			err := json.Unmarshal([]byte(chatBytes), &chat)
			if err != nil {
				return nil, err
			}
			chats = append(chats, &chat)
		}
	} else {
		// Cache the data.
		// create a map of field-value pairs for each chat object

		for _, chat := range chats {
			chatBytes, err := json.Marshal(chat)
			if err != nil {
				return nil, err
			}
			chatFields[chat.ID] = string(chatBytes)
		}

		// set the chat objects in Redis using HMSet command
		err := repo.redis.HMSet(ctx, cacheKey, chatFields).Err()
		if err != nil {
			return nil, err
		}
	}

	return chats, nil
}
