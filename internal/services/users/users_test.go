package users_test

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
	"github.com/iskanye/mirea-queue/internal/services/users"
	"github.com/iskanye/mirea-queue/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const timeout = time.Second

var (
	userCreator  *mocks.MockUserCreator
	userRemover  *mocks.MockUserRemover
	userModifier *mocks.MockUserModifier
	userProvider *mocks.MockUserProvider
)

// Создаем сервис пользователей с приклеплёнными к нему моками
func newService(t *testing.T) (interfaces.UsersService, context.Context) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	userCreator = mocks.NewMockUserCreator(t)
	userRemover = mocks.NewMockUserRemover(t)
	userModifier = mocks.NewMockUserModifier(t)
	userProvider = mocks.NewMockUserProvider(t)

	t.Cleanup(func() {
		// Очищаем ожидаемые вызовы моков
		userCreator.ExpectedCalls = nil
		userRemover.ExpectedCalls = nil
		userModifier.ExpectedCalls = nil
		userProvider.ExpectedCalls = nil
		cancel()
	})

	return users.New(
		slog.New(slog.DiscardHandler),
		userCreator,
		userRemover,
		userModifier,
		userProvider,
	), ctx
}

func TestUsersCreate_Success(t *testing.T) {
	service, ctx := newService(t)

	chatID := gofakeit.Int64()
	user := models.User{
		Name: gofakeit.Username(),
	}

	userCreator.EXPECT().CreateUser(ctx, chatID, user).Return(nil)

	res, err := service.CreateUser(ctx, chatID, user)
	require.NoError(t, err)
	assert.Equal(t, user, res)
}

func TestUsersCreate_Failure(t *testing.T) {
	service, ctx := newService(t)

	chatID := gofakeit.Int64()
	user := models.User{Name: gofakeit.Username()}

	expectedErr := errors.New("внезапная ошибка на стороне базы данных")
	userCreator.EXPECT().CreateUser(ctx, chatID, user).Return(expectedErr)

	res, err := service.CreateUser(ctx, chatID, user)
	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Empty(t, res)
}

func TestUsersRemove_Success(t *testing.T) {
	service, ctx := newService(t)

	chatID := gofakeit.Int64()

	userRemover.EXPECT().RemoveUser(ctx, chatID).Return(nil)

	err := service.RemoveUser(ctx, chatID)
	require.NoError(t, err)
}

func TestUsersRemove_Failure(t *testing.T) {
	service, ctx := newService(t)

	chatID := gofakeit.Int64()

	expectedErr := errors.New("внезапная ошибка на стороне базы данных")
	userRemover.EXPECT().RemoveUser(ctx, chatID).Return(expectedErr)

	err := service.RemoveUser(ctx, chatID)
	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
}

func TestUsersUpdate_Success(t *testing.T) {
	service, ctx := newService(t)

	chatID := gofakeit.Int64()
	user := models.User{Name: gofakeit.Username()}

	userModifier.EXPECT().UpdateUser(ctx, chatID, user).Return(nil)

	res, err := service.UpdateUser(ctx, chatID, user)
	require.NoError(t, err)
	assert.Equal(t, user, res)
}

func TestUsersUpdate_Failure(t *testing.T) {
	service, ctx := newService(t)

	chatID := gofakeit.Int64()
	user := models.User{Name: gofakeit.Username()}

	expectedErr := errors.New("внезапная ошибка на стороне базы данных")
	userModifier.EXPECT().UpdateUser(ctx, chatID, user).Return(expectedErr)

	res, err := service.UpdateUser(ctx, chatID, user)
	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Empty(t, res)
}

func TestUsersGet_Success(t *testing.T) {
	service, ctx := newService(t)

	chatID := gofakeit.Int64()
	user := models.User{Name: gofakeit.Username()}

	userProvider.EXPECT().GetUser(ctx, chatID).Return(user, nil)

	res, err := service.GetUser(ctx, chatID)
	require.NoError(t, err)
	assert.Equal(t, user, res)
}

func TestUsersGet_NotFound(t *testing.T) {
	service, ctx := newService(t)

	chatID := gofakeit.Int64()

	userProvider.EXPECT().GetUser(ctx, chatID).Return(models.User{}, repositories.ErrNotFound)

	res, err := service.GetUser(ctx, chatID)
	require.Error(t, err)
	assert.ErrorIs(t, err, services.ErrNotFound)
	assert.Empty(t, res)
}
