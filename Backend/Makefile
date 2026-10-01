.DEFAULT_GOAL := help
SHELL := /bin/sh

APP_NAME    := goclean
MAIN_PKG    := ./cmd/api

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

