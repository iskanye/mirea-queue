package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Token      string        `env:"BOT_TOKEN"`
	BotTimeout time.Duration `env:"BOT_TIMEOUT"`
	AdminToken string        `env:"ADMIN_TOKEN"`
	Postgres   postgresConfig
	Redis      redisConfig
}

type postgresConfig struct {
	User      string `env:"POSTGRES_USER"`
	Password  string `env:"POSTGRES_PASSWORD"`
	Host      string `env:"POSTGRES_HOST"`
	Port      int    `env:"POSTGRES_PORT"`
	DB        string `env:"POSTGRES_DB"`
	PoolConns int    `env:"POSTGRES_POOL_CONNS"`
}

type redisConfig struct {
	User     string        `env:"REDIS_USER"`
	Password string        `env:"REDIS_USER_PASSWORD"`
	Addr     string        `env:"REDIS_ADDR"`
	Timeout  time.Duration `env:"REDIS_TIMEOUT"`
}

func MustLoadConfig() *Config {
	var config Config

	err := cleanenv.ReadEnv(&config)
	if err != nil {
		panic(err)
	}

	return &config
}
