.PHONY: help build run dev lint migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make dev         - Run the application in development mode"
	@echo "  make lint        - Run linter on the codebase"
	@echo "  make format      - Format the code and re-arrange imports"
	@echo "  make migrate-up  - Apply database migrations"
	@echo "  make migrate-down- Rollback database migrations"


build:
	go build -o bin/app ./cmd/api

run:
	go run ./cmd/api

dev:
	go run ./cmd/api

lint: format
	golangci-lint run ./...

format:
	@gofmt -s -w .
	@goimports -w .

migrate-up:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable" down

docker-up:
	docker compose -f docker/docker-compose.yml up -d

docker-down:
	docker compose -f docker/docker-compose.yml down
