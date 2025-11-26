package redis

import (
	"context"

	"github.com/iskanye/mirea-queue/internal/config"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	cl *redis.Client
}

func New(cfg config.Config) (*Storage, error) {
	const op = "repositories.redis.New"

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Username:     cfg.Redis.User,
		Password:     cfg.Redis.Password,
		ReadTimeout:  cfg.Redis.Timeout,
		WriteTimeout: cfg.Redis.Timeout,
	})

	// Проверяем подключение к Redis
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return &Storage{
		cl: client,
	}, nil
}
