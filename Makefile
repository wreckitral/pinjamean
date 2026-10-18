TEST_DB_URL ?= postgres://postgres:postgres@localhost:5432/pinjamean_test?sslmode=disable

.PHONY: run up down test test-unit test-integration lint

run:
	docker compose up --build

up:
	docker compose up -d postgres
	@echo "Waiting for postgres to be ready..."
	@docker compose exec postgres sh -c 'until pg_isready -U postgres; do sleep 1; done'

down:
	docker compose down

test-unit:
	cd internal/loans && go test -race ./domain/... ./app/...

test-integration: up
	TEST_DATABASE_URL="$(TEST_DB_URL)" cd internal/loans && go test -race -v ./adapters/...

test: test-unit test-integration

lint:
	cd internal/common && go vet ./...
	cd internal/loans && go vet ./...
	cd internal/identity && go vet ./...
