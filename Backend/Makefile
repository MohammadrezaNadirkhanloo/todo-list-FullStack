-include .env
export

.DEFAULT_GOAL := help
SHELL := /bin/sh

APP_NAME    := todo
MAIN_PKG    := ./cmd/api
MIGRATE_DSN ?= postgres://$(APP_POSTGRES_USER):$(APP_POSTGRES_PASSWORD)@localhost:5432/$(APP_POSTGRES_DBNAME)?sslmode=disable


.PHONY: help
help: 
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## run dev 
	go run $(MAIN_PKG)

.PHONY: tidy
tidy: ## run tidy
	go mod tidy

.PHONY: migrate-up
migrate-up: ## up migrate
	migrate -path migrations -database "$(MIGRATE_DSN)" up

.PHONY: migrate-down
migrate-down: ## down migrate
	migrate -path migrations -database "$(MIGRATE_DSN)" down 1

.PHONY: migrate-create
migrate-create: ## create migrate
	@test -n "$(name)" || { echo "استفاده: make migrate-create name=add_orders"; exit 1; }
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: migrate-version
migrate-version: ## version migrate
	migrate -path migrations -database "$(MIGRATE_DSN)" version

