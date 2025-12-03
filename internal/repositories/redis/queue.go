package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/iskanye/mirea-queue/internal/models"
	"github.com/iskanye/mirea-queue/internal/repositories"
	"github.com/redis/go-redis/v9"
)

func (s *Storage) Push(
	ctx context.Context,
	queue models.Queue,
	entry models.QueueEntry,
) error {
	const op = "redis.Push"

	// Пытаемся найти данный айди в очереди
	// Если есть, значит пользователь уже есть в очереди
	_, err := s.cl.LPop(ctx, queue.Key()).Result()
	if err == nil {
		return fmt.Errorf("%s: %w", op, repositories.ErrAlreadyInQueue)
	} else if !errors.Is(err, redis.Nil) {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.cl.RPush(ctx, queue.Key(), entry.ChatID).Result()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Pop(
	ctx context.Context,
	queue models.Queue,
) (models.QueueEntry, error) {
	const op = "redis.Pop"

	student, err := s.cl.LPop(ctx, queue.Key()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return models.QueueEntry{}, fmt.Errorf("%s: %w", op, repositories.ErrNotFound)
		}
		return models.QueueEntry{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.QueueEntry{
		ChatID: student,
	}, nil
}

func (s *Storage) Range(
	ctx context.Context,
	queue models.Queue,
	n int64,
) ([]models.QueueEntry, error) {
	const op = "redis.Range"

	students, err := s.cl.LRange(ctx, queue.Key(), 0, n-1).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("%s: %w", op, repositories.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	entries := make([]models.QueueEntry, 0, n)
	for _, student := range students {
		entries = append(entries, models.QueueEntry{
			ChatID: student,
		})
	}

	return entries, nil
}

func (s *Storage) Clear(
	ctx context.Context,
	queue models.Queue,
) error {
	const op = "redis.Pop"

	_, err := s.cl.Del(ctx, queue.Key()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("%s: %w", op, repositories.ErrNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetPosition(
	ctx context.Context,
	queue models.Queue,
	entry models.QueueEntry,
) (int64, error) {
	const op = "redis.GetPosition"

	pos, err := s.cl.LPos(ctx, queue.Key(), entry.ChatID, redis.LPosArgs{}).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, fmt.Errorf("%s: %w", op, repositories.ErrNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Отсчёт позиции должен начинаться с 1
	return pos + 1, nil
}

func (s *Storage) Len(
	ctx context.Context,
	queue models.Queue,
) (int64, error) {
	const op = "redis.Len"

	len, err := s.cl.LLen(ctx, queue.Key()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, fmt.Errorf("%s: %w", op, repositories.ErrNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return len, nil
}

func (s *Storage) LetAhead(
	ctx context.Context,
	queue models.Queue,
	entry models.QueueEntry,
) error {
	const op = "redis.LetAhead"

	pos, err := s.cl.LPos(ctx, queue.Key(), entry.ChatID, redis.LPosArgs{}).Result()
	if err != nil {
		// Списка нет или элемента нет в списке
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("%s: %w", op, repositories.ErrNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	ahead, err := s.cl.LIndex(ctx, queue.Key(), pos+1).Result()
	if err != nil {
		// Элемент не найден так как изначальный элемент в конце списка
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("%s: %w", op, repositories.ErrNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// Свапаем элементы
	_, err = s.cl.LSet(ctx, queue.Key(), pos, ahead).Result()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.cl.LSet(ctx, queue.Key(), pos+1, entry.ChatID).Result()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
