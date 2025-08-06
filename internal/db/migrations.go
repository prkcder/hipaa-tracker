package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
)

func RunMigrations(db *sql.DB) error {
	slog.Info("Running migrations from scripts")

	sqlBytes, err := os.ReadFile("scripts/seed_db.sql")
	if err != nil {
		return err
	}

	if _, err := db.Exec(string(sqlBytes)); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	slog.Info("Migrations ran successfully")
	return err
}
