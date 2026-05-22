.PHONY: run test build migrate-up migrate-down migrate-status

run:
	go run ./cmd/api/

test:
	go test ./... -count=1

build:
	go build -o bin/api ./cmd/api/

migrate-up:
	go run ./cmd/migrate/ up

migrate-down:
	go run ./cmd/migrate/ down

migrate-status:
	go run ./cmd/migrate/ status
