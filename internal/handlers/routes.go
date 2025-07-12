package handlers


import (
    "database/sql"
    "net/http"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
    mux.HandleFunc("/", RootHandler)
    mux.HandleFunc("/healthz", HealthHandler)
    mux.HandleFunc("/event", NewEventHandler(db))
	mux.HandleFunc("/events", GetEventsHandler(db))
}

// import (
//     "database/sql"
//     "net/http"

//     "github.com/freshpaint/hipaa-tracker/handlers"
// )

// func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
//     mux.HandleFunc("/", handlers.RootHandler)
//     mux.HandleFunc("/healthz", handlers.HealthHandler)
//     mux.HandleFunc("/event", handlers.NewEventHandler(db)) // accepts POST
// }