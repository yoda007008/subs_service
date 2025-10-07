package migrator

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log/slog"
	"path/filepath"
)

func RunMigrations(connString string, migrationsPath string) error {
	abcPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("Please get absolute PATH: %w", err)
	}

	slog.Info("Apply migrations on these path:", "path", abcPath)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return fmt.Errorf("Error opening DB: %w", err)
	}

	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("Error driver Postgres: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+abcPath,
		"postgres", driver,
	)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up failed: %w", err)
	}

	slog.Info("Migrations was succesfully")
	return nil
}
