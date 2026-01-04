package queue_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/iskanye/mirea-queue/internal/interfaces"
	"github.com/iskanye/mirea-queue/internal/models"
	"github.com/iskanye/mirea-queue/internal/repositories"
	"github.com/iskanye/mirea-queue/internal/services"
	"github.com/iskanye/mirea-queue/internal/services/queue"
	"github.com/iskanye/mirea-queue/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	queueRange = int64(10)
	timeout    = time.Second
)

var (
	queueBase    *mocks.MockQueue
	queueViewer  *mocks.MockQueueViewer
	queuePos     *mocks.MockQueuePosition
	queueLength  *mocks.MockQueueLength
	queueSwap    *mocks.MockQueueSwap
	queueRemover *mocks.MockQueueRemover
)

// Создаем сервис очереди с приклеплёнными к нему моками
func newService(t *testing.T) (interfaces.QueueService, context.Context) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	queueBase = mocks.NewMockQueue(t)
	queueViewer = mocks.NewMockQueueViewer(t)
	queuePos = mocks.NewMockQueuePosition(t)
	queueLength = mocks.NewMockQueueLength(t)
	queueSwap = mocks.NewMockQueueSwap(t)
	queueRemover = mocks.NewMockQueueRemover(t)
	cache = mocks.NewMockCache(t)

	t.Cleanup(func() {
		// Очищаем ожидаемые вызовы моков
		queueBase.ExpectedCalls = nil
		queueViewer.ExpectedCalls = nil
		queuePos.ExpectedCalls = nil
		queueLength.ExpectedCalls = nil
		queueSwap.ExpectedCalls = nil
		queueRemover.ExpectedCalls = nil
		cache.ExpectedCalls = nil
		cancel()
	})

	return queue.New(
		slog.New(slog.DiscardHandler),
		queueRange,
		queueBase,
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
	queueBase.EXPECT().Push(ctx, subjectQueue, entry).Return(nil)
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

	queueBase.EXPECT().Push(ctx, subjectQueue, entry).Return(repositories.ErrAlreadyInQueue)

	pos, err := service.Push(ctx, subjectQueue, entry)
	require.ErrorIs(t, err, services.ErrAlreadyInQueue)
	assert.Empty(t, pos)
}

func TestQueuePush_Failure(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	chatID := gofakeit.ID()
	entry := models.QueueEntry{
		ChatID: chatID,
	}

	expectedErr := errors.New("внезапная ошибка на стороне базы данных")
	queueBase.EXPECT().Push(ctx, subjectQueue, entry).Return(expectedErr)

	pos, err := service.Push(ctx, subjectQueue, entry)
	require.ErrorIs(t, err, expectedErr)
	assert.Empty(t, pos)
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

	// Попаем айдишник
	queueBase.EXPECT().Pop(ctx, subjectQueue).Return(entry, nil)

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
	queueBase.EXPECT().Pop(ctx, subjectQueue).Return(models.QueueEntry{}, repositories.ErrNotFound)

	popedEntry, err := service.Pop(ctx, subjectQueue)
	require.ErrorIs(t, err, services.ErrNotFound)
	assert.Empty(t, popedEntry)
}

func TestQueuePop_Failure(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	expectedErr := errors.New("внезапная ошибка на стороне базы данных")
	queueBase.EXPECT().Pop(ctx, subjectQueue).Return(models.QueueEntry{}, expectedErr)

	popedEntry, err := service.Pop(ctx, subjectQueue)
	require.ErrorIs(t, err, expectedErr)
	assert.Empty(t, popedEntry)
}

// QueueService.Clear

func TestQueueClear_Success(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	// Очищаем очередь
	queueBase.EXPECT().Clear(ctx, subjectQueue).Return(nil)

	err := service.Clear(ctx, subjectQueue)
	require.Empty(t, err)
}

func TestQueueClear_NotFound(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	queueBase.EXPECT().Clear(ctx, subjectQueue).Return(repositories.ErrNotFound)

	err := service.Clear(ctx, subjectQueue)
	require.ErrorIs(t, err, services.ErrNotFound)
}

func TestQueueClear_Failure(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	expectedErr := errors.New("внезапная ошибка на стороне базы данных")
	queueBase.EXPECT().Clear(ctx, subjectQueue).Return(expectedErr)

	err := service.Clear(ctx, subjectQueue)
	require.ErrorIs(t, err, expectedErr)
}

// QueueService.Range

func TestQueueRange_Success(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	entries := []models.QueueEntry{
		{ChatID: gofakeit.ID()},
		{ChatID: gofakeit.ID()},
		{ChatID: gofakeit.ID()},
	}

	queueViewer.EXPECT().Range(ctx, subjectQueue, queueRange).Return(entries, nil)

	result, err := service.Range(ctx, subjectQueue)
	require.Empty(t, err)
	assert.Equal(t, entries, result)
}

func TestQueueRange_NotFound(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	queueViewer.EXPECT().Range(ctx, subjectQueue, queueRange).Return(nil, repositories.ErrNotFound)

	result, err := service.Range(ctx, subjectQueue)
	require.ErrorIs(t, err, services.ErrNotFound)
	assert.Empty(t, result)
}

func TestQueueRange_Failure(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	expectedErr := errors.New("внезапная ошибка на стороне базы данных")
	queueViewer.EXPECT().Range(ctx, subjectQueue, queueRange).Return(nil, expectedErr)

	result, err := service.Range(ctx, subjectQueue)
	require.ErrorIs(t, err, expectedErr)
	assert.Empty(t, result)
}

