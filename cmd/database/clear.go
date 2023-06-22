/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package database

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
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

		// Set connection pool configuration options
		config, err := pgxpool.ParseConfig(cfg.PostgresqlUrl)
		if err != nil {
			panic(err)
		}

		config.MaxConns = 20                     // Set the maximum number of connections in the pool
		config.MaxConnIdleTime = time.Minute * 5 // Set the maximum idle time for connections

		// pgx
		conn, err := pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			go log.Fatal().Msgf("Pgx connector Error: %v", err)
		}

		defer conn.Close()

		rows, err := conn.Query(context.Background(), "DELETE FROM \"user\";")
		if err != nil {
			go log.Fatal().Msg(err.Error())
		}
		defer rows.Close()
		rows, err = conn.Query(context.Background(), "DELETE FROM \"api_key\";")
		if err != nil {
			go log.Fatal().Msg(err.Error())
		}
		defer rows.Close()
		rows, err = conn.Query(context.Background(), "DELETE FROM \"refresh_token\";")
		if err != nil {
			go log.Fatal().Msg(err.Error())
		}
		defer rows.Close()
		rows, err = conn.Query(context.Background(), "DELETE FROM \"access_token\";")
		if err != nil {
			go log.Fatal().Msg(err.Error())
		}
		defer rows.Close()
		rows, err = conn.Query(context.Background(), "DELETE FROM \"message\";")
		if err != nil {
			go log.Fatal().Msg(err.Error())
		}
		defer rows.Close()
		rows, err = conn.Query(context.Background(), "DELETE FROM \"chat\";")
		if err != nil {
			go log.Fatal().Msg(err.Error())
		}
		defer rows.Close()
		go log.Info().Msg("Database cleaned!")
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
