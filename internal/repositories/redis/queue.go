package redis

import (
	"context"
	"fmt"

	"github.com/iskanye/mirea-queue/internal/models"
	"github.com/redis/go-redis/v9"
)

func (s *Storage) Push(
	ctx context.Context,
	queue models.Queue,
	entry models.QueueEntry,
) error {
	const op = "redis.Push"

	_, err := s.cl.RPush(ctx, queue.Key(), entry.Student).Result()
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
		return models.QueueEntry{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.QueueEntry{
		Student: student,
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
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	entries := make([]models.QueueEntry, 0, n)
	for _, student := range students {
		entries = append(entries, models.QueueEntry{
			Student: student,
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

	pos, err := s.cl.LPos(ctx, queue.Key(), entry.Student, redis.LPosArgs{}).Result()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return pos + 1, nil
}
