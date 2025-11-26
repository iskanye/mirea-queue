package postgres

import (
	"context"
	"fmt"

	"github.com/iskanye/mirea-queue/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(cfg config.Config) (*Storage, error) {
	const op = "postgres.New"

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=%d",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Password,
		cfg.Postgres.PoolConns,
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Подключаемся к Postgres
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		pool: pool,
	}, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}
