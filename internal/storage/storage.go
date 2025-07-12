package storage

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/freshpaint/hipaa-tracker/internal/models"
)

// SaveEvent stores a sanitized event into the database
func SaveEvent(db *sql.DB, e *models.Event) error {

	e.CreatedAt = time.Now()

	payloadBytes, err := json.Marshal(e.Payload)
	
	if err != nil {
		return err
	}


	query := `
		INSERT INTO events (event_type, payload, sanitized, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return db.QueryRow(query, e.EventType, payloadBytes, e.Sanitized, e.CreatedAt).
		Scan(&e.ID, &e.CreatedAt)
}