package config

import "os"

type Config struct {
	Token       string
	PostgresUrl string
}

func MustLoadConfig() *Config {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		panic("cant get bot token")
	}

	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		panic("cant get database url")
	}

	return &Config{
		Token:       token,
		PostgresUrl: connStr,
	}
}
