package queue

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/models"
)

func (q *QueueService) NewUser(
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
