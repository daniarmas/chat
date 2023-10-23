/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package server

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/daniarmas/chat/gen"
	"github.com/daniarmas/chat/internal/config"
	"github.com/daniarmas/chat/internal/datasource/cacheds"
	"github.com/daniarmas/chat/internal/datasource/databaseds"
	"github.com/daniarmas/chat/internal/datasource/hashds"
	"github.com/daniarmas/chat/internal/datasource/jwtds"
	"github.com/daniarmas/chat/internal/datasource/stream"
	"github.com/daniarmas/chat/internal/delivery/graph"
	"github.com/daniarmas/chat/internal/delivery/graph/middleware"
	"github.com/daniarmas/chat/internal/repository"
	"github.com/daniarmas/chat/internal/usecases"
	ownredis "github.com/daniarmas/chat/pkg/own-redis"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// UNIX Time is faster and smaller than most timestamps
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

		cfg := config.NewConfig()

		// Set connection pool configuration options
		config, err := pgxpool.ParseConfig(cfg.PostgresqlUrl)
		if err != nil {
			panic(err)
		}

		config.MaxConns = 20                     // Set the maximum number of connections in the pool
		config.MaxConnIdleTime = time.Minute * 5 // Set the maximum idle time for connections

		// pgx
		db, err := pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			go log.Fatal().Msgf("Error connecting to PostgreSQL server: %v", err)
		}

		defer db.Close()

		// sqlc client
		dbsqlc, err := sql.Open("postgres", cfg.PostgresqlDsn)
		if err != nil {
			go log.Fatal().Msgf("Error connecting to postgres: %s", err)
		}

		// sqlc client
		sqlcQueries := gen.New(dbsqlc)

		// Redis client
		redis, err := ownredis.NewRedis(cfg)
		if err != nil {
			go log.Fatal().Msgf("Error connecting to Redis server: %s", err)
		}

		// NATS client
		nc, _ := nats.Connect(cfg.NatsUrl)
		if err != nil {
			go log.Fatal().Msgf("Error connecting to NATS server: %s", err)
			return
		}

		defer nc.Close()

		// Check if the connection is still active
		if !nc.IsConnected() {
			go log.Fatal().Msg("Connection to NATS server is lost")
			return
		}

		// Hash Datasource
		hashDs := hashds.NewBcryptHash()

		// Database Datasources
		chatDatabaseDs := databaseds.NewChat(db)
		accessTokenDatabaseDs := databaseds.NewAccessToken(db)
		refreshTokenDatabaseDs := databaseds.NewRefreshToken(db)
		userDatabaseDs := databaseds.NewUser(db, hashDs, sqlcQueries)
		messageDatabaseDs := databaseds.NewMessage(db)

		// Cache Datasources
		chatCacheDs := cacheds.NewChatCacheDatasource(redis)
		userCacheDs := cacheds.NewUserCacheDatasource(redis, cfg)
		accessTokenCacheDs := cacheds.NewAccessTokenCacheDatasource(redis, cfg)

		// Jwt Datasource
		jwtDs := jwtds.NewJwtDatasource(cfg)

		// Stream Datasource
		// messageStreamDatasource := stream.NewMessageStreamRedisDatasource(redis)
		messageStreamDatasource := stream.NewMessageStreamNatsDatasource(nc)

		// Repositories
		userRepo := repository.NewUser(userDatabaseDs, userCacheDs, sqlcQueries)
		refreshTokenRepo := repository.NewRefreshToken(refreshTokenDatabaseDs)
		accessTokenRepo := repository.NewAccessToken(accessTokenDatabaseDs, accessTokenCacheDs)
		messageRepo := repository.NewMessage(messageDatabaseDs, messageStreamDatasource)
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

		router.Use(middleware.AuthorizationMiddleware(jwtDs, accessTokenRepo))

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
	},
}

func init() {
	ServerCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clearCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clearCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
