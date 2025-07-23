package handlers


import (
	"database/sql"
	"encoding/json"
	"net/http"
	"log/slog"

	"github.com/freshpaint/hipaa-tracker/internal/db"
	
	"github.com/freshpaint/hipaa-tracker/internal/models"
	"github.com/freshpaint/hipaa-tracker/internal/sanitize"
	"github.com/freshpaint/hipaa-tracker/internal/forwarder"
	"github.com/freshpaint/hipaa-tracker/internal/storage"
)

func NewEventHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Incoming /event request", "method", r.Method, "remote", r.RemoteAddr)

		if r.Method != http.MethodPost {

			slog.Warn("Invalid method for /event", "method", r.Method)

			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var input models.Event

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {

			slog.Warn("Missing required fields in payload", "event_type", input.EventType)

			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if input.EventType == "" || input.Payload == nil {

			slog.Warn("Missing required fields in payload", "event_type", input.EventType)

			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		slog.Info("Event received", "event_type", input.EventType, "payload_keys", keys(input.Payload))

		input.Payload, input.Sanitized = sanitize.Sanitize(input.Payload)

		slog.Info("Sanitized payload", "sanitized", input.Sanitized, "cleaned_keys", keys(input.Payload))

		// Insert using db package
		if err := storage.SaveEventFunc(database, &input); err != nil {
			slog.Error("Database insert failed", "error", err)
			http.Error(w, "Failed to store event", http.StatusInternalServerError)
			return
		}

		//  Forward to downstream system
		if err := forwarder.ForwardEventFunc(input); err != nil {
			// add log here 
			// Forwarding failures shouldn't stop successful DB inserts
			slog.Error("Forwarding failed", "event_id", input.ID, "error", err)
		} else {
			slog.Info("📤 Event forwarded", "event_id", input.ID)
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

// helper function for readable logs
func keys(m map[string]any) []string {
	k := make([]string, 0, len(m))
	for key := range m {
		k = append(k, key)
	}
	return k
}


func GetEventsHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("GET /events received", "method", r.Method, "remote_addr", r.RemoteAddr)

		events, err := db.GetAllEvents(database)
		
		if err != nil {
			slog.Error("Failed to fetch events from DB", "error", err)
			http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
			return
		}

		slog.Info("Events fetched", "count", len(events))
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(events); err != nil {
			slog.Error("Failed to encode events to JSON", "error", err)
		}
	}
}