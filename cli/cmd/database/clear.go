/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package database

import (
	"github.com/daniarmas/chat/config"
	"github.com/daniarmas/chat/pkg/sqldatabase"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
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
		if err := db.Gorm.Exec("DELETE FROM \"user\";").Error; err != nil {
			log.Fatal().Msg(err.Error())
		}
		if err := db.Gorm.Exec("DELETE FROM \"api_key\";").Error; err != nil {
			log.Fatal().Msg(err.Error())
		}
		if err := db.Gorm.Exec("DELETE FROM \"refresh_token\";").Error; err != nil {
			log.Fatal().Msg(err.Error())
		}
		if err := db.Gorm.Exec("DELETE FROM \"access_token\";").Error; err != nil {
			log.Fatal().Msg(err.Error())
		}
		if err := db.Gorm.Exec("DELETE FROM \"message\";").Error; err != nil {
			log.Fatal().Msg(err.Error())
		}
		log.Info().Msg("Database cleaned!")
	},
}

func init() {
	DatabaseCmd.AddCommand(clearCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clearCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clearCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
