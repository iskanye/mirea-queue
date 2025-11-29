package interfaces

import (
	"context"

	"github.com/iskanye/mirea-queue/internal/models"
)

type QueueService interface {
	Push(
		ctx context.Context,
		queue models.Queue,
		entry models.QueueEntry,
	) ([]models.QueueEntry, error)
	Pop(
		ctx context.Context,
		queue models.Queue,
	) (models.QueueEntry, []models.QueueEntry, error)
	Clear(
		ctx context.Context,
		queue models.Queue,
		key string,
	) error
	GetCurrentPosition(
		ctx context.Context,
		queue models.Queue,
		entry models.QueueEntry,
	) (int, []models.QueueEntry, error)
}

type UsersService interface {
	CreateUser(
		ctx context.Context,
		chatID int64,
		user models.User,
	) (models.User, error)
	RemoveUser(
		ctx context.Context,
		chatID int64,
	) error
	UpdateUser(
		ctx context.Context,
		chatID int64,
		user models.User,
	) (models.User, error)
	GetUser(
		ctx context.Context,
		chatID int64,
	) (models.User, error)
}
