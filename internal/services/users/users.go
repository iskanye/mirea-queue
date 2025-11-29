package users

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/interfaces"
	"github.com/iskanye/mirea-queue/internal/models"
)

type Users struct {
	log *slog.Logger

	userCreator  interfaces.UserCreator
	userRemover  interfaces.UserRemover
	userModifier interfaces.UserModifier
	userProvider interfaces.UserProvider
}

func New(
	log *slog.Logger,

	userCreator interfaces.UserCreator,
	userRemover interfaces.UserRemover,
	userModifier interfaces.UserModifier,
	userProvider interfaces.UserProvider,
) *Users {
	return &Users{
		log: log,

		userCreator:  userCreator,
		userRemover:  userRemover,
		userModifier: userModifier,
		userProvider: userProvider,
	}
}

func (q *Users) NewUser(
	ctx context.Context,
	chatID int64,
	user models.User,
) error {
	const op = "NewUser"

	log := q.log.With(
		slog.String("op", op),
		slog.String("username", user.Name),
	)

	log.Info("Trying to create new user")

	err := q.userCreator.CreateUser(ctx, chatID, user)
	if err != nil {
		log.Error("Failed to create user")
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully created new user")

	return nil
}
