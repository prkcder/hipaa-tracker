package main

// import (
//     "fmt"
//     "log"
//     "net/http"
//     "os"
// )

// func main() {
//     port := os.Getenv("PORT")
//     if port == "" {
//         port = "8080"
//     }

// 	log.Println(("starting the app"))

//     http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//         fmt.Fprintln(w, "Backend is running!")
//     })

//     log.Printf("Listening on port %s...", port)
//     err := http.ListenAndServe(":" + port, nil)
//     if err != nil {
//         log.Fatal(err)
//     }
// }

// import "fmt"

// import "rsc.io/quote"

// func main() {
//     fmt.Println(quote.Go())
// 	fmt.Println("hello")
// }

// import (
//     "database/sql"
//     "fmt"
//     "log"
//     "net/http"
//     "os"

//     _ "github.com/lib/pq"
// )

// func main() {


// 	db, err := setupDatabase()

//     if err != nil {
//         log.Fatalf("Failed to connect to DB: %v", err)
//     }
// 	defer db.Close()


//     http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//         fmt.Fprintln(w, "Server is up and DB connected")
//     })

//     http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
//         fmt.Fprintln(w, " POST your event data here")
//     })

//     http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
//         fmt.Fprintln(w, "ok")
//     })

//     log.Println(" Starting server on :8080")
//     if err := http.ListenAndServe(":8080", nil); err != nil {
//         log.Fatalf("Server failed: %v", err)
//     }

// }


// func setupDatabase() (*sql.DB, error) {
// 	dbHost := os.Getenv("DB_HOST")
//     dbPort := os.Getenv("DB_PORT")
//     dbUser := os.Getenv("DB_USER")
//     dbPassword := os.Getenv("DB_PASSWORD")
//     dbName := os.Getenv("DB_NAME")

//     dsn := fmt.Sprintf(
//         "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
//         dbHost, dbPort, dbUser, dbPassword, dbName,
//     )

// 	db, err := sql.Open("postgres", dsn)
// 	if err != nil {
//         return nil, fmt.Errorf("failed to connect: %w", err)
//     }

//     err = db.Ping()
//     if err != nil {
//         log.Fatalf("Database is unreachable: %v", err)
//     }

//     log.Println("Connected to Postgres!")
//     return db, nil
	
// }



import (
    "log"
    "net/http"

    "github.com/freshpaint/hipaa-tracker/internal/db"
    "github.com/freshpaint/hipaa-tracker/internal/handlers"
    "github.com/freshpaint/hipaa-tracker/internal/sanitize"
    // "github.com/freshpaint/hipaa-tracker/internal/forwarder"

)


func main() {

    log.Println("starting tracker")

    if err := sanitize.LoadSensitiveFields("sensitive_fields.yaml"); err != nil {
        log.Fatalf("Failed to load sensitive fields config: %v", err)
    }

    database, err := db.Connect()
    if err != nil {
        log.Fatalf("failed to connect to DB: %v", err)
    }
    defer database.Close()

    mux := http.NewServeMux()

    // mux.HandleFunc("/", handlers.RootHandler)
    // mux.HandleFunc("/healthz", handlers.HealthHandler)
    // mux.HandleFunc("/event", handlers.NewEventHandler(database)) // returns a func

    handlers.RegisterRoutes(mux, database)



    // mux := handlers.InitRoutes(db)

    log.Println("🚀 Server is listening on :8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatalf("Server failed: %v", err)
    }


}