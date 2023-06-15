package sqldatabase

import (
	"database/sql"
	"os"
	"time"

	logg "log"

	"github.com/daniarmas/chat/internal/config"
	"github.com/rs/zerolog/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Sql struct {
	Gorm  *gorm.DB
	SqlDb *sql.DB
}

func New(cfg *config.Config) (*Sql, error) {
	// Starting a database
	newLogger := logger.New(
		logg.New(os.Stdout, "\r\n", logg.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Millisecond * 200,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)
	gorm, err := gorm.Open(postgres.Open(cfg.PostgresqlDsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	connect, err := gorm.DB()
	if err != nil {
		return nil, err
	}
	maxAttempts := 3
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		err = connect.Ping()
		if err == nil {
			break
		}
		log.Error().Err(err).Msg("")
		time.Sleep(time.Duration(attempts) * time.Second)
	}
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	return &Sql{
		SqlDb: connect,
		Gorm:  gorm,
	}, nil
}

// Close -.
func (g *Sql) Close() {
	if g.SqlDb != nil {
		g.SqlDb.Close()
		log.Info().Msg("Sql server connection closed!")
	}
}
