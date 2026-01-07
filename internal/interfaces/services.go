package interfaces

import (
	"context"

	"github.com/iskanye/mirea-queue/internal/models"
)

// Сервис очереди
//
//mockery:generate: false
type QueueService interface {
	// Пушает пользователя в очередь
	Push(
		ctx context.Context,
		queue models.Queue,
		entry models.QueueEntry,
	) error
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
	) error
	// Получает текущую позицию пользователя в очереди
	Pos(
		ctx context.Context,
		queue models.Queue,
		entry models.QueueEntry,
	) (int64, error)
	// Пропускает пользователя в очереди
	LetAhead(
		ctx context.Context,
		queue models.Queue,
		entry models.QueueEntry,
	) error
	// Получает содержимое очереди (первые несколько элементов)
	Range(
		ctx context.Context,
		queue models.Queue,
	) ([]models.QueueEntry, error)
	// Сохраняет очередь в кеш
	SaveToCache(
		ctx context.Context,
		chatID int64,
		queue models.Queue,
	) error
	// Получает очеред из кеша
	GetFromCache(
		ctx context.Context,
		chatID int64,
	) (models.Queue, error)
	// Удаляет пользователя из очереди
	Remove(
		ctx context.Context,
		queue models.Queue,
		entry models.QueueEntry,
	) error
}

// Сервис пользователей
//
//mockery:generate: false
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
//
//mockery:generate: false
type AdminService interface {
	ValidateToken(string) bool
}

// Сервис расписания и групп
//
//mockery:generate: false
type ScheduleService interface {
	// Получает список доступных групп в расписании по названию группы
	GetGroups(
		ctx context.Context,
		group string,
	) ([]models.Group, error)
	// Получает предметы конкретной группы
	GetSubjects(
		ctx context.Context,
		group models.Group,
	) ([]string, error)
}
