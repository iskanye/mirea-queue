package queue

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/models"
	"github.com/iskanye/mirea-queue/internal/repositories"
	"github.com/iskanye/mirea-queue/internal/services"
)

func (q *Queue) SaveToCache(
	ctx context.Context,
	chatID int64,
	queue models.Queue,
) error {
	const op = "queue.SaveToCache"

	log := q.log.With(
		slog.String("op", op),
		slog.String("queue_group", queue.Group),
		slog.String("queue_subject", queue.Subject),
	)

	log.Info("Trying to save queue to cache")

	err := q.cache.Set(ctx, fmt.Sprint(chatID), queue.Key())
	if err != nil {
		log.Error("Failed to save queue",
			slog.String("err", err.Error()),
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully saved")

	return nil
}

func (q *Queue) GetFromCache(
	ctx context.Context,
	chatID int64,
) (models.Queue, error) {
	const op = "queue.GetFromCache"

	log := q.log.With(
		slog.String("op", op),
	)

	log.Info("Trying to get queue from cache")

	queueKey, err := q.cache.Get(ctx, fmt.Sprint(chatID))
	if err != nil {
		log.Error("Failed to get queue",
			slog.String("err", err.Error()),
		)

		if errors.Is(err, repositories.ErrCacheMiss) {
			return models.Queue{}, fmt.Errorf("%s: %w", op, services.ErrNotFound)
		}
		return models.Queue{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully got queue")

	return models.QueueFromKey(queueKey), nil
}
