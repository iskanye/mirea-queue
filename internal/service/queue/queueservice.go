package queue

import (
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/interfaces"
)

type QueueService struct {
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
) *QueueService {
	return &QueueService{
		log: log,

		queue:       queue,
		queueViewer: queueViewer,

		userCreator:  userCreator,
		userRemover:  userRemover,
		userModifier: userModifier,
		userProvider: userProvider,
	}
}
