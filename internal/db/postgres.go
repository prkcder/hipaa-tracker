package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {

	// dsn := fmt.Sprintf(
	// 	"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	//     os.Getenv("DB_HOST"),
	//     os.Getenv("DB_PORT"),
	//     os.Getenv("DB_USER"),
	//     os.Getenv("DB_PASSWORD"),
	//     os.Getenv("DB_NAME"),
	// )

	// Using Render db
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
