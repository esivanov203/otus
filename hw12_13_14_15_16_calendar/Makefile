BIN_CALENDAR := ./bin/calendar
BIN_SCHEDULER := ./bin/scheduler
BIN_SENDER := ./bin/sender

RELEASE ?= develop
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -s -w \
    -X main.release=$(RELEASE) \
    -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) \
    -X main.gitHash=$(GIT_HASH)

export RELEASE
export LDFLAGS

MIGRATIONS_DIR := ./migrations
ENV_FILE := ./.env

PROTO_DIR := ./proto
PROTO_OUT := ./proto

# ==============================
# Build
# ==============================
build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

# ==============================
# Run
# ==============================
run: build
	$(BIN_CALENDAR) --config=./configs/config.yaml

run-scheduler: build
	$(BIN_SCHEDULER) --config=./configs/scheduler_config.yaml

run-sender: build
	$(BIN_SENDER) --config=./configs/sender_config.yaml

# ==============================
# Protobuf
# ==============================
generate:
	protoc -I=$(PROTO_DIR) --go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

# ==============================
# Docker
# ==============================
build-img:
	docker build --build-arg=LDFLAGS="$(LDFLAGS)" -t calendar-app -f build/Dockerfile.calendar .
	docker build --build-arg=LDFLAGS="$(LDFLAGS)" -t calendar-scheduler -f build/Dockerfile.scheduler .
	docker build --build-arg=LDFLAGS="$(LDFLAGS)" -t calendar-sender -f build/Dockerfile.sender .

run-img: build-img
	docker run $(DOCKER_IMG)

# ==============================
# Version
# ==============================
version: build
	$(BIN_CALENDAR) version
	$(BIN_SCHEDULER) version
	$(BIN_SENDER) version

# ==============================
# Testing / Lint
# ==============================
test:
	go test -race ./internal/...

lint:
	golangci-lint run ./...

# ==============================
# Database migrations
# ==============================
# use as: make migrate-create NAME=create_events_table
migrate-create:
ifndef NAME
	$(error NAME is not set, e.g. make migrate-create NAME=create_events_table)
endif
	@echo "Creating new migration $(NAME)..."
	@goose -dir $(MIGRATIONS_DIR) create $(NAME) sql

migrate-up:
	@echo "Applying migrations..."
	@export $$(grep -v '^#' $(ENV_FILE) | xargs) && \
	goose -dir $(MIGRATIONS_DIR) postgres "user=$$POSTGRES_USER password=$$POSTGRES_PASSWORD dbname=$$POSTGRES_DB host=localhost port=$$POSTGRES_PORT sslmode=disable" up

migrate-down:
	@echo "Rolling back last migration..."
	@export $$(grep -v '^#' $(ENV_FILE) | xargs) && \
	goose -dir $(MIGRATIONS_DIR) postgres "user=$$POSTGRES_USER password=$$POSTGRES_PASSWORD dbname=$$POSTGRES_DB host=localhost port=$$POSTGRES_PORT sslmode=disable" down

# ==============================
# Release
# ==============================
up:
	@echo "****** Release $(RELEASE) is building *****"
	@COMPOSE_BAKE=1 docker-compose up -d --build
	@echo "****** Release $(RELEASE) has started *****"

down:
	@docker-compose down
	@docker image prune -f

# ==============================
# Integration Tests
# ==============================
integration-tests:
	@bash -c '\
		set -e; \
	  	echo "************ Integration tests ************"; \
		COMPOSE_BAKE=1 docker-compose -f docker-compose.yaml -f docker-compose.integration.override.yaml up -d --build; \
		docker wait integration_tests > /dev/null; \
		docker-compose -f docker-compose.yaml -f docker-compose.integration.override.yaml logs integration-tests; \
		CODE=$$(docker inspect -f "{{.State.ExitCode}}" integration_tests); \
		echo "************ Tests finish with exit code = $$CODE ************"; \
		docker-compose -f docker-compose.yaml -f docker-compose.integration.override.yaml down -v; \
		docker image prune -f; \
		exit $$CODE \
	'

# ==============================
# Phony targets
# ==============================
.PHONY: build run run-scheduler run-sender generate build-img run-img version test lint \
        migrate-create migrate-up migrate-down up down integration-tests
