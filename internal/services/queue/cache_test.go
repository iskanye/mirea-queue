package queue

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/iskanye/mirea-queue/internal/models"
	"github.com/iskanye/mirea-queue/internal/repositories"
	"github.com/iskanye/mirea-queue/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// QueueService.SaveToCache

func TestSaveToCache_Success(t *testing.T) {
	service, ctx := newService(t)

	chatID := gofakeit.Int64()
	subjectQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	cache.EXPECT().Set(ctx, fmt.Sprint(chatID), subjectQueue.Key()).Return(nil)

	err := service.SaveToCache(ctx, chatID, subjectQueue)
	require.Empty(t, err)
}

func TestSaveToCache_MultipleCalls_Success(t *testing.T) {
	service, ctx := newService(t)

	chatID1 := gofakeit.Int64()
	chatID2 := gofakeit.Int64()

	queue1 := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}
	queue2 := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	cache.EXPECT().Set(ctx, fmt.Sprint(chatID1), queue1.Key()).Return(nil)
	cache.EXPECT().Set(ctx, fmt.Sprint(chatID2), queue2.Key()).Return(nil)

	err := service.SaveToCache(ctx, chatID1, queue1)
	require.Empty(t, err)

	err = service.SaveToCache(ctx, chatID2, queue2)
	require.Empty(t, err)
}

// QueueService.GetFromCache

func TestGetFromCache_Success(t *testing.T) {
	service, ctx := newService(t)

	chatID := gofakeit.Int64()
	expectedQueue := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	cache.EXPECT().Get(ctx, fmt.Sprint(chatID)).Return(expectedQueue.Key(), nil)

	queue, err := service.GetFromCache(ctx, chatID)
	require.Empty(t, err)
	assert.Equal(t, expectedQueue, queue)
}

func TestGetFromCache_CacheMiss(t *testing.T) {
	service, ctx := newService(t)

	chatID := gofakeit.Int64()

	cache.EXPECT().Get(ctx, fmt.Sprint(chatID)).Return("", repositories.ErrCacheMiss)

	queue, err := service.GetFromCache(ctx, chatID)
	require.ErrorIs(t, err, services.ErrNotFound)
	assert.Empty(t, queue)
}

func TestGetFromCache_MultipleCalls_Success(t *testing.T) {
	service, ctx := newService(t)

	chatID1 := gofakeit.Int64()
	chatID2 := gofakeit.Int64()

	queue1 := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}
	queue2 := models.Queue{
		Group:   gofakeit.ID(),
		Subject: gofakeit.Noun(),
	}

	cache.EXPECT().Get(ctx, fmt.Sprint(chatID1)).Return(queue1.Key(), nil)
	cache.EXPECT().Get(ctx, fmt.Sprint(chatID2)).Return(queue2.Key(), nil)

	result1, err1 := service.GetFromCache(ctx, chatID1)
	require.Empty(t, err1)
	assert.Equal(t, queue1, result1)

	result2, err2 := service.GetFromCache(ctx, chatID2)
	require.Empty(t, err2)
	assert.Equal(t, queue2, result2)
}
