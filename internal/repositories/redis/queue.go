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

	_, err := s.cl.LPush(ctx, queue.Key(), entry.Student).Result()
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
