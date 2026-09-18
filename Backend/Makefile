# فرمان‌های متداول توسعه.
#
# اجرای `make` یا `make help` فهرست فرمان‌ها را نشان می‌دهد.

.DEFAULT_GOAL := help
SHELL := /bin/sh

APP_NAME    := goclean
MAIN_PKG    := ./cmd/api
BIN_DIR     := bin
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS     := -s -w -X main.version=$(VERSION)

COMPOSE     := docker compose -f deploy/docker-compose.yml
MIGRATE_DSN ?= postgres://$(APP_POSTGRES_USER):$(APP_POSTGRES_PASSWORD)@localhost:5432/$(APP_POSTGRES_DBNAME)?sslmode=disable

.PHONY: help
help: ## نمایش این راهنما
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# ─────────────────────────────── ساخت و اجرا ──────────────────────

.PHONY: build
build: ## ساخت باینری
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(APP_NAME) $(MAIN_PKG)
	@echo "ساخته شد: $(BIN_DIR)/$(APP_NAME)  (نسخه $(VERSION))"

.PHONY: run
run: ## اجرای برنامه به‌صورت محلی
	go run $(MAIN_PKG)

.PHONY: tidy
tidy: ## مرتب‌سازی وابستگی‌ها
	go mod tidy

# ─────────────────────────────── کیفیت کد ─────────────────────────

.PHONY: fmt
fmt: ## قالب‌بندی کد
	gofmt -w .
	@command -v goimports >/dev/null 2>&1 && goimports -w . || true

.PHONY: vet
vet: ## اجرای go vet
	go vet ./...

.PHONY: lint
lint: ## اجرای golangci-lint
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint نصب نیست. نصب:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; }
	golangci-lint run ./...

.PHONY: test
test: ## اجرای تست‌ها
	go test -race -count=1 ./...

.PHONY: cover
cover: ## گزارش پوشش تست
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "گزارش: coverage.html"

.PHONY: check
check: fmt vet lint test ## اجرای همه‌ی بررسی‌ها

# ─────────────────────────────── امنیت ────────────────────────────

.PHONY: audit
audit: ## بررسی آسیب‌پذیری وابستگی‌ها
	@command -v govulncheck >/dev/null 2>&1 || \
		go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

# ⚠️ این هدف را در CI اجرا کنید.
#
# پروژه‌های مبدأ هیچ بررسی آسیب‌پذیری نداشتند و به کتابخانه‌هایی
# وابسته بودند که دیگر پشتیبانی نمی‌شوند — از جمله
# github.com/golang-jwt/jwt نسخه‌ی ۳ که برای کتابخانه‌ی احراز هویت
# انتخاب بدی است.

.PHONY: secret-scan
secret-scan: ## جست‌وجوی راز در تاریخچه‌ی گیت
	@command -v gitleaks >/dev/null 2>&1 || { \
		echo "gitleaks نصب نیست: https://github.com/gitleaks/gitleaks"; exit 1; }
	gitleaks detect --source . --verbose

# ─────────────────────────────── مهاجرت ───────────────────────────

.PHONY: migrate-up
migrate-up: ## اجرای مهاجرت‌ها
	migrate -path migrations -database "$(MIGRATE_DSN)" up

.PHONY: migrate-down
migrate-down: ## بازگشت یک مهاجرت
	migrate -path migrations -database "$(MIGRATE_DSN)" down 1

.PHONY: migrate-create
migrate-create: ## ساخت مهاجرت جدید (name=... )
	@test -n "$(name)" || { echo "استفاده: make migrate-create name=add_orders"; exit 1; }
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: migrate-version
migrate-version: ## نمایش نسخه‌ی فعلی مهاجرت
	migrate -path migrations -database "$(MIGRATE_DSN)" version

# ─────────────────────────────── داکر ─────────────────────────────

.PHONY: up
up: ## بالا آوردن پشته‌ی هسته (app + postgres + redis)
	$(COMPOSE) up -d --build

.PHONY: down
down: ## توقف پشته
	$(COMPOSE) down

.PHONY: down-clean
down-clean: ## توقف پشته و حذف داده‌ها (⚠️ داده پاک می‌شود)
	$(COMPOSE) down -v

.PHONY: logs
logs: ## دنبال کردن لاگ‌ها
	$(COMPOSE) logs -f app

.PHONY: ps
ps: ## وضعیت سرویس‌ها
	$(COMPOSE) ps
