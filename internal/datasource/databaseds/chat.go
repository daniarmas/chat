package databaseds

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/entity"
	myerror "github.com/daniarmas/chat/pkg/my_error"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type ChatDbDatasource interface {
	CreateChat(ctx context.Context, chat *entity.Chat) (*entity.Chat, error)
	GetChatById(ctx context.Context, chatId string) (*entity.Chat, error)
	GetChat(ctx context.Context, userId string, otherUserId string) (*entity.Chat, error)
	GetChats(ctx context.Context, userId string, updateTimeCursor time.Time) ([]*entity.Chat, error)
}

type chatPostgresDatasource struct {
	pgxConn  *pgxpool.Pool
}

func NewChat(pgxConn *pgxpool.Pool) ChatDbDatasource {
	return &chatPostgresDatasource{
		pgxConn:  pgxConn,
	}
}

func (data chatPostgresDatasource) CreateChat(ctx context.Context, chat *entity.Chat) (*entity.Chat, error) {
	var res entity.Chat
	createUpdateTime := time.Now().UTC()
	err := data.pgxConn.QueryRow(context.Background(), "INSERT INTO \"chat\" (first_user_id, second_user_id, create_time, update_time) VALUES ($1, $2, $3, $4) RETURNING id, first_user_id, second_user_id, create_time, update_time", chat.FirstUserId, chat.SecondUserId, createUpdateTime, createUpdateTime).Scan(&res.ID, &res.FirstUserId, &res.SecondUserId, &res.CreateTime, &res.UpdateTime)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &res, nil
}

func (data chatPostgresDatasource) GetChat(ctx context.Context, userId string, otherUserId string) (*entity.Chat, error) {
	var chat entity.Chat
	row := data.pgxConn.QueryRow(context.Background(), "SELECT id, first_user_id, second_user_id, create_time, update_time FROM \"chat\" WHERE (first_user_id = $1 OR second_user_id = $1) AND (first_user_id = $2 OR second_user_id = $2);", userId, otherUserId)
	err := row.Scan(&chat.ID, &chat.FirstUserId, &chat.SecondUserId, &chat.CreateTime, &chat.UpdateTime)
	if err != nil {
		if err.Error() == "no rows in result set" {
			log.Error().Msg(err.Error())
			return nil, myerror.NotFoundError{}
		} else {
			log.Error().Msg(err.Error())
			return nil, err
		}
	}
	return &chat, nil
}

func (data chatPostgresDatasource) GetChatById(ctx context.Context, chatId string) (*entity.Chat, error) {
	var chat entity.Chat
	row := data.pgxConn.QueryRow(context.Background(), "SELECT id, first_user_id, second_user_id, create_time, update_time FROM \"chat\" WHERE id = $1;", chatId)
	err := row.Scan(&chat.ID, &chat.FirstUserId, &chat.SecondUserId, &chat.CreateTime, &chat.UpdateTime)
	if err != nil {
		if err.Error() == "no rows in result set" {
			log.Error().Msg(err.Error())
			return nil, myerror.NotFoundError{}
		} else {
			log.Error().Msg(err.Error())
			return nil, err
		}
	}
	return &chat, nil
}

func (data chatPostgresDatasource) GetChats(ctx context.Context, userId string, updateTimeCursor time.Time) ([]*entity.Chat, error) {
	var cursor time.Time
	if updateTimeCursor.IsZero() {
		cursor = time.Now().UTC()
	} else {
		cursor = updateTimeCursor
	}

	var chats []*entity.Chat

	rows, err := data.pgxConn.Query(context.Background(), "SELECT id, first_user_id, second_user_id, create_time, update_time FROM chat WHERE (first_user_id = $1 OR second_user_id = $1) AND (update_time < $2) ORDER BY update_time DESC LIMIT 11;", userId, cursor)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var chat entity.Chat
		err := rows.Scan(&chat.ID, &chat.FirstUserId, &chat.SecondUserId, &chat.CreateTime, &chat.UpdateTime)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		chats = append(chats, &chat)
	}

	if err = rows.Err(); err != nil {
		log.Error().Msg(err.Error())
	}
	return chats, nil
}
