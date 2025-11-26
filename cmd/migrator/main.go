package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/iskanye/mirea-queue/internal/config"
)

func main() {
	var migrationsPath, migrationsTable string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()

	cfg := config.MustLoadConfig()

	uri := fmt.Sprintf("%s:%s@%s:%d/%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DB,
	)

	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("pgx5://%s?x-migrations-table=%s&sslmode=disable", uri, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	mustMigrate(m.Up())
}

func mustMigrate(err error) {
	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("no migrations to apply")
		return
	}
	if err != nil {
		panic(err)
	}

	fmt.Println("migrations applied")
}
