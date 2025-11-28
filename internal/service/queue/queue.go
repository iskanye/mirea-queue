package queue

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/models"
)

func (q *QueueService) Push(
	ctx context.Context,
	queue models.Queue,
	entry models.QueueEntry,
) error {
	const op = "Push"

	log := q.log.With(
		slog.String("op", op),
		slog.String("queue_group", queue.Group),
		slog.String("queue_subject", queue.Subject),
	)

	log.Info("Trying to push user to queue")

	err := q.queue.Push(ctx, queue, entry)
	if err != nil {
		log.Error("Failed to create user")
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully pushed")

	return nil
}
