package forwarder


import (
	"log"
	// "encoding/json"
	// "net/http"
	// "time"

	"github.com/freshpaint/hipaa-tracker/internal/models"
)


func ForwardEvent(e models.Event) error {

	log.Printf("Forwarding event ID %d (type: %s, sanitized: %v)", e.ID, e.EventType, e.Sanitized)

	return nil
}