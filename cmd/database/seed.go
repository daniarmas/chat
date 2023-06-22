/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package database

import (
	"context"
	"time"

	"github.com/daniarmas/chat/internal/config"
	"github.com/daniarmas/chat/internal/datasource/databaseds"
	"github.com/daniarmas/chat/internal/datasource/hashds"
	"github.com/daniarmas/chat/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
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
		var users = []*models.User{{Email: "user1@example.com", Password: "prueba1234", Username: "user1", Fullname: "User1"}, {Email: "user2@example.com", Password: "prueba1234", Username: "user2", Fullname: "User2"}, {Email: "admin@example.com", Password: "prueba1234", Username: "admin", Fullname: "Admin"}}
		hashDs := hashds.NewBcryptHash()
		userDs := databaseds.NewUser(conn, hashDs)
		userDs.BulkCreateUser(context.Background(), users)
		go log.Info().Msg("Database migrations complete!")
	},
}

func init() {
	DatabaseCmd.AddCommand(seedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// seedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// seedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
