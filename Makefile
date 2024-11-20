# goose binary
GOOSE_BIN=goose

# Migration directory
MIGRATION_DIR=migrations

# Load .env variables file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

run: help

.PHONY: help
help:
	@echo "Makefile for managing Project tools"
	@echo 
	@echo "Tools:"
	@echo "  setup				- Download and install all necessary tools"
	@echo "  install-goose			- Download and install goose"
	@echo "  install-swag			- Download and install swag"
	@echo 
	@echo "Actions:"
	@echo "  test				- Run tests"
	@echo "  swagger			- Generate Swagger configuration files"
	@echo 
	@echo "Database:"
	@echo "  rollback			- Rollback the last migration"
	@echo "  rebuild			- Rebuild migrations"
	@echo "  reset				- Rollback all the migrations"
	@echo "  status			- Show the current migration status"
	@echo "  migrate			- Apply all pending migrations"
	@echo "  create-migration		- Create a new migration with a user-provided name"

.PHONY: install-goose
install-goose:
	@echo "Downloading and installing goose"
	go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "To use the tool please add $GOPATH/bin to the path"

.PHONY: install-swag
install-swag:
	@echo "Downloading and installing swag"
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "To use the tool please add $GOPATH/bin to the path"

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

.PHONY: swagger
swagger:
	@echo "Generating..."
	@swag init -g cmd/identityService/main.go -o docs/

.PHONY: rollback
rollback:
	@echo "Rolling back the last migration..."
	@$(GOOSE_BIN) -dir $(MIGRATION_DIR) postgres "user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) sslmode=disable" down

.PHONY: rebuild
rebuild: reset migrate

.PHONY: reset
reset:
	@echo "Rolling back all migrations..."
	@$(GOOSE_BIN) -dir $(MIGRATION_DIR) postgres "user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) sslmode=disable" reset

.PHONY: status
status:
	@echo "Checking migration status..."
	@$(GOOSE_BIN) -dir $(MIGRATION_DIR) postgres "user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) sslmode=disable" status

.PHONY: migrate
migrate:
	@echo "Applying migrations..."
	@$(GOOSE_BIN) -dir $(MIGRATION_DIR) postgres "user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) sslmode=disable" up

.PHONY: create-migration
create-migration:
	@read -p "Enter migration name: " name; \
	@echo "Creating migration with name: $$name"; \
	@$(GOOSE_BIN) -dir $(MIGRATION_DIR) create $$name sql
