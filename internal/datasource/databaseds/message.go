package databaseds

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type MessageDbDatasource interface {
	CreateMessage(ctx context.Context, message entity.Message) (*entity.Message, error)
	GetMessagesByChatId(ctx context.Context, chatId string, createTimeCursor time.Time) ([]*entity.Message, error)
}

type messageDbDatasource struct {
	pgxConn  *pgxpool.Pool
}

func NewMessage(pgxConn *pgxpool.Pool) MessageDbDatasource {
	return &messageDbDatasource{
		pgxConn:  pgxConn,
	}
}

func (repo messageDbDatasource) GetMessagesByChatId(ctx context.Context, chatId string, createTimeCursor time.Time) ([]*entity.Message, error) {
	var cursor time.Time
	if createTimeCursor.IsZero() {
		cursor = time.Now().UTC()
	} else {
		cursor = createTimeCursor
	}

	var messages []*entity.Message

	rows, err := repo.pgxConn.Query(context.Background(), "SELECT id, chat_id, user_id, content, create_time FROM message WHERE create_time < $1 ORDER BY create_time DESC LIMIT 11;", cursor)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var message entity.Message
		err := rows.Scan(&message.ID, &message.ChatId, &message.UserId, &message.Content, &message.CreateTime)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		messages = append(messages, &message)
	}

	if err = rows.Err(); err != nil {
		log.Error().Msg(err.Error())
	}
	return messages, nil
}

func (repo messageDbDatasource) CreateMessage(ctx context.Context, message entity.Message) (*entity.Message, error) {
	var res entity.Message
	err := repo.pgxConn.QueryRow(context.Background(), "INSERT INTO \"message\" (chat_id, user_id, content, create_time) VALUES ($1, $2, $3, $4) RETURNING id, chat_id, user_id, content, create_time", message.ChatId, message.UserId, message.Content, time.Now().UTC()).Scan(&res.ID, &res.ChatId, &res.UserId, &res.Content, &res.CreateTime)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &res, nil
}
