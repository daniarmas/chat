package databaseds

import (
	"context"
	"time"

	"github.com/daniarmas/chat/gen"
	"github.com/daniarmas/chat/internal/datasource/hashds"
	"github.com/daniarmas/chat/internal/entity"
	"github.com/daniarmas/chat/internal/models"
	myerror "github.com/daniarmas/chat/pkg/my_error"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type UserDbDatasource interface {
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, email string, password string, username string, fullname string) (*entity.User, error)
	BulkCreateUser(ctx context.Context, users []*models.User) ([]*models.User, error)
}

type userDbDatasource struct {
	pgxConn *pgxpool.Pool
	hashDs  hashds.HashDatasource
	queries *gen.Queries
}

func NewUser(pgxConn *pgxpool.Pool, hashDs hashds.HashDatasource,
	queries *gen.Queries) UserDbDatasource {
	return &userDbDatasource{
		pgxConn: pgxConn,
		hashDs:  hashDs,
		queries: queries,
	}
}

func (repo *userDbDatasource) GetUserById(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	row := repo.pgxConn.QueryRow(context.Background(), "SELECT id, email, fullname, username, password, create_time FROM \"user\" WHERE id = $1;", id)
	err := row.Scan(&user.ID, &user.Email, &user.Fullname, &user.Username, &user.Password, &user.CreateTime)
	if err != nil {
		if err.Error() == "no rows in result set" {
			log.Error().Msg(err.Error())
			return nil, myerror.NotFoundError{}
		} else {
			log.Error().Msg(err.Error())
			return nil, err
		}
	}
	return &user, nil
}

func (repo *userDbDatasource) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	res, err := repo.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &entity.User{
		ID:         res.ID.String(),
		Email:      res.Email,
		Password:   res.Password,
		Fullname:   res.Fullname,
		Username:   res.Username,
		CreateTime: res.CreateTime,
	}, nil
}

func (repo *userDbDatasource) CreateUser(ctx context.Context, email string, password string, username string, fullname string) (*entity.User, error) {
	var user entity.User
	passwordHashed, _ := repo.hashDs.Hash(password)
	err := repo.pgxConn.QueryRow(context.Background(), "INSERT INTO \"user\" (email, fullname, username, password, create_time) VALUES ($1, $2, $3) RETURNING id, email, fullname, username, password, create_time", email, fullname, passwordHashed, time.Now().UTC()).Scan(&user.ID, &user.Email, &user.Fullname, &user.Username, &user.Password, &user.CreateTime)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return &user, nil
}

func (repo *userDbDatasource) BulkCreateUser(ctx context.Context, users []*models.User) ([]*models.User, error) {
	var res []*models.User

	// Create a Batch object for the transaction
	batch := &pgx.Batch{}

	// Add an INSERT statement for each row of data
	for _, item := range users {
		passwordHashed, _ := repo.hashDs.Hash(item.Password)
		batch.Queue("INSERT INTO \"user\" (email, fullname, username, password, create_time) VALUES ($1, $2, $3, $4, $5) RETURNING id, email, fullname, username, password, create_time", item.Email, item.Fullname, item.Username, passwordHashed, time.Now().UTC())
	}

	// Execute the transaction and get the results
	results := repo.pgxConn.SendBatch(context.Background(), batch)
	defer results.Close()

	// Iterate over the results and scan the data into the struct
	for i := 0; i < len(users); i++ {
		var user models.User
		// Get the result of the statement and check for errors
		err := results.QueryRow().Scan(&user.ID, &user.Email, &user.Fullname, &user.Username, &user.Password, &user.CreateTime)
		if err != nil {
			log.Error().Msg(err.Error())
			return nil, err
		}
		res = append(res, &user)
	}

	return res, nil
}
