package main

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/daniarmas/chat/config"
	"github.com/daniarmas/chat/graph"
	"github.com/daniarmas/chat/internal/repository"
	"github.com/daniarmas/chat/internal/usecases"
	"github.com/daniarmas/chat/middleware"
	"github.com/daniarmas/chat/pkg/sqldatabase"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func main() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cfg := config.NewConfig()

	db, err := sqldatabase.New(cfg)
	if err != nil {
		log.Fatal().Msgf("Postgres Error: %v", err)
	}

	defer db.Close()

	userRepository := repository.NewUserRepository(db)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db)
	accessTokenRepository := repository.NewAccessTokenRepository(db)
	messageRepository := repository.NewMessageRepository(db)

	authUsecase := usecases.NewAuthUsecase(userRepository, refreshTokenRepository, accessTokenRepository, cfg)
	messageUsecase := usecases.NewMessageUsecase(userRepository, messageRepository, cfg)

	router := chi.NewRouter()

	router.Use(middleware.AuthorizationMiddleware(*cfg))

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{AuthUsecase: authUsecase, MessageUsecase: messageUsecase}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.GraphqlPort)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}

}
