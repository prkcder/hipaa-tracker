# ======================
# Frontend Builder
# ======================
FROM node:22-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm install
COPY web/ ./
RUN npm run build

# ======================
# Backend Builder
# ======================
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app

# copying the file needed for modul initialization first
COPY cmd/main.go ./cmd/
COPY go.mod go.sum* ./
# # if i have internal packages
# COPY internal/ ./internal/ 

# Initialize module and get dependencies
RUN if [ ! -f go.mod ]; then \
      go mod init github.com/freshpaint/hipaa-tracker; \
    fi && \
    go mod tidy && \
    go mod download && \
    go mod verify  # Added verification for production-grade safety

# orignial below 
# COPY go.mod go.sum ./
# RUN go mod download
    
COPY . .
# Build
RUN go build -o /app/main ./cmd/main.go

# ======================
# Development Stage 
# ======================
# FROM backend-builder AS dev

# # Install air and frontend dev tools
# RUN go install github.com/cosmtrek/air@latest && \
#     apk add --no-cache curl

# # Set working directory to where your main.go is
# WORKDIR /app/cmd

# # Start air (will watch for changes)
# CMD ["air"]
# # Start both backend (air) and frontend (vite) dev servers
# CMD ["sh", "-c", "cd web && npm run dev & air -c .air.toml"]

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

# # AWS environment variables
# ENV AWS_REGION=us-east-1
# ENV AWS_ACCESS_KEY_ID=""
# ENV AWS_SECRET_ACCESS_KEY=""

COPY --from=backend-builder /app/main .
COPY --from=frontend-builder /app/frontend/dist ./static

EXPOSE 8080
CMD ["./main"]