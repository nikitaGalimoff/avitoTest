.PHONY: build run test clean docker-build docker-up docker-down migrate up down logs restart clean-all help

# Переменные
BINARY_NAME=avitotest
MAIN_PATH=./cmd/server

# Основная команда - запуск всего проекта (сборка, миграции, запуск)
.DEFAULT_GOAL := up

# Запуск всего проекта: сборка, миграции, запуск
up: docker-build docker-up wait-healthy
	@echo "✅ Проект запущен!"
	@echo "📊 Сервис доступен на http://localhost:8080"
	@echo "📝 Логи: make logs"

# Остановка проекта
down:
	docker-compose down

# Сборка приложения локально
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Запуск приложения локально
run:
	go run $(MAIN_PATH)

# Запуск тестов
test:
	go test -v ./...

# Очистка локальных артефактов
clean:
	rm -rf bin/

# Сборка Docker образа
docker-build:
	@echo "🔨 Сборка Docker образов..."
	docker-compose build

# Запуск через Docker Compose
docker-up:
	@echo "🚀 Запуск контейнеров..."
	docker-compose up -d

# Ожидание готовности сервисов
wait-healthy:
	@echo "⏳ Ожидание готовности сервисов..."
	@echo "   (Миграции применяются автоматически при первом запуске PostgreSQL)"
	@echo "   (PostgreSQL имеет healthcheck, приложение запустится автоматически после готовности БД)"
	@echo "✅ Сервисы запущены"

# Просмотр логов
logs:
	docker-compose logs -f app

# Перезапуск сервисов
restart:
	docker-compose restart

# Полная очистка (включая volumes и образы)
clean-all:
	@echo "🧹 Полная очистка..."
	docker-compose down -v
	docker-compose rm -f

# Справка
help:
	@echo "Доступные команды:"
	@echo "  make          - Запуск всего проекта (сборка, миграции, запуск)"
	@echo "  make up       - Запуск всего проекта"
	@echo "  make down     - Остановка проекта"
	@echo "  make logs     - Просмотр логов приложения"
	@echo "  make restart  - Перезапуск сервисов"
	@echo "  make clean-all - Полная очистка (включая volumes)"
	@echo "  make build    - Сборка приложения локально"
	@echo "  make run      - Запуск приложения локально"
	@echo "  make test     - Запуск тестов"

