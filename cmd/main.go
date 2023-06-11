package main

import (
	"github.com/daniarmas/chat/config"
	"github.com/daniarmas/chat/pkg/sqldatabase"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cfg := config.NewConfig()

	db, err := sqldatabase.New(cfg)
	if err != nil {
		log.Fatal().Msgf("Postgres Error: %v", err)
	}

	defer db.Close()
}
