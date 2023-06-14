/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package database

import (
	"github.com/daniarmas/chat/config"
	"github.com/daniarmas/chat/internal/models"
	"github.com/daniarmas/chat/pkg/sqldatabase"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		db, err := sqldatabase.New(cfg)
		if err != nil {
			log.Fatal().Msgf("Postgres Error: %v", err)
		}
		if err := db.Gorm.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
			log.Fatal().Msg(err.Error())
		}
		db.Gorm.AutoMigrate(&models.UserOrm{}, &models.ApiKeyOrm{}, &models.RefreshTokenOrm{}, &models.AccessTokenOrm{}, &models.MessageOrm{}, &models.ChatOrm{})
		log.Info().Msg("Database migrations complete!")
	},
}

func init() {
	DatabaseCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
