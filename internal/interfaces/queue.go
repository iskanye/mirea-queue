package interfaces

import "github.com/iskanye/mirea-queue/internal/models"

type Queue interface {
	Push(entry models.QueueEntry) error
	Pop() (models.QueueEntry, error)
	Clear() error
}
