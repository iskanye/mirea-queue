package interfaces

import (
	"context"

	"github.com/iskanye/mirea-queue/internal/models"
)

type Queue interface {
	Push(
		ctx context.Context,
		entry models.QueueEntry,
	) error

	Pop(ctx context.Context) (models.QueueEntry, error)

	Clear(ctx context.Context) error
}
