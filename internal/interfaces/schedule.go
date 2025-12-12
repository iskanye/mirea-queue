package interfaces

import (
	"context"

	"github.com/iskanye/mirea-queue/internal/models"
)

type GroupProvider interface {
	GetGroups(
		ctx context.Context,
		group string,
		limit int,
	) ([]models.Group, error)
}

type SubjectsProvider interface {
	GetSubjects(
		ctx context.Context,
		group models.Group,
	) ([]string, error)
}
