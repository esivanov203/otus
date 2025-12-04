BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

MIGRATIONS_DIR := ./migrations
ENV_FILE := ./.env

PROTO_DIR="./proto"
PROTO_OUT="./proto"

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-scheduler:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/scheduler

build-sender:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(BIN) --config=./configs/config.yaml

run-scheduler: build-scheduler
	$(BIN) --config=./configs/scheduler_config.yaml

run-sender: build-sender
	$(BIN) --config=./configs/sender_config.yaml

generate:
	protoc -I=$(PROTO_DIR) --go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

lint:
	golangci-lint run ./...

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

.PHONY: build run build-img run-img version test lint migrate-up migrate-down install-lint-deps
