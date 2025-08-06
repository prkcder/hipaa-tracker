package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/freshpaint/hipaa-tracker/internal/db"
	"github.com/freshpaint/hipaa-tracker/internal/handlers"
	"github.com/freshpaint/hipaa-tracker/internal/middleware"
	"github.com/freshpaint/hipaa-tracker/internal/sanitize"
	// "github.com/freshpaint/hipaa-tracker/internal/forwarder"
)

func main() {

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	log.Println("starting tracker")

	if err := sanitize.LoadSensitiveFields("sensitive_fields.yaml"); err != nil {
		log.Fatalf("Failed to load sensitive fields config: %v", err)
	}

	database, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer database.Close()

	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("Failed to run DB migrations: %v", err)
	}

	mux := http.NewServeMux()

	handlers.RegisterRoutes(mux, database)

	handler := middleware.WithCORS(mux)

	log.Println("🚀 Server is listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
