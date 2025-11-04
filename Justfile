# Load environment variables from .env file
set dotenv-load

# goose binary
goose_bin := "goose"

# Migration directory
migration_dir := "migrations"

default:
    @just --list

# Run the application
serve:
    go run cmd/identity-service/main.go

# Run the application with nix
serve-nix:
    nix run .#default

# Download and install all necessary tools
setup: install-goose install-swag

# Download and install goose
install-goose:
	@echo "Downloading and installing goose"
	go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "To use the tool please add $GOPATH/bin to the path"

# Download and install swag
install-swag:
	@echo "Downloading and installing swag"
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "To use the tool please add $GOPATH/bin to the path"

# Run tests
test:
	@echo "Running tests..."
	@go test -cover ./...

# Generate Coverage Report
coverage:
	@echo "Running tests..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

# Generate Swagger configuration files
swagger:
	@echo "Generating..."
	@swag init -g cmd/identity-service/main.go -o docs/

# Rollback the last migration
rollback:
	@echo "Rolling back the last migration..."
	@{{goose_bin}} -dir {{migration_dir}} postgres "user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} host=${POSTGRES_HOST} port=${POSTGRES_PORT} sslmode=disable" down

# Rebuild migrations
rebuild: reset migrate

# Rollback all migrations
reset:
	@echo "Rolling back all migrations..."
	@{{goose_bin}} -dir {{migration_dir}} postgres "user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} host=${POSTGRES_HOST} port=${POSTGRES_PORT} sslmode=disable" reset

# Show the current migration status
status:
	@echo "Checking migration status..."
	@{{goose_bin}} -dir {{migration_dir}} postgres "user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} host=${POSTGRES_HOST} port=${POSTGRES_PORT} sslmode=disable" status

# Apply all pending migrations
migrate:
	@echo "Applying migrations..."
	@{{goose_bin}} -dir {{migration_dir}} postgres "user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} host=${POSTGRES_HOST} port=${POSTGRES_PORT} sslmode=disable" up

# Create a new migration with a user-provided name
create-migration name:
	@echo "Creating migration with name: {{name}}"
	@{{goose_bin}} -dir {{migration_dir}} create {{name}} sql
