package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresqlDsn           string `mapstructure:"POSTGRESQL_DSN"`
	JwtSecret               string `mapstructure:"JWT_SECRET"`
	RefreshTokenExpireHours int    `mapstructure:"REFRESH_TOKEN_EXPIRE_HOURS"`
	AccessTokenExpireHours  int    `mapstructure:"ACCESS_TOKEN_EXPIRE_HOURS"`
	GraphqlPort             string `mapstructure:"GRAPHQL_PORT"`
}

func NewConfig() *Config {
	env := Config{}
	viper.AddConfigPath(".")
	viper.SetConfigFile("app.env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Msgf("Can't find the file .env : %v", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal().Msgf("Environment can't be loaded: %v", err)
	}

	return &env
}
