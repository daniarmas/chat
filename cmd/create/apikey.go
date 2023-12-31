/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package create

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/config"
	"github.com/daniarmas/chat/internal/datasource/databaseds"
	"github.com/daniarmas/chat/internal/datasource/jwtds"
	"github.com/daniarmas/chat/internal/inputs"
	"github.com/daniarmas/chat/internal/repository"
	"github.com/daniarmas/chat/internal/usecases"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	appVersion string
	revoked    bool
)

// apikeyCmd represents the apikey command
var apikeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
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
			go log.Fatal().Msgf("Pgx connector Error: %v", err)
		}

		defer db.Close()

		jwtDs := jwtds.NewJwtDatasource(cfg)
		apiKeyDbDatasource := databaseds.NewApiKey(db)
		apiKeyRepo := repository.NewApiKey(apiKeyDbDatasource)
		apiKeyUsecase := usecases.NewApiKey(apiKeyRepo, jwtDs)
		apiKey, err := apiKeyUsecase.CreateApiKey(ctx, inputs.CreateApiKeyInput{AppVersion: appVersion, Revoked: revoked})
		if err != nil {
			go log.Error().Msg(err.Error())
		}
		go log.Info().Msg("Api key created sucessfully.")
		go log.Info().Msgf("Api key: %s", apiKey.Jwt)
	},
}

func init() {
	apikeyCmd.Flags().StringVarP(&appVersion, "appVersion", "a", "", "The appVersion of the apikey")
	if err := apikeyCmd.MarkFlagRequired("appVersion"); err != nil {
		go log.Fatal().Msg(err.Error())
	}
	apikeyCmd.Flags().BoolVarP(&revoked, "revoked", "r", false, "Revoke the apikey")

	CreateCmd.AddCommand(apikeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apikeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apikeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
