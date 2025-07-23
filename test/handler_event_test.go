package test


import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"errors"
	"time"

	"github.com/freshpaint/hipaa-tracker/internal/handlers"
	"github.com/freshpaint/hipaa-tracker/internal/models"
	"github.com/freshpaint/hipaa-tracker/internal/sanitize"
	"github.com/freshpaint/hipaa-tracker/internal/storage"
	"github.com/freshpaint/hipaa-tracker/internal/forwarder"
)

// override SaveEvent and ForwardEvent before each test
func setupTestMocks() {
	storage.SaveEventFunc = func(_ *sql.DB, e *models.Event) error {
		e.ID = 999
		e.CreatedAt = time.Now()
		return nil
	}

	forwarder.ForwardEventFunc = func(e models.Event) error {
		// Simulate success
		return nil
	}
}

func TestNewEventHandler_ValidPayload(t *testing.T) {
	setupTestMocks()
	sanitize.InitTestConfig([]string{"email", "password"})

	payload := map[string]interface{}{
		"event_type": "login",
		"payload": map[string]interface{}{
			"email":    "secret@example.com",
			"password": "123456",
			"username": "tester",
		},
	}
	jsonBody, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := handlers.NewEventHandler(nil)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status 201 Created, got %d", status)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["sanitized"] != true {
		t.Errorf("Expected sanitized = true, got %v", response["sanitized"])
	}
}

func TestNewEventHandler_InvalidMethod(t *testing.T) {
	setupTestMocks()

	req := httptest.NewRequest(http.MethodGet, "/event", nil)
	rr := httptest.NewRecorder()

	handler := handlers.NewEventHandler(nil)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", rr.Code)
	}
}

func TestNewEventHandler_InvalidJSON(t *testing.T) {
	setupTestMocks()

	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer([]byte(`not-json`)))
	rr := httptest.NewRecorder()

	handler := handlers.NewEventHandler(nil)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", rr.Code)
	}
}

func TestNewEventHandler_MissingFields(t *testing.T) {
	setupTestMocks()

	body := map[string]interface{}{
		"payload": map[string]interface{}{"x": 1},
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := handlers.NewEventHandler(nil)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", rr.Code)
	}
}


func TestNewEventHandler_DBInsertFails(t *testing.T) {
	// Override storage to simulate failure
	storage.SaveEventFunc = func(_ *sql.DB, e *models.Event) error {
		return errors.New("mock DB failure")
	}

	// Sanitize config (won’t matter here)
	sanitize.InitTestConfig([]string{"email"})

	payload := map[string]interface{}{
		"event_type": "test_error",
		"payload": map[string]interface{}{
			"email": "x@example.com",
		},
	}
	jsonBody, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := handlers.NewEventHandler(nil)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 Internal Server Error, got %d", rr.Code)
	}
}


func TestNewEventHandler_MissingPayloadField(t *testing.T) {
	setupTestMocks()

	payload := map[string]interface{}{
		"event_type": "incomplete_event",
		// payload is missing
	}
	jsonBody, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := handlers.NewEventHandler(nil)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for missing payload field, got %d", rr.Code)
	}
}

func TestNewEventHandler_ForwarderCalled(t *testing.T) {
	called := false

	storage.SaveEventFunc = func(_ *sql.DB, e *models.Event) error {
		e.ID = 1001
		e.CreatedAt = time.Now()
		return nil
	}

	forwarder.ForwardEventFunc = func(e models.Event) error {
		called = true
		return nil
	}

	sanitize.InitTestConfig([]string{"token"})

	payload := map[string]interface{}{
		"event_type": "test_forward",
		"payload": map[string]interface{}{
			"token": "abc123",
			"name":  "forward_test",
		},
	}
	jsonBody, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := handlers.NewEventHandler(nil)
	handler.ServeHTTP(rr, req)

	if !called {
		t.Errorf("Expected forwarder to be called, but it wasn't")
	}
}

