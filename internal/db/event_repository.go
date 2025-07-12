package db

import (
	"database/sql"
	"encoding/json"


	"github.com/freshpaint/hipaa-tracker/internal/models" 
)


func GetAllEvents(db *sql.DB) ([]models.Event, error) {

	rows, err := db.Query(`
		SELECT id, event_type, payload, sanitized, created_at
		FROM events
		ORDER BY created_at DESC`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var e models.Event

		var payload []byte

		if err := rows.Scan(&e.ID, &e.EventType, &payload, &e.Sanitized, &e.CreatedAt); err != nil {
			continue
		}

		if err := json.Unmarshal(payload, &e.Payload); err != nil {
			continue
		}

		events = append(events, e)
	}

	return events, nil
}