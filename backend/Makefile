.PHONY: run test build migrate-up migrate-down migrate-status

run:
	go run ./cmd/api/

test:
	go test ./... -count=1

test-leads:
	go test ./internal/interfaces/http/... -count=1 -run 'LeadsHTTP'

test-activities:
	go test ./internal/interfaces/http/... -count=1 -run 'Activities'

test-contacts:
	go test ./internal/interfaces/http/... -count=1 -run 'Contacts'

build:
	go build -o bin/api ./cmd/api/

migrate-up:
	go run ./cmd/migrate/ up

migrate-down:
	go run ./cmd/migrate/ down

migrate-status:
	go run ./cmd/migrate/ status
