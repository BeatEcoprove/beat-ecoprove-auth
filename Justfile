# Load environment variables from .env file
set dotenv-load

# goose binary
goose_bin := "goose"

# Migration directory
migration_dir := "migrations"

# Default recipe (runs when you type 'just')
default: help

# Show available recipes
help:
	@echo "Justfile for managing Project tools"
	@echo
	@echo "Tools:"
	@echo "  setup                         - Download and install all necessary tools"
	@echo "  install-goose                 - Download and install goose"
	@echo "  install-swag                  - Download and install swag"
	@echo
	@echo "Actions:"
	@echo "  test                          - Run tests"
	@echo "  coverage                      - Generate Coverage Report"
	@echo "  swagger                       - Generate Swagger configuration files"
	@echo
	@echo "Database:"
	@echo "  rollback                      - Rollback the last migration"
	@echo "  rebuild                       - Rebuild migrations"
	@echo "  reset                         - Rollback all the migrations"
	@echo "  status                        - Show the current migration status"
	@echo "  migrate                       - Apply all pending migrations"
	@echo "  create-migration              - Create a new migration with a user-provided name"

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
