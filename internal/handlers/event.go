package handlers


import (
	"database/sql"
	"encoding/json"
	"net/http"
	"log"

	"github.com/freshpaint/hipaa-tracker/internal/db"
	
	"github.com/freshpaint/hipaa-tracker/internal/models"
	"github.com/freshpaint/hipaa-tracker/internal/sanitize"
	"github.com/freshpaint/hipaa-tracker/internal/forwarder"
	"github.com/freshpaint/hipaa-tracker/internal/storage"
)

func NewEventHandler(database *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var input models.Event

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if input.EventType == "" || input.Payload == nil {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		// sanitize
		input.Payload, input.Sanitized = sanitize.Sanitize(input.Payload)

		// Insert using db package
		if err := storage.SaveEvent(database, &input); err != nil {
			log.Printf("❌ DB error: %v", err) 
			http.Error(w, "Failed to store event", http.StatusInternalServerError)
			return
		}

		//  Forward to downstream system
		if err := forwarder.ForwardEvent(input); err != nil {
			// add log here 
			// Forwarding failures shouldn't stop successful DB inserts
			log.Printf("Failed to forward event ID %d: %v", input.ID, err)
		}

		// Return the inserted record's ID and timestamp
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(map[string]any{
			"id":         input.ID,
			"created_at": input.CreatedAt,
			"sanitized": input.Sanitized,
		})
	}
}


func GetEventsHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		events, err := db.GetAllEvents(database)
		
		if err != nil {
			http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(events)
	}
}