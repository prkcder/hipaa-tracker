# HIPAA Tracker 


## Purpose

This small project is designed to mimic parts of Freshpaint's potential technical stack and workflows. It's a learning exercise that demonstrates:
- Sanitizing sensitive data from event payloads
- Ingesting and storing structured events
- Forwarding events to downstream systems
- Viewing collected events through a simple frontend

The goal is to learn by building something similar to what Support/Integration Engineers or Data Privacy roles might work with.

---

## Project Structure

```
Old Planning Tree Structure
hipaa-tracker/
├── cmd/
│   └── main.go                    # Entry point
├── internal/
│   └── tracker/
│       ├── handler.go             # HTTP Handlers
│       ├── sanitizer.go           # Strip sensitive data
│       ├── forwarder.go           # Simulate sending to analytics destination
│       └── storage.go             # DB interface
├── config/
│   └── sensitive_fields.yaml      # Configurable field list to sanitize
├── web/
│   ├── index.html                 # Simple frontend UI
│   ├── events.js                  # JS to emit events
│   └── style.css
├── scripts/
│   └── seed_db.sql                # Optional DB setup
├── test/
│   ├── sanitizer_test.go
│   └── handler_test.go
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── .env                           # For DB creds and config
└── README.md


New Tree Structure
hipaa-tracker/
├── Dockerfile
├── README.md
├── cmd
│   ├── main.go
│   └── mock
│       └── mock_server.go
├── config
├── docker-compose.yml
├── event.json
├── go.mod
├── go.sum
├── internal
│   ├── db
│   │   ├── event_repository.go
│   │   └── postgres.go
│   ├── forwarder
│   │   └── forwarder.go
│   ├── handlers
│   │   ├── event.go
│   │   ├── root.go
│   │   └── routes.go
│   ├── middleware
│   │   └── cors.go
│   ├── models
│   │   └── event.go
│   ├── sanitize
│   │   └── sanitizer.go
│   └── storage
│       └── storage.go
├── scripts
│   └── seed_db.sql
├── sensitive_fields.yaml
├── test
│   ├── handler_event_test.go
│   ├── handlers_routes_test.go
│   ├── sanitize_load_test.go
│   └── sanitizer_test.go
└── web
    ├── README.md
    ├── eslint.config.mjs
    ├── next-env.d.ts
    ├── next.config.ts
    ├── package-lock.json
    ├── package.json
    ├── postcss.config.mjs
    ├── public
    │   ├── file.svg
    │   ├── globe.svg
    │   ├── next.svg
    │   └── window.svg
    ├── src
    │   ├── app
    │   │   ├── favicon.ico
    │   │   ├── globals.css
    │   │   ├── layout.tsx
    │   │   └── page.tsx
    │   └── components
    │       ├── SubmitEvent.tsx
    │       └── ViewEvents.tsx
    └── tsconfig.json
```


---

## 🛠️ Tech Stack

- Go (Golang)
- PostgreSQL
- Docker + Docker Compose
- YAML (for config)
- HTML/CSS/JS frontend
- `slog` for structured logging

---

## 🧪 Running Locally

### Prerequisite

Before running the application, you need to set up your environment variables:

Create your .env file:
```bash
cp .env.example .env
```

Update the .env file with your configuration:
Open the .env file and update the values according to your setup. The file should contain database credentials and other configuration variables as shown in .env.example.

### Option 1: Docker Compose

```bash
docker-compose up --build -d
```

This sets up the Go API and Postgres database. You can then access:
- `http://localhost:8080` — Frontend UI
- `http://localhost:8080/event` — API endpoint to POST events
- `http://localhost:8080/events` — JSON list of events

### Option 2: Run manually

1. Ensure Postgres is running
2. Set your `.env` variables (see Prerequisites above)
3. Run:
```bash
go run cmd/main.go
```

---

## 📤 Sending an Event

```bash
curl -X POST http://localhost:8080/event \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "signup",
    "payload": {
      "email": "test@example.com",
      "username": "tester123",
      etc...
    }
  }'
```

You’ll get a response like:

```json
{
  "id": 1,
  "created_at": "2025-07-18T19:00:00Z",
  "sanitized": true
}
```

---

## 🔍 What Gets Sanitized?

The list of fields (e.g., `email`, `ssn`, `password`) is defined in:

```
config/sensitive_fields.yaml
```

You can update this file to simulate different privacy policies.

---

## 🧪 Running Tests

```bash
go test ./...
```

Tests cover:
- Payload sanitization
- Valid/invalid HTTP requests
- DB + forwarding mock behavior

---

## 📡 Forwarding

Events are also forwarded (simulated) to a mock analytics destination via the `internal/forwarder` package. This behavior can be modified or extended to use real services.

---

## 📝 Notes

- This is **not production code**, but a learning tool (**could become real**)
- Good starting point for understanding event ingestion, processing, and data hygiene
- Inspired by different job descriptions from Freshpaint

---

## ✨ Future Ideas

- Add authentication layer
- Forward to AWS / webhook endpoint
- Add UI filters/search for events
- CLI tools for managing event schema

---
