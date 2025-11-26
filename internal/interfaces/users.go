package interfaces

import (
	"context"

	"github.com/iskanye/mirea-queue/internal/models"
)

type UserCreator interface {
	CreateUser(
		ctx context.Context,
		chatID int64,
		user models.User,
	) error
}

type UserRemover interface {
	RemoveUser(
		ctx context.Context,
		chatID int64,
	) error
}

type UserModifier interface {
	UpdateUser(
		ctx context.Context,
		chatID int64,
		user models.User,
	) error
}

type UserProvider interface {
	GetUser(
		ctx context.Context,
		chatID int64,
	) (models.User, error)
}
