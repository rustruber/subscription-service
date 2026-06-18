# ============================================================
# Makefile для subscription-service
# ============================================================

.PHONY: help init run build test migrate-up migrate-down \
        docker-up docker-down swagger tidy clean

# ------------------------- HELP -----------------------------
help:  ## Показать все команды
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ------------------------- INIT -----------------------------
init:  ## Инициализация проекта (go mod init + tidy)
	go mod init subscription-service
	go mod tidy

# ------------------------- BUILD ----------------------------
build: ## Собрать бинарник
	go build -o bin/subscription-server cmd/server/main.go

run:   ## Запустить сервис локально
	go run cmd/server/main.go

test:  ## Запустить все тесты
	go test -v ./...

# ------------------------- MIGRATIONS -----------------------
migrate-up:   ## Применить миграции
	psql "postgres://postgres:postgres@localhost:5432/subscription?sslmode=disable" \
		-f migrations/001_create_subscriptions_table.up.sql

migrate-down: ## Откатить миграции
	psql "postgres://postgres:postgres@localhost:5432/subscription?sslmode=disable" \
		-f migrations/001_create_subscriptions_table.down.sql

# ------------------------- DOCKER ---------------------------
docker-up:   ## Запустить всё через Docker Compose
	docker-compose up -d

docker-down: ## Остановить Docker Compose
	docker-compose down

docker-logs: ## Показать логи
	docker-compose logs -f

# ------------------------- SWAGGER --------------------------
swagger: ## Сгенерировать Swagger документацию
	swag init -g cmd/server/main.go -o api/

# ------------------------- UTILS ----------------------------
tidy:  ## Подчистить зависимости
	go mod tidy

clean: ## Очистить бинарники
	rm -rf bin/
	rm -rf tmp/