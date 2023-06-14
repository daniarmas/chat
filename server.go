package main

import (
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/daniarmas/chat/config"
	"github.com/daniarmas/chat/graph"
	"github.com/daniarmas/chat/internal/datasource/cache"
	"github.com/daniarmas/chat/internal/datasource/dbdatasource"
	"github.com/daniarmas/chat/internal/repository"
	"github.com/daniarmas/chat/internal/usecases"
	"github.com/daniarmas/chat/middleware"
	ownredis "github.com/daniarmas/chat/pkg/own-redis"
	"github.com/daniarmas/chat/pkg/sqldatabase"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
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

	redis, err := ownredis.NewRedis(cfg)
	if err != nil {
		log.Fatal().Msgf("Redis Error: %v", err)
	}

	// Datasources
	chatDbDatasource := dbdatasource.NewChatDbDatasource(db)
	accessTokenDbDatasource := dbdatasource.NewAccessTokenDbDatasource(db)
	chatCacheDatasource := cache.NewChatCacheDatasource(redis)

	userRepository := repository.NewUserRepository(db)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db)
	accessTokenRepository := repository.NewAccessTokenRepository(accessTokenDbDatasource)
	messageRepository := repository.NewMessageRepository(db)
	chatRepository := repository.NewChatRepository(chatCacheDatasource, chatDbDatasource)

	authUsecase := usecases.NewAuthUsecase(userRepository, refreshTokenRepository, accessTokenRepository, cfg)
	messageUsecase := usecases.NewMessageUsecase(userRepository, messageRepository, chatRepository, cfg, redis)
	chatUsecase := usecases.NewChatUsecase(chatRepository)

	router := chi.NewRouter()

	// CORS setup, allow any for now
	// https://gqlgen.com/recipes/cors/
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	})

	router.Use(middleware.AuthorizationMiddleware(*cfg))

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{AuthUsecase: authUsecase, MessageUsecase: messageUsecase, ChatUsecase: chatUsecase}}))

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
			return middleware.AuthorizationWebsocketMiddleware(ctx, cfg, initPayload)
		},
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", c.Handler(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.GraphqlPort)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}

}
