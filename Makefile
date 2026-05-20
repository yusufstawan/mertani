-include .env
export

APP_PORT ?= 8080
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= mertani
DB_SSLMODE ?= disable
PGWEB_PORT ?= 8081
STEPS ?= 1

.PHONY: help run test fmt compose-up compose-down db-up pgweb-up logs migrate-up migrate-down migrate-steps migrate-force migrate-version

help:
	@printf "%s\n" \
		"Available commands:" \
		"  make run                         Run API server" \
		"  make test                        Run Go tests" \
		"  make fmt                         Format Go files" \
		"  make compose-up                  Start PostgreSQL and pgweb" \
		"  make compose-down                Stop Docker Compose services" \
		"  make db-up                       Start PostgreSQL only" \
		"  make pgweb-up                    Start pgweb only" \
		"  make logs                        Follow Docker Compose logs" \
		"  make migrate-up                  Apply all pending migrations" \
		"  make migrate-down [STEPS=1]       Roll back migration steps" \
		"  make migrate-steps STEPS=1        Move migration version by N steps" \
		"  make migrate-force VERSION=2      Force migration version metadata" \
		"  make migrate-version             Show current migration version"

run:
	go run ./cmd/api

test:
	go test ./...

fmt:
	@gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

compose-up:
	docker compose up -d

compose-down:
	docker compose down

db-up:
	docker compose up -d postgres

pgweb-up:
	docker compose up -d pgweb

logs:
	docker compose logs -f

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down $(STEPS)

migrate-steps:
	go run ./cmd/migrate steps $(STEPS)

migrate-force:
	@test -n "$(VERSION)" || (echo "usage: make migrate-force VERSION=2"; exit 1)
	go run ./cmd/migrate force $(VERSION)

migrate-version:
	go run ./cmd/migrate version
