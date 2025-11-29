package interfaces

import (
	"context"

	"github.com/iskanye/mirea-queue/internal/models"
)

type Queue interface {
	Push(
		ctx context.Context,
		queue models.Queue,
		entry models.QueueEntry,
	) error
	Pop(
		ctx context.Context,
		queue models.Queue,
	) (models.QueueEntry, error)
	Clear(
		ctx context.Context,
		queue models.Queue,
	) error
}

type QueueViewer interface {
	Range(
		ctx context.Context,
		queue models.Queue,
		n int64,
	) ([]models.QueueEntry, error)
}

type QueuePosition interface {
	GetCurrentPosition(
		ctx context.Context,
		queue models.Queue,
		entry models.QueueEntry,
	) (int, error)
}
