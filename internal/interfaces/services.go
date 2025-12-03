package interfaces

import (
	"context"

	"github.com/iskanye/mirea-queue/internal/models"
)

// Сервис очереди
type QueueService interface {
	// Пушает пользователя в очередь
	Push(
		ctx context.Context,
		queue models.Queue,
		entry models.QueueEntry,
	) (int64, error)
	// Получает пользователя из начала очереди
	// и удаляет его из очереди
	Pop(
		ctx context.Context,
		queue models.Queue,
	) (models.QueueEntry, error)
	// Очищает очередь
	Clear(
		ctx context.Context,
		queue models.Queue,
		key string,
	) error
	// Получает длину очереди
	Len(
		ctx context.Context,
		queue models.Queue,
	) (int64, error)
	// Получает текущую позицию пользователя в очереди
	Pos(
		ctx context.Context,
		queue models.Queue,
		entry models.QueueEntry,
	) (int64, error)
}

// Сервис пользователей
type UsersService interface {
	// Создает нового пользователя и возвращает его
	CreateUser(
		ctx context.Context,
		chatID int64,
		user models.User,
	) (models.User, error)
	// Удаляет пользователя
	RemoveUser(
		ctx context.Context,
		chatID int64,
	) error
	// Обновляет данные пользователя
	UpdateUser(
		ctx context.Context,
		chatID int64,
		user models.User,
	) (models.User, error)
	// Получает пользователя
	// Если его нет то возвращает ErrNotFound
	GetUser(
		ctx context.Context,
		chatID int64,
	) (models.User, error)
}

// Сервис проверки прав админа
type AdminService interface {
	ValidateToken(string) bool
}
