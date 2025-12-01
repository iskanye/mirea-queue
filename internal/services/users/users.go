package users

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

func (q *Users) CreateUser(
	ctx context.Context,
	chatID int64,
	user models.User,
) (models.User, error) {
	const op = "users.NewUser"

	log := q.log.With(
		slog.String("op", op),
		slog.String("username", user.Name),
	)

	log.Info("Trying to create new user")

	err := q.userCreator.CreateUser(ctx, chatID, user)
	if err != nil {
		// Не проверяем на то, существует ли уже юзер или нет
		// Это проверка находится на уровне обработчиков бота
		log.Error("Failed to create user",
			slog.String("err", err.Error()),
		)
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully created new user")

	return user, nil
}

func (q *Users) RemoveUser(
	ctx context.Context,
	chatID int64,
) error {
	const op = "users.RemoveUser"

	log := q.log.With(
		slog.String("op", op),
	)

	log.Info("Trying to remove user")

	err := q.userRemover.RemoveUser(ctx, chatID)
	if err != nil {
		// Не проверяем на то, существует ли уже юзер или нет
		// Это проверка находится на уровне обработчиков бота
		log.Error("Failed to remove user",
			slog.String("err", err.Error()),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User successfully removed")

	return nil
}

func (q *Users) UpdateUser(
	ctx context.Context,
	chatID int64,
	user models.User,
) (models.User, error) {
	const op = "users.UpdateUser"

	log := q.log.With(
		slog.String("op", op),
		slog.String("username", user.Name),
	)

	log.Info("Trying to update user data")

	err := q.userModifier.UpdateUser(ctx, chatID, user)
	if err != nil {
		// Не проверяем на то, существует ли уже юзер или нет
		// Это проверка находится на уровне обработчиков бота
		log.Error("Failed to update user",
			slog.String("err", err.Error()),
		)
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully updated user data")

	return user, nil
}

func (q *Users) GetUser(
	ctx context.Context,
	chatID int64,
) (models.User, error) {
	const op = "users.GetUser"

	log := q.log.With(
		slog.String("op", op),
	)

	log.Info("Trying to get user")

	user, err := q.userProvider.GetUser(ctx, chatID)
	if err != nil {
		log.Error("Failed to get user",
			slog.String("err", err.Error()),
		)

		if errors.Is(err, repositories.ErrNotFound) {
			return models.User{}, fmt.Errorf("%s: %w", op, services.ErrNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully got user")

	return user, nil
}
