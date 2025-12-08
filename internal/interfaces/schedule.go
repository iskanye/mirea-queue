package interfaces

import (
	"context"

	"github.com/iskanye/mirea-queue/internal/models"
)

type GroupProvider interface {
	GetGroup(
		ctx context.Context,
		group string,
	) (models.Group, error)
}

type ScheduleProvider interface {
	GetSchedule(
		ctx context.Context,
		group models.Group,
	) ([]string, error)
}
