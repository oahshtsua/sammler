# Load environment variables
include .env

## help: Show this help message
.PHONY: help
help:
	@echo 'Available commands:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/  /'

## run: Start the application server
.PHONY: run
run:
	go run ./cmd/web/ -port=${PORT} -dsn=${DB_DSN}

## sqlc: Generate code using sqlc
.PHONY: sqlc
sqlc:
	sqlc generate

## migration: Create a new SQL migration (make migration name=add_users_table)
.PHONY: migration
migration:
	goose sqlite3 ${DB_DSN} -dir=migrations create ${name} sql

## migrate: Apply all pending migrations to the database
.PHONY: migrate
migrate:
	goose sqlite3 ${DB_DSN} -dir=migrations up

