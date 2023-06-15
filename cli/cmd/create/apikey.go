/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package create

import (
	"context"

	"github.com/daniarmas/chat/config"
	"github.com/daniarmas/chat/internal/datasource/databaseds"
	"github.com/daniarmas/chat/internal/inputs"
	"github.com/daniarmas/chat/internal/repository"
	"github.com/daniarmas/chat/internal/usecases"
	"github.com/daniarmas/chat/pkg/sqldatabase"
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
		db, err := sqldatabase.New(cfg)
		if err != nil {
			log.Fatal().Msgf("Postgres Error: %v", err)
		}

		defer db.Close()
		apiKeyDbDatasource := databaseds.NewApiKey(db)
		apiKeyRepo := repository.NewApiKey(apiKeyDbDatasource)
		apiKeyUsecase := usecases.NewApiKey(apiKeyRepo)
		apiKey, err := apiKeyUsecase.CreateApiKey(ctx, inputs.CreateApiKeyInput{AppVersion: appVersion, Revoked: revoked})
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		log.Info().Msg("Api key created sucessfully.")
		log.Info().Msgf("Api key: %s", apiKey.Jwt)
	},
}

func init() {
	apikeyCmd.Flags().StringVarP(&appVersion, "appVersion", "a", "", "The appVersion of the apikey")
	if err := apikeyCmd.MarkFlagRequired("appVersion"); err != nil {
		log.Fatal().Msg(err.Error())
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
