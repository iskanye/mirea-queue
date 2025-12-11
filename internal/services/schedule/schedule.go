package schedule

import (
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/interfaces"
)

type Schedule struct {
	log *slog.Logger

	groupProvider    interfaces.GroupProvider
	scheduleProvider interfaces.SubjectsProvider
}

func New(
	log *slog.Logger,
	groupProvider interfaces.GroupProvider,
	scheduleProvider interfaces.SubjectsProvider,
) *Schedule {
	return &Schedule{
		log: log,

		groupProvider:    groupProvider,
		scheduleProvider: scheduleProvider,
	}
}
