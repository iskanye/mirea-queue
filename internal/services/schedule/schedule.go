package schedule

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/client/schedule"
	"github.com/iskanye/mirea-queue/internal/interfaces"
	"github.com/iskanye/mirea-queue/internal/models"
	"github.com/iskanye/mirea-queue/internal/services"
)

type Schedule struct {
	log *slog.Logger

	// Пагинация групп
	groupPagination int

	groupProvider    interfaces.GroupProvider
	subjectsProvider interfaces.SubjectsProvider
}

func New(
	log *slog.Logger,
	groupPagination int,
	groupProvider interfaces.GroupProvider,
	subjectsProvider interfaces.SubjectsProvider,
) *Schedule {
	return &Schedule{
		log: log,

		groupPagination: groupPagination,

		groupProvider:    groupProvider,
		subjectsProvider: subjectsProvider,
	}
}

func (s *Schedule) GetGroups(
	ctx context.Context,
	group string,
) ([]models.Group, error) {
	const op = "schedule.GetGroups"

	log := s.log.With(
		slog.String("op", op),
		slog.String("group_name", group),
	)

	log.Info("Trying to get groups list")

	groups, err := s.groupProvider.GetGroups(ctx, group, s.groupPagination)
	if err != nil {
		log.Error("Failed to get groups list",
			slog.String("err", err.Error()),
		)

		if errors.Is(err, schedule.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, services.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully got")

	return groups, nil
}

func (s *Schedule) GetSubjects(
	ctx context.Context,
	group models.Group,
) ([]string, error) {
	const op = "schedule.GetSubjects"

	log := s.log.With(
		slog.String("op", op),
		slog.String("group_name", group.Name),
	)

	log.Info("Trying to get group subjects")

	subjects, err := s.subjectsProvider.GetSubjects(ctx, group)
	if err != nil {
		log.Error("Failed to get group subjects",
			slog.String("err", err.Error()),
		)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully got")

	return subjects, nil
}