// QueueService.Remove

func TestQueueRemove_Success(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	entry := models.QueueEntry{
		ChatID: gofakeit.ID(),
	}

	queueRemover.EXPECT().Remove(ctx, subjectQueue, entry).Return(nil)

	err := service.Remove(ctx, subjectQueue, entry)
	require.Empty(t, err)
}

func TestQueueRemove_Failure(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	entry := models.QueueEntry{
		ChatID: gofakeit.ID(),
	}

	expectedErr := errors.New("внезапная ошибка на стороне базы данных")
	queueRemover.EXPECT().Remove(ctx, subjectQueue, entry).Return(expectedErr)

	err := service.Remove(ctx, subjectQueue, entry)
	require.ErrorIs(t, err, expectedErr)
}

// QueueService.LetAhead

func TestQueueLetAhead_Success(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	entry := models.QueueEntry{
		ChatID: gofakeit.ID(),
	}

	queuePos.EXPECT().GetPosition(ctx, subjectQueue, entry).Return(2, nil)
	queueLength.EXPECT().Len(ctx, subjectQueue).Return(5, nil)
	queueSwap.EXPECT().LetAhead(ctx, subjectQueue, entry).Return(nil)

	err := service.LetAhead(ctx, subjectQueue, entry)
	require.Empty(t, err)
}

func TestQueueLetAhead_NotFound(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	entry := models.QueueEntry{
		ChatID: gofakeit.ID(),
	}

	queuePos.EXPECT().GetPosition(ctx, subjectQueue, entry).Return(0, repositories.ErrNotFound)

	err := service.LetAhead(ctx, subjectQueue, entry)
	require.ErrorIs(t, err, services.ErrNotFound)
}

func TestQueueLetAhead_QueueEnd(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	entry := models.QueueEntry{
		ChatID: gofakeit.ID(),
	}

	queuePos.EXPECT().GetPosition(ctx, subjectQueue, entry).Return(5, nil)
	queueLength.EXPECT().Len(ctx, subjectQueue).Return(5, nil)

	err := service.LetAhead(ctx, subjectQueue, entry)
	require.NotEmpty(t, err)
	assert.ErrorIs(t, err, services.ErrQueueEnd)
}

func TestQueueLetAhead_Failure1(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	entry := models.QueueEntry{
		ChatID: gofakeit.ID(),
	}

	expectedErr := errors.New("внезапная ошибка на стороне базы данных")
	queuePos.EXPECT().GetPosition(ctx, subjectQueue, entry).Return(0, expectedErr)

	err := service.LetAhead(ctx, subjectQueue, entry)
	require.ErrorIs(t, err, expectedErr)
}

func TestQueueLetAhead_Failure2(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	entry := models.QueueEntry{
		ChatID: gofakeit.ID(),
	}

	expectedErr := errors.New("внезапная ошибка на стороне базы данных")
	queuePos.EXPECT().GetPosition(ctx, subjectQueue, entry).Return(1, nil)
	queueLength.EXPECT().Len(ctx, subjectQueue).Return(0, expectedErr)

	err := service.LetAhead(ctx, subjectQueue, entry)
	require.ErrorIs(t, err, expectedErr)
}

func TestQueueLetAhead_Failure3(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	entry := models.QueueEntry{
		ChatID: gofakeit.ID(),
	}

	expectedErr := errors.New("внезапная ошибка на стороне базы данных")
	queuePos.EXPECT().GetPosition(ctx, subjectQueue, entry).Return(2, nil)
	queueLength.EXPECT().Len(ctx, subjectQueue).Return(5, nil)
	queueSwap.EXPECT().LetAhead(ctx, subjectQueue, entry).Return(expectedErr)

	err := service.LetAhead(ctx, subjectQueue, entry)
	require.ErrorIs(t, err, expectedErr)
}

// QueueService.GetPosition

func TestQueueGetPosition_Success(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	entry := models.QueueEntry{
		ChatID: gofakeit.ID(),
	}

	expectedPos := int64(3)
	queuePos.EXPECT().GetPosition(ctx, subjectQueue, entry).Return(expectedPos, nil)

	pos, err := service.Pos(ctx, subjectQueue, entry)
	require.Empty(t, err)
	assert.Equal(t, expectedPos, pos)
}

func TestQueueGetPosition_NotFound(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	entry := models.QueueEntry{
		ChatID: gofakeit.ID(),
	}

	queuePos.EXPECT().GetPosition(ctx, subjectQueue, entry).Return(0, repositories.ErrNotFound)

	pos, err := service.Pos(ctx, subjectQueue, entry)
	require.ErrorIs(t, err, services.ErrNotFound)
	assert.Empty(t, pos)
}

func TestQueueGetPosition_Failure(t *testing.T) {
	service, ctx := newService(t)

	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	entry := models.QueueEntry{
		ChatID: gofakeit.ID(),
	}

	expectedErr := errors.New("внезапная ошибка на стороне базы данных")
	queuePos.EXPECT().GetPosition(ctx, subjectQueue, entry).Return(0, expectedErr)

	pos, err := service.Pos(ctx, subjectQueue, entry)
	require.ErrorIs(t, err, expectedErr)
	assert.Empty(t, pos)
}
