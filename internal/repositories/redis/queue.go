package redis

import (
	"context"
	"fmt"

	"github.com/iskanye/mirea-queue/internal/models"
)

func (s *Storage) Push(
	ctx context.Context,
	queue models.Queue,
	entry models.QueueEntry,
) error {
	const op = "redis.Push"

	queueKey := fmt.Sprintf("%s:%s", queue.Group, queue.Subject)

	_, err := s.cl.LPush(ctx, queueKey, entry.Student).Result()
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

	queueKey := fmt.Sprintf("%s:%s", queue.Group, queue.Subject)

	student, err := s.cl.LPop(ctx, queueKey).Result()
	if err != nil {
		return models.QueueEntry{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.QueueEntry{
		Student: student,
	}, nil
}

func (s *Storage) Clear(
	ctx context.Context,
	queue models.Queue,
) error {
	const op = "redis.Pop"

	queueKey := fmt.Sprintf("%s:%s", queue.Group, queue.Subject)

	_, err := s.cl.Del(ctx, queueKey).Result()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
