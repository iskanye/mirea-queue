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
}

func New(
	log *slog.Logger,
	queueRange int64,
	queue interfaces.Queue,
	queueViewer interfaces.QueueViewer,
) *Queue {
	return &Queue{
		log: log,

		queueRange: queueRange,

		queue:       queue,
		queueViewer: queueViewer,
	}
}

// Пушает пользователя в очередь
func (q *Queue) Push(
	ctx context.Context,
	queue models.Queue,
	entry models.QueueEntry,
) ([]models.QueueEntry, error) {
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
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	queueEntries, err := q.queueViewer.Range(ctx, queue, q.queueRange)
	if err != nil {
		log.Error("Failed to get queue",
			slog.String("err", err.Error()),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully pushed")

	return queueEntries, nil
}

// Получает пользователя из начала очереди
// и удаляет его из очереди
func (q *Queue) Pop(
	ctx context.Context,
	queue models.Queue,
) (models.QueueEntry, []models.QueueEntry, error) {
	return models.QueueEntry{}, nil, nil
}

// Очищает очередь
func (q *Queue) Clear(
	ctx context.Context,
	queue models.Queue,
	key string,
) error {
	return nil
}

// Получает текущую позицию пользователя в очереди
func (q *Queue) GetCurrentPosition(
	ctx context.Context,
	queue models.Queue,
	entry models.QueueEntry,
) (int, error) {
	return 0, nil
}
