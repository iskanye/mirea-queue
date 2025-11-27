package queue

import (
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/interfaces"
)

type Queue struct {
	log *slog.Logger

	queue       interfaces.Queue
	queueViewer interfaces.QueueViewer

	userCreator  interfaces.UserCreator
	userRemover  interfaces.UserRemover
	userModifier interfaces.UserModifier
	userProvider interfaces.UserProvider
}

func New(
	log *slog.Logger,
	queue interfaces.Queue,
	queueViewer interfaces.QueueViewer,
	userCreator interfaces.UserCreator,
	userRemover interfaces.UserRemover,
	userModifier interfaces.UserModifier,
	userProvider interfaces.UserProvider,
) *Queue {
	return &Queue{
		log: log,

		queue:       queue,
		queueViewer: queueViewer,

		userCreator:  userCreator,
		userRemover:  userRemover,
		userModifier: userModifier,
		userProvider: userProvider,
	}
}
