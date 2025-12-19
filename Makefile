.PHONY: init build start mocks lint test checks stop compush stage
MAKEFLAGS += --no-print-directory
GIT_BRANCH := $(shell git branch --show-current)

CHECK_EMOJI := ✅
ERROR_EMOJI := ❌
INFO_EMOJI := ℹ️
ARROW_UP := ⬆️

PROJECT_NAME ?= todo_task

MIGRATION_NAME ?= migration

MIGRATIONS_PATH_POSTGRES := migrations/postgres

DB_URL_POSTGRES := postgres://admin:password@localhost:5432/todo_task?sslmode=disable

help:
	@echo "Available commands:"
	@echo "  build      - Build the binary file"
	@echo "  init       - Initialize the project"
	@echo "  start      - Start the application"
	@echo "  mocks      - Generate mocks"
	@echo "  lint       - Run linters"
	@echo "  test       - Run tests"
	@echo "  checks     - Run all checks (mocks, lint, test)"
	@echo "  stop       - Stop containers"
	@echo "  compush    - Run pre-commit checks, commit, and push changes"
	@echo "  stage      - Stage changes and push to Git"

build:
	@echo "$(INFO_EMOJI) Building the binary file..."
	@GOOS=linux CGO_ENABLED=0 go build -o portal ./cmd/main.go
	@echo "$(CHECK_EMOJI) Project built successfully!"

init:
	@echo "$(INFO_EMOJI) Initializing project..."
	@docker-compose up -d
	@echo "$(CHECK_EMOJI) Project initialized successfully!"

start:
	@echo "$(INFO_EMOJI) Starting the application..."
	@docker-compose up -d
	@trap 'exit 0' SIGINT; go run $(CURDIR)/cmd/main.go
	@echo "$(CHECK_EMOJI) Application started successfully!"

mocks:
	@echo "$(INFO_EMOJI) Generating mocks..."
	@(go generate ./... > mocks.log 2>&1 || \
		(cat mocks.log && echo "$(ERROR_EMOJI) Error generating mocks! Check logs $(ARROW_UP)" && exit 1))
	@rm -f mocks.log
	@echo "$(CHECK_EMOJI) All mocks were generated successfully!"

lint:
	@echo "$(INFO_EMOJI) Running linters..."
	@(golangci-lint run ./... > lints.log 2>&1 \
		|| (cat lints.log && echo "$(ERROR_EMOJI) Linter found issues! Check logs $(ARROW_UP)" && exit 1))
	@rm -f lints.log
	@echo "$(CHECK_EMOJI) No lint errors found!"

test:
	@echo "$(INFO_EMOJI) Running tests..."
	@(go test ./... > tests.log 2>&1 || \
		(cat tests.log && echo "$(ERROR_EMOJI) Tests failed! Check logs $(ARROW_UP)" && exit 1))
	@rm -f tests.log
	@echo "$(CHECK_EMOJI) All tests PASSED!"

checks: mocks lint test

stop:
	@echo "$(INFO_EMOJI) Stopping containers..."
	@docker-compose down
	@echo "$(CHECK_EMOJI) Containers stopped successfully!"

compush:
	@echo "$(INFO_EMOJI) Running pre-commit checks..."
	@$(MAKE) checks
	@echo "$(CHECK_EMOJI) Pre-commit checks passed! Moving to staging..."
	@$(MAKE) stage

stage:
	@echo "$(INFO_EMOJI) Staging changes..."
	@git add .
	@git commit -m "$(m)"
	@git push origin $(GIT_BRANCH)
	@echo "$(CHECK_EMOJI) Changes pushed to $(GIT_BRANCH)!"

proto-gen:
	@echo "$(INFO_EMOJI) Generating protobuf code..."
	@protoc \
		--proto_path= \
		--go_out=./api/grpc \
		--go_opt=module=gl.iteco.com/technology/go_services/iteco_notifgatev2/$(PROJECT_NAME)/api/grpc \
		--go-grpc_out=./api/grpc \
		--go-grpc_opt=module=gl.iteco.com/technology/go_services/iteco_notifgatev2/$(PROJECT_NAME)/api/grpc \
		--experimental_allow_proto3_optional \
		./api/$(PROJECT_NAME)_messages.proto \
		./api/$(PROJECT_NAME)_service.proto
	@echo "$(CHECK_EMOJI) Protobuf code generated successfully!"


m-pg-create:
	@echo "$(INFO_EMOJI) Creating new migration for postgres '$(MIGRATION_NAME)'..."
	@goose -dir $(MIGRATIONS_PATH_POSTGRES) create $(MIGRATION_NAME) sql
	@echo "$(CHECK_EMOJI) Migration files created!"

m-pg-up:
	@echo "$(INFO_EMOJI) Applying migrations for postgres UP..."
	@goose -dir $(MIGRATIONS_PATH_POSTGRES) postgres "$(DB_URL_POSTGRES)" up
	@echo "$(CHECK_EMOJI) Migrations applied!"

m-pg-down:
	@echo "$(INFO_EMOJI) Rolling back LAST migration for postgres..."
	@goose -dir $(MIGRATIONS_PATH_POSTGRES) postgres "$(DB_URL_POSTGRES)" down
	@echo "$(CHECK_EMOJI) Rolled back one migration!"

m-pg-down-all:
	@echo "$(INFO_EMOJI) Rolling back ALL migrations for postgres..."
	@goose -dir $(MIGRATIONS_PATH_POSTGRES) postgres "$(DB_URL_POSTGRES)" reset
	@echo "$(CHECK_EMOJI) Rolled back all migrations!"

m-pg-status:
	@goose -dir $(MIGRATIONS_PATH_POSTGRES) postgres "$(DB_URL_POSTGRES)" status