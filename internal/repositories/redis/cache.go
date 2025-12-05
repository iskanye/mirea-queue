package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/iskanye/mirea-queue/internal/repositories"
	"github.com/redis/go-redis/v9"
)

func (s *Storage) Set(
	ctx context.Context,
	key string,
	val string,
) error {
	const op = "redis.Set"

	// Записываем данные в кеш
	_, err := s.cl.Set(ctx, key, val, 0).Result()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Get(
	ctx context.Context,
	key string,
) (string, error) {
	const op = "redis.Get"

	// Получаем данные из кеша
	val, err := s.cl.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("%s: %w", op, repositories.ErrCacheMiss)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return val, nil
}
