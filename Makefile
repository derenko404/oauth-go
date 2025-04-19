include .env

APP_NAME = go-auth
MAIN = ./main.go
BUILD_DIR = bin
BINARY = $(BUILD_DIR)/$(APP_NAME)

MIGRATIONS_DIR = migrations
TARGET_VERSION ?= 1
DB_URL = postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable

.PHONY: all build run test clean format migrate-up migrate-down migrate-create

all: build

run:
	go run $(MAIN)

start-dev:
	docker-compose up -d && air

build:
	go build -o $(BINARY) $(MAIN)

test:
	go test ./...

clean:
	rm -rf $(BUILD_DIR)

format:
	gofmt -s -w ./

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	migrate -path migrations -database $(DB_URL) up $(TARGET_VERSION)

migrate-down:
	migrate -path migrations -database $(DB_URL) down $(TARGET_VERSION)

swagger:
	swag init