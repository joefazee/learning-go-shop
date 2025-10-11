# Learning Go Shop

Learning Go Shop is a modern e-commerce backend written in Go, featuring REST and GraphQL APIs, background workers, and integrations with PostgreSQL, AWS services, and local file storage. The codebase is the hands-on companion for the Udemy course [Learn Golang Like Google Engineers Do: 100% Project-Based](https://www.udemy.com/course/draft/6836245/).

## Prerequisites
- Go 1.23 or newer
- PostgreSQL (default connection `postgres://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable`)
- Optional: Docker and Docker Compose for local dependencies
- Optional: AWS credentials or LocalStack when working with S3 uploads and event queues

Application settings are read from environment variables. You can create a `.env` file in the project root to override the defaults defined in `internal/config/config.go`.

## Run with Go Commands (No Make Required)
Clone the repository, then from the project root:

```bash
# install dependencies
go mod download

# run database migrations (requires https://github.com/golang-migrate/migrate)
migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable" up

# start the HTTP API
go run ./cmd/api
```

The notifier service can be started with:

```bash
go run ./cmd/notifier
```

## Additional Useful Commands
- Build binaries: `go build -o bin/api ./cmd/api` and `go build -o bin/notifier ./cmd/notifier`
- Format code: `gofmt -s -w .` and `goimports -w .`
- Lint (requires `golangci-lint`): `golangci-lint run ./...`
- Generate Swagger docs (requires `swag`): `swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal --exclude .git,docs,docker,db`
- Generate GraphQL code: `go run github.com/99designs/gqlgen generate`
- Run tests: `go test ./...`

## Docker Support
If you prefer to run dependencies with Docker, start the local stack with:

```bash
docker compose -f docker/docker-compose.yml up -d
```

Shut everything down when finished:

```bash
docker compose -f docker/docker-compose.yml down
```
