# ==============================================================================
# Configuration
# ==============================================================================

-include .env
export

.DEFAULT_GOAL := help
SHELL := /bin/sh

APP_NAME    := todo
MAIN_PKG    := ./cmd/api
MIGRATE_DSN ?= postgres://$(APP_POSTGRES_USER):$(APP_POSTGRES_PASSWORD)@localhost:5432/$(APP_POSTGRES_DBNAME)?sslmode=disable


# ==============================================================================
# Help
# ==============================================================================

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'


# ==============================================================================
# Development
# ==============================================================================

.PHONY: run
run: ## Run dev server
	go run $(MAIN_PKG)

.PHONY: tidy
tidy: ## Run go mod tidy
	go mod tidy


# ==============================================================================
# Database Migrations
# ==============================================================================

.PHONY: migrate-up
migrate-up: ## Apply all up migrations
	migrate -path migrations -database "$(MIGRATE_DSN)" up

.PHONY: migrate-down
migrate-down: ## Roll back the last migration
	migrate -path migrations -database "$(MIGRATE_DSN)" down 1

.PHONY: migrate-create
migrate-create: ## Create a new migration (usage: make migrate-create name=add_orders)
	@test -n "$(name)" || { echo "Usage: make migrate-create name=add_orders"; exit 1; }
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: migrate-version
migrate-version: ## Show current migration version
	migrate -path migrations -database "$(MIGRATE_DSN)" version


# ==============================================================================
# Docker 
# ==============================================================================

.PHONY: up
up: ## Start the core stack (app + postgres + redis)
	$(COMPOSE) up -d --build

.PHONY: down
down: ## Stop the stack
	$(COMPOSE) down

.PHONY: down-clean
down-clean: ## Stop the stack and remove data (⚠️ data will be deleted)
	$(COMPOSE) down -v

.PHONY: logs
logs: ## Follow logs
	$(COMPOSE) logs -f app

.PHONY: ps
ps: ## Show services status
	$(COMPOSE) ps