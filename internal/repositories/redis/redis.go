package redis

import (
	"context"
	"fmt"

	"github.com/iskanye/mirea-queue/internal/config"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	cl *redis.Client
}

func New(cfg *config.Config) (*Storage, error) {
	const op = "redis.New"

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		ReadTimeout:  cfg.Redis.Timeout,
		WriteTimeout: cfg.Redis.Timeout,
	})

	// Проверяем подключение к Redis
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		cl: client,
	}, nil
}

func (s *Storage) Close() error {
	return s.cl.Close()
}

func (s *Storage) FlushDB(
	ctx context.Context,
) error {
	const op = "redis.Remove"

	_, err := s.cl.FlushDB(ctx).Result()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
