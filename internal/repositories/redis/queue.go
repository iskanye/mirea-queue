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
	_, err := s.cl.ZRank(ctx, queue.Key(), entry.ChatID).Result()
	if err == nil {
		return fmt.Errorf("%s: %w", op, repositories.ErrAlreadyInQueue)
	} else if !errors.Is(err, redis.Nil) {
		return fmt.Errorf("%s: %w", op, err)
	}

	if entry.Position == 0 {
		// Ищем последнюю позицию на которую можно поставить элемент
		entries, err := s.cl.ZRangeWithScores(ctx, queue.Key(), 0, -1).Result()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		pos := 1
		// Проходимся по списку позиций, пока не находим пустое место
		for _, e := range entries {
			if e.Score != float64(pos) {
				break
			}
			pos++
		}
		entry.Position = pos
	} else {
		// Проверяем что на данное место можно встать
		entry, err := s.cl.ZRangeByScore(ctx, queue.Key(), &redis.ZRangeBy{
			Min: fmt.Sprint(entry.Position),
			Max: fmt.Sprint(entry.Position),
		}).Result()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if len(entry) != 0 {
			// На данное место уже можно встать => неверно выбрана позиция
			return fmt.Errorf("%s: %w", op, repositories.ErrPlaceTaken)
		}
	}

	_, err = s.cl.ZAdd(ctx, queue.Key(), entry.ToRedis()).Result()
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

	student, err := s.cl.ZPopMin(ctx, queue.Key()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return models.QueueEntry{}, fmt.Errorf("%s: %w", op, repositories.ErrNotFound)
		}
		return models.QueueEntry{}, fmt.Errorf("%s: %w", op, err)
	}

	// Уменьшаем позиции в очереди на 1
	entries, err := s.cl.ZRange(ctx, queue.Key(), 0, -1).Result()
	if err != nil {
		return models.QueueEntry{}, fmt.Errorf("%s: %w", op, err)
	}

	for _, entry := range entries {
		_, err := s.cl.ZIncrBy(ctx, queue.Key(), -1, entry).Result()
		if err != nil {
			return models.QueueEntry{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	return models.QueueEntry{
		ChatID: student[0].Member.(string),
	}, nil
}

func (s *Storage) Range(
	ctx context.Context,
	queue models.Queue,
	n int64,
) ([]models.QueueEntry, error) {
	const op = "redis.Range"

	students, err := s.cl.ZRangeWithScores(ctx, queue.Key(), 0, n-1).Result()
	// Очередь не создана
	if len(students) == 0 {
		return nil, fmt.Errorf("%s: %w", op, repositories.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	entries := make([]models.QueueEntry, 0, n)
	for _, student := range students {
		entries = append(entries, models.QueueEntry{
			Position: int(student.Score),
			ChatID:   student.Member.(string),
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

	pos, err := s.cl.ZScore(ctx, queue.Key(), entry.ChatID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, fmt.Errorf("%s: %w", op, repositories.ErrNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int64(pos), nil
}

func (s *Storage) Len(
	ctx context.Context,
	queue models.Queue,
) (int64, error) {
	const op = "redis.Len"

	len, err := s.cl.ZCard(ctx, queue.Key()).Result()
	if err != nil {
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

	pos, err := s.cl.ZRank(ctx, queue.Key(), entry.ChatID).Result()
	if err != nil {
		// Элемента нет в списке
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("%s: %w", op, repositories.ErrNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	ahead, err := s.cl.ZRangeByScore(ctx, queue.Key(), &redis.ZRangeBy{
		Min: fmt.Sprint(pos + 1),
		Max: fmt.Sprint(pos + 1),
	}).Result()
	if errors.Is(err, redis.Nil) {
		// Элемент не найден так как изначальный элемент в конце списка
		// или перед ним дырка => можем увеличить ранг без последствий
		_, err := s.cl.ZIncrBy(ctx, queue.Key(), 1, entry.ChatID).Result()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	} else if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Элемент спереди найден => меняем им ранги
	_, err = s.cl.ZIncrBy(ctx, queue.Key(), 1, entry.ChatID).Result()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.cl.ZIncrBy(ctx, queue.Key(), -1, ahead[0]).Result()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Remove(
	ctx context.Context,
	queue models.Queue,
	entry models.QueueEntry,
) error {
	const op = "redis.Remove"

	_, err := s.cl.ZRem(ctx, queue.Key(), entry.ChatID).Result()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
