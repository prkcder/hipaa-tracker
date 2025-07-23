package forwarder


import (
	"log"
	"encoding/json"
	"net/http"
	"time"
	"bytes"
	"github.com/freshpaint/hipaa-tracker/internal/models"
)


var ForwardEventFunc = ForwardEvent // for mocking in tests

// ForwardEvent simulates sending an event to an external system
func ForwardEvent(e models.Event) error {
	// Simulate with logging for now
	log.Printf("Forwarding event ID %d (type: %s, sanitized: %v)", e.ID, e.EventType, e.Sanitized)

	// --- Optional: simulate HTTP POST (for Phase 3) ---
	forwardURL := "http://172.17.0.1:9000/mock-endpoint"

	// forwardURL := "https://your-api-id.execute-api.us-west-2.amazonaws.com/prod/receive"
	// req.Header.Set("x-api-key", os.Getenv("AWS_API_KEY"))

	jsonBytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", forwardURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Printf("Forwarding failed with status %d", resp.StatusCode)
	} else {
		log.Printf("✅ Event forwarded successfully to %s", forwardURL)
	}

	return nil
}


// import (
// 	"bytes"
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/freshpaint/hipaa-tracker/internal/models"
// )

// var ForwardEventFunc = ForwardEvent

// func ForwardEvent(e models.Event) error {
// 	// Simulate with dummy endpoint
// 	forwardURL := "https://httpbin.org/post" // replace with real analytics system if needed

// 	jsonBytes, err := json.Marshal(e)
// 	if err != nil {
// 		return err
// 	}

// 	req, err := http.NewRequest("POST", forwardURL, bytes.NewBuffer(jsonBytes))
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{Timeout: 3 * time.Second}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode >= 300 {
// 		log.Printf("Forwarding failed with status %d", resp.StatusCode)
// 	} else {
// 		log.Printf("Successfully forwarded event ID %d", e.ID)
// 	}

// 	return nil
// }
