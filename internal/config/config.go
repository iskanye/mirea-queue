package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Token    string
	Postgres postgresConfig
}

type postgresConfig struct {
	User     string `env:"POSTGRES_USERNAME"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	DBName   string `env:"POSTGRES_DB"`
}

func MustLoadConfig() *Config {
	var config *Config

	err := cleanenv.ReadEnv(&config)
	if err != nil {
		panic(err)
	}

	return config
}
