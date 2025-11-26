package postgres

import (
	"context"
	"fmt"

	"github.com/iskanye/mirea-queue/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(cfg config.Config) (*Storage, error) {
	const op = "repositories.postgres.New"

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Password,
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return &Storage{
		pool: pool,
	}, nil
}
