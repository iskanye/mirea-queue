package queue

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/interfaces"
	"github.com/iskanye/mirea-queue/internal/models"
)

type Queue struct {
	log *slog.Logger

	// Пагинация очереди
	queueRange int64

	queue       interfaces.Queue
	queueViewer interfaces.QueueViewer
	queuePos    interfaces.QueuePosition
}

func New(
	log *slog.Logger,
	queueRange int64,
	queue interfaces.Queue,
	queueViewer interfaces.QueueViewer,
	queuePos interfaces.QueuePosition,
) *Queue {
	return &Queue{
		log: log,

		queueRange: queueRange,

		queue:       queue,
		queueViewer: queueViewer,
		queuePos:    queuePos,
	}
}

func (q *Queue) Push(
	ctx context.Context,
	queue models.Queue,
	entry models.QueueEntry,
) (int64, error) {
	const op = "Push"

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
	return models.QueueEntry{}, nil
}

func (q *Queue) Clear(
	ctx context.Context,
	queue models.Queue,
	key string,
) error {
	return nil
}

func (q *Queue) GetCurrentPosition(
	ctx context.Context,
	queue models.Queue,
	entry models.QueueEntry,
) (int64, error) {
	return 0, nil
}
