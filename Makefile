# Makefile

# Default target
start-build:
	docker compose -f docker-compose.dev.yml up --build

start-up:
	docker compose -f docker-compose.dev.yml up -d
start:
	docker compose -f docker-compose.dev.yml up --build -d

stop:
	docker compose -f docker-compose.dev.yml down

stop-clean:
	docker compose -f docker-compose.dev.yml down -v

logs:
	docker compose -f docker-compose.dev.yml logs -f

restart:
	docker compose -f docker-compose.dev.yml down && docker compose -f docker-compose.dev.yml up --build

# Optional: start prod
prod:
	docker compose -f docker-compose.yml up --build -d
