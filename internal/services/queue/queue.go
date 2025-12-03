package queue

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/interfaces"
	"github.com/iskanye/mirea-queue/internal/models"
	"github.com/iskanye/mirea-queue/internal/repositories"
	"github.com/iskanye/mirea-queue/internal/services"
)

type Queue struct {
	log *slog.Logger

	// Пагинация очереди
	queueRange int64

	queue       interfaces.Queue
	queueViewer interfaces.QueueViewer
	queuePos    interfaces.QueuePosition
	queueLength interfaces.QueueLength
}

func New(
	log *slog.Logger,
	queueRange int64,
	queue interfaces.Queue,
	queueViewer interfaces.QueueViewer,
	queuePos interfaces.QueuePosition,
	queueLength interfaces.QueueLength,
) *Queue {
	return &Queue{
		log: log,

		queueRange: queueRange,

		queue:       queue,
		queueViewer: queueViewer,
		queuePos:    queuePos,
		queueLength: queueLength,
	}
}

func (q *Queue) Push(
	ctx context.Context,
	queue models.Queue,
	entry models.QueueEntry,
) (int64, error) {
	const op = "queue.Push"

	log := q.log.With(
		slog.String("op", op),
		slog.String("queue_group", queue.Group),
		slog.String("queue_subject", queue.Subject),
	)

	log.Info("Trying to push to queue")

	err := q.queue.Push(ctx, queue, entry)
	if err != nil {
		log.Error("Failed to push",
			slog.String("err", err.Error()),
		)

		if errors.Is(err, repositories.ErrAlreadyInQueue) {
			return 0, fmt.Errorf("%s: %w", op, services.ErrAlreadyInQueue)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	pos, err := q.queuePos.GetPosition(ctx, queue, entry)
	if err != nil {
		log.Error("Failed to get entry position",
			slog.String("err", err.Error()),
		)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully pushed")

	return pos, nil
}

func (q *Queue) Pop(
	ctx context.Context,
	queue models.Queue,
) (models.QueueEntry, error) {
	const op = "queue.Pop"

	log := q.log.With(
		slog.String("op", op),
		slog.String("queue_group", queue.Group),
		slog.String("queue_subject", queue.Subject),
	)

	log.Info("Trying to pop from queue")

	entry, err := q.queue.Pop(ctx, queue)
	if err != nil {
		log.Error("Failed to pop",
			slog.String("err", err.Error()),
		)

		if errors.Is(err, repositories.ErrNotFound) {
			return models.QueueEntry{}, fmt.Errorf("%s: %w", op, services.ErrNotFound)
		}
		return models.QueueEntry{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully poped")

	return entry, nil
}

func (q *Queue) Clear(
	ctx context.Context,
	queue models.Queue,
	key string,
) error {
	const op = "queue.Clear"

	log := q.log.With(
		slog.String("op", op),
		slog.String("queue_group", queue.Group),
		slog.String("queue_subject", queue.Subject),
	)

	log.Info("Trying to clear queue")

	err := q.queue.Clear(ctx, queue)
	if err != nil {
		log.Error("Failed to clear queue",
			slog.String("err", err.Error()),
		)

		if errors.Is(err, repositories.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, services.ErrNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully cleared")

	return nil
}

func (q *Queue) Pos(
	ctx context.Context,
	queue models.Queue,
	entry models.QueueEntry,
) (int64, error) {
	const op = "queue.Pos"

	log := q.log.With(
		slog.String("op", op),
		slog.String("queue_group", queue.Group),
		slog.String("queue_subject", queue.Subject),
	)

	log.Info("Trying to get position in queue")

	pos, err := q.queuePos.GetPosition(ctx, queue, entry)
	if err != nil {
		log.Error("Failed to get entry position",
			slog.String("err", err.Error()),
		)

		if errors.Is(err, repositories.ErrNotFound) {
			return 0, fmt.Errorf("%s: %w", op, services.ErrNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully got position")

	return pos, nil
}

func (q *Queue) Len(
	ctx context.Context,
	queue models.Queue,
) (int64, error) {
	const op = "queue.Len"

	log := q.log.With(
		slog.String("op", op),
		slog.String("queue_group", queue.Group),
		slog.String("queue_subject", queue.Subject),
	)

	log.Info("Trying to get length of queue")

	len, err := q.queueLength.Len(ctx, queue)
	if err != nil {
		log.Error("Failed to get length of queue",
			slog.String("err", err.Error()),
		)

		if errors.Is(err, repositories.ErrNotFound) {
			return 0, fmt.Errorf("%s: %w", op, services.ErrNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully got length")

	return len, nil
}
