package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Token    string
	Postgres postgresConfig
	Redis    redisConfig
}

type postgresConfig struct {
	User     string `env:"POSTGRES_USERNAME"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	DBName   string `env:"POSTGRES_DB"`
}

type redisConfig struct {
	User     string        `env:"REDIS_USER"`
	Password string        `env:"REDIS_PASSWORD"`
	Host     string        `env:"REDIS_HOST"`
	Port     int           `env:"REDIS_PORT"`
	DBName   string        `env:"REDIS_DB"`
	Timeout  time.Duration `env:"REDIS_TIMEOUT"`
}

func MustLoadConfig() *Config {
	var config *Config

	err := cleanenv.ReadEnv(&config)
	if err != nil {
		panic(err)
	}

	return config
}
