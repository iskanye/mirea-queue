package queue

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/iskanye/mirea-queue/internal/interfaces"
	"github.com/iskanye/mirea-queue/internal/models"
	"github.com/iskanye/mirea-queue/internal/repositories"
	"github.com/iskanye/mirea-queue/internal/services"
	"github.com/iskanye/mirea-queue/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	queueRange = 10
	timeout    = time.Second
)

var (
	// НЕ ПУТАТЬ С service, queue - это мок интерфейса Queue,
	// для тестирования надо вызывать функции переменной service!
	queue        *mocks.MockQueue
	queueViewer  *mocks.MockQueueViewer
	queuePos     *mocks.MockQueuePosition
	queueLength  *mocks.MockQueueLength
	queueSwap    *mocks.MockQueueSwap
	queueRemover *mocks.MockQueueRemover
	cache        *mocks.MockCache
)

// Создаем сервис очереди с приклеплёнными к нему моками
func newService(t *testing.T) (interfaces.QueueService, context.Context) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	queue = mocks.NewMockQueue(t)
	queueViewer = mocks.NewMockQueueViewer(t)
	queuePos = mocks.NewMockQueuePosition(t)
	queueLength = mocks.NewMockQueueLength(t)
	queueSwap = mocks.NewMockQueueSwap(t)
	queueRemover = mocks.NewMockQueueRemover(t)
	cache = mocks.NewMockCache(t)

	t.Cleanup(func() {
		// Очищаем ожидаемые вызовы моков
		queue.ExpectedCalls = nil
		queueViewer.ExpectedCalls = nil
		queuePos.ExpectedCalls = nil
		queueLength.ExpectedCalls = nil
		queueSwap.ExpectedCalls = nil
		queueRemover.ExpectedCalls = nil
		cache.ExpectedCalls = nil
		cancel()
	})

	return New(
		slog.New(slog.DiscardHandler),
		queueRange,
		queue,
		queueViewer,
		queuePos,
		queueLength,
		queueSwap,
		queueRemover,
		cache,
	), ctx
}

// QueueService.Push

func TestQueuePush_Success(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	chatID := gofakeit.ID()
	entry := models.QueueEntry{
		ChatID: chatID,
	}

	expectedPos := int64(1)
	queue.EXPECT().Push(ctx, subjectQueue, entry).Return(nil)
	queuePos.EXPECT().GetPosition(ctx, subjectQueue, entry).Return(expectedPos, nil)

	pos, err := service.Push(ctx, subjectQueue, entry)
	require.Empty(t, err)
	assert.Equal(t, expectedPos, pos)
}

func TestQueuePush_AlreadyInQueue(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	chatID := gofakeit.ID()
	entry := models.QueueEntry{
		ChatID: chatID,
	}

	expectedPos := int64(1)
	queue.EXPECT().Push(ctx, subjectQueue, entry).Return(nil).Once()
	queuePos.EXPECT().GetPosition(ctx, subjectQueue, entry).Return(expectedPos, nil)

	pos, err := service.Push(ctx, subjectQueue, entry)
	require.Empty(t, err)
	assert.Equal(t, expectedPos, pos)

	// Переназначаем ожидаемое возвращаемое значение мока
	queue.EXPECT().Push(ctx, subjectQueue, entry).Return(repositories.ErrAlreadyInQueue)

	pos, err = service.Push(ctx, subjectQueue, entry)
	assert.Error(t, err)
	assert.ErrorIs(t, err, services.ErrAlreadyInQueue)
	assert.Equal(t, int64(0), pos)
}

// QueueService.Pop

func TestQueuePop_Success(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	chatID := gofakeit.ID()
	entry := models.QueueEntry{
		ChatID: chatID,
	}

	expectedPos := int64(1)
	queue.EXPECT().Push(ctx, subjectQueue, entry).Return(nil)
	queuePos.EXPECT().GetPosition(ctx, subjectQueue, entry).Return(expectedPos, nil)

	pos, err := service.Push(ctx, subjectQueue, entry)
	require.Empty(t, err)
	assert.Equal(t, expectedPos, pos)

	// Попаем айдишник
	queue.EXPECT().Pop(ctx, subjectQueue).Return(entry, nil)

	popedEntry, err := service.Pop(ctx, subjectQueue)
	require.Empty(t, err)
	assert.Equal(t, entry, popedEntry)
}

func TestQueuePop_NotFound(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	// Попаем и получаем ErrNotFound
	queue.EXPECT().Pop(ctx, subjectQueue).Return(models.QueueEntry{}, repositories.ErrNotFound)

	popedEntry, err := service.Pop(ctx, subjectQueue)
	require.NotEmpty(t, err)
	assert.ErrorIs(t, err, services.ErrNotFound)
	assert.Empty(t, popedEntry)
}

// QueueService.Clear

func TestQueueClear_Success(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	chatID := gofakeit.ID()
	entry := models.QueueEntry{
		ChatID: chatID,
	}

	// Пушаем первого юзера
	expectedPos := int64(1)
	queue.EXPECT().Push(ctx, subjectQueue, entry).Return(nil)
	queuePos.EXPECT().GetPosition(ctx, subjectQueue, entry).Return(expectedPos, nil)

	pos, err := service.Push(ctx, subjectQueue, entry)
	require.Empty(t, err)
	assert.Equal(t, expectedPos, pos)

	// Пушаем второго
	chatID = gofakeit.ID()
	entry = models.QueueEntry{
		ChatID: chatID,
	}

	expectedPos = int64(2)
	queue.EXPECT().Push(ctx, subjectQueue, entry).Return(nil)
	queuePos.EXPECT().GetPosition(ctx, subjectQueue, entry).Return(expectedPos, nil)

	pos, err = service.Push(ctx, subjectQueue, entry)
	require.Empty(t, err)
	assert.Equal(t, expectedPos, pos)

	// Очищаем очередь
	queue.EXPECT().Clear(ctx, subjectQueue).Return(nil)

	err = service.Clear(ctx, subjectQueue)
	require.Empty(t, err)
}

func TestQueueClear_NotFound(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	queue.EXPECT().Clear(ctx, subjectQueue).Return(repositories.ErrNotFound)

	err := service.Clear(ctx, subjectQueue)
	require.ErrorIs(t, err, services.ErrNotFound)
}
