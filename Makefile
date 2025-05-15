MIGRATIONS_DIR=./migrations
CONFIG=.env
DSN=$(shell yq e '.db_postgres.dsn' $(CONFIG))

.PHONY: migrate-up migrate-down migrate-status migrate-create


check-docker-postgres:
	docker compose exec postgres psql -U admin -d users

docker-up:
	docker compose up --build

migrate-create-%:
	goose -dir $(MIGRATIONS_DIR) create $(subst migrate-create-,,$@) sql

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DSN)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DSN)" down

migrate-reset:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DSN)" reset

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DSN)" status
test:
	go test -v -cover -coverpkg ./... ./...	