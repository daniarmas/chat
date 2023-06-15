package main

import (
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/daniarmas/chat/internal/config"
	"github.com/daniarmas/chat/internal/datasource/cacheds"
	"github.com/daniarmas/chat/internal/datasource/databaseds"
	"github.com/daniarmas/chat/internal/datasource/hashds"
	"github.com/daniarmas/chat/internal/datasource/jwtds"
	"github.com/daniarmas/chat/internal/delivery/graph"
	"github.com/daniarmas/chat/internal/delivery/graph/middleware"
	"github.com/daniarmas/chat/internal/repository"
	"github.com/daniarmas/chat/internal/usecases"
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
		go log.Fatal().Msgf("Postgres Error: %v", err)
	}

	defer db.Close()

	redis, err := ownredis.NewRedis(cfg)
	if err != nil {
		go log.Fatal().Msgf("Redis Error: %v", err)
	}

	// Database Datasources
	chatDatabaseDs := databaseds.NewChat(db)
	accessTokenDatabaseDs := databaseds.NewAccessToken(db)
	refreshTokenDatabaseDs := databaseds.NewRefreshToken(db)
	userDatabaseDs := databaseds.NewUser(db)
	messageDatabaseDs := databaseds.NewMessage(db)

	// Cache Datasources
	chatCacheDs := cacheds.NewChatCacheDatasource(redis)

	// Jwt Datasource
	jwtDs := jwtds.NewJwtDatasource(cfg)

	// Hash Datasource
	hashDs := hashds.NewBcryptHash()

	// Repositories
	userRepo := repository.NewUser(userDatabaseDs)
	refreshTokenRepo := repository.NewRefreshToken(refreshTokenDatabaseDs)
	accessTokenRepo := repository.NewAccessToken(accessTokenDatabaseDs)
	messageRepo := repository.NewMessage(messageDatabaseDs)
	chatRepo := repository.NewChat(chatCacheDs, chatDatabaseDs)

	// Usecases
	authUsecase := usecases.NewAuth(userRepo, refreshTokenRepo, accessTokenRepo, jwtDs, hashDs, cfg)
	messageUsecase := usecases.NewMessage(userRepo, messageRepo, chatRepo, cfg, redis)
	chatUsecase := usecases.NewChat(chatRepo)

	router := chi.NewRouter()

	// CORS setup, allow any for now
	// https://gqlgen.com/recipes/cors/
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	})

	router.Use(middleware.AuthorizationMiddleware(jwtDs))

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
			return middleware.AuthorizationWebsocketMiddleware(ctx, jwtDs, initPayload)
		},
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", c.Handler(srv))

	go log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.GraphqlPort)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}

}
