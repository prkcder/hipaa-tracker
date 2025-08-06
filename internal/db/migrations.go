package db

import (
    "database/sql"
    "os"

    _ "github.com/lib/pq"
)

func RunMigrations(db *sql.DB) error {
    sqlBytes, err := os.ReadFile("scripts/seed_db.sql") // update path if needed
    if err != nil {
        return err
    }

    _, err = db.Exec(string(sqlBytes))
    return err
}