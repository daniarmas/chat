package config

type Config struct {
	PostgresqlDsn           string `mapstructure:"POSTGRESQL_DSN"`
	JwtSecret               string `mapstructure:"JWT_SECRET"`
	RefreshTokenExpireHours int    `mapstructure:"REFRESH_TOKEN_EXPIRE_HOURS"`
	AccessTokenExpireHours  int    `mapstructure:"ACCESS_TOKEN_EXPIRE_HOURS"`
}
