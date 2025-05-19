include .env
export

.PHONY: build-worker
MIGRATE_VERSION := v4.18.2
MIGRATE_BIN := ./bin/migrate
# Compose connection string
DB_CONN_STRING := postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Download migrate CLI into ./bin
setup-migrate:
	mkdir -p ./bin
	curl -L https://github.com/golang-migrate/migrate/releases/download/$(MIGRATE_VERSION)/migrate.$(shell uname -s | tr A-Z a-z)-amd64.tar.gz | tar xz -C ./bin
	chmod +x $(MIGRATE_BIN)
	@echo "✅ migrate installed at $(MIGRATE_BIN)"

run-api:
	go run ./cmd/api

run-worker:
	. ./scripts/.radar-oficial/bin/activate && \
	go run ./cmd/worker

migrate:
	go run ./cmd/migrate

# Rollback the last migration
rollback-migration:
	$(MIGRATE_BIN) -path cmd/migrate/migrations -database "$(DB_CONN_STRING)" -verbose down 1

create-migration:
	@read -p "Enter migration name: " name; \
	$(MIGRATE_BIN) create -ext sql -dir cmd/migrate/migrations -seq $$name

psql:
	@PGPASSWORD=$(DB_PASS) psql -p $(DB_PORT) -U $(DB_USER)  -h $(DB_HOST) $(DB_NAME)