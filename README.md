hipaa-tracker/
├── cmd/
│   └── main.go                 # Entry point
├── internal/
│   ├── tracker/
│   │   ├── handler.go          # HTTP Handlers
│   │   ├── sanitizer.go        # Strip sensitive data
│   │   ├── forwarder.go        # Simulate sending to analytics destination
│   │   └── storage.go          # DB interface
├── config/
│   └── sensitive_fields.yaml   # Configurable field list to sanitize
├── web/
│   ├── index.html              # Simple frontend UI
│   ├── events.js               # JS to emit events
│   └── style.css
├── scripts/
│   └── seed_db.sql             # Optional DB setup
├── test/
│   ├── sanitizer_test.go
│   └── handler_test.go
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── .env                        # For DB creds and config
└── README.md



