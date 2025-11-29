package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/iskanye/mirea-queue/internal/models"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) CreateUser(
	ctx context.Context,
	chatID int64,
	user models.User,
) error {
	const op = "postgres.CreateUser"

	// Вставляем группу
	getGroupID := s.pool.QueryRow(
		ctx,
		`
		WITH inserted AS (
			INSERT INTO student_groups (group_name)
			VALUES ($1)
			ON CONFLICT DO NOTHING
			RETURNING group_id
		)
		SELECT group_id FROM inserted
		UNION ALL
		SELECT group_id FROM student_groups 
		WHERE group_name = $1
		LIMIT 1;
		`,
		user.Group,
	)

	var groupID int
	err := getGroupID.Scan(&groupID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Вставляем пользователя
	_, err = s.pool.Exec(
		ctx,
		`
		INSERT INTO users (chat_id, name, group_id, queue_access)
		VALUES ($1, $2, $3, $4);
		`,
		chatID, user.Name, groupID, user.QueueAccess,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) RemoveUser(
	ctx context.Context,
	chatID int64,
) error {
	const op = "postgres.RemoveUser"

	// Удаляем пользователя
	_, err := s.pool.Exec(
		ctx,
		`
		DELETE FROM users
		WHERE chat_id = $1;
		`,
		chatID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateUser(
	ctx context.Context,
	chatID int64,
	user models.User,
) error {
	const op = "postgres.UpdateUser"

	// Вставляем группу
	getGroupID := s.pool.QueryRow(
		ctx,
		`
		WITH inserted AS (
			INSERT INTO student_groups (group_name)
			VALUES ($1)
			ON CONFLICT DO NOTHING
			RETURNING group_id;
		)
		SELECT group_id FROM inserted
		UNION ALL
		SELECT group_id FROM student_groups 
		WHERE group_name = $1
		LIMIT 1;
		`,
		user.Group,
	)

	var groupID int
	err := getGroupID.Scan(&groupID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Обновляем пользователя
	_, err = s.pool.Exec(
		ctx,
		`
		UPDATE users
		SET name = $1, group_id = $2, queue_access = $3
		WHERE chat_id = $4;
		`,
		user.Name, groupID, user.QueueAccess, chatID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetUser(
	ctx context.Context,
	chatID int64,
) (models.User, error) {
	const op = "postgres.GetUser"

	// Получаем пользователя
	getUser := s.pool.QueryRow(
		ctx,
		`
		SELECT name, group_name, queue_access
		FROM users
		JOIN student_groups USING(group_id)
		WHERE chatID = $1;
		`,
		chatID,
	)

	var user models.User
	err := getUser.Scan(&user.Name, &user.Group, &user.QueueAccess)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
