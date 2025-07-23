package test



import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/freshpaint/hipaa-tracker/internal/db"
	"github.com/freshpaint/hipaa-tracker/internal/handlers"
	"github.com/freshpaint/hipaa-tracker/internal/models"
)

// mock GetAllEvents override
func mockGetAllEvents(_ *sql.DB) ([]models.Event, error) {
	return []models.Event{
		{
			ID:        1,
			EventType: "mock",
			Payload: map[string]any{
				"mock": true,
			},
			Sanitized: true,
			CreatedAt: time.Now(),
		},
	}, nil
}

func TestGetEventsHandler_ReturnsEvents(t *testing.T) {
	// swap out real function with mock
	original := db.GetAllEvents
	db.GetAllEvents = mockGetAllEvents
	defer func() { db.GetAllEvents = original }()

	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	rr := httptest.NewRecorder()

	handler := handlers.GetEventsHandler(nil)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", rr.Code)
	}

	var events []models.Event
	if err := json.NewDecoder(rr.Body).Decode(&events); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(events) != 1 || events[0].EventType != "mock" {
		t.Errorf("Unexpected response: %+v", events)
	}
}