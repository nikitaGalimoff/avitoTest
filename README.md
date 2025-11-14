# PR Reviewer Assignment Service

Сервис для автоматического назначения ревьюеров на Pull Request'ы с управлением командами и участниками.

## Архитектура

Проект реализован с использованием чистой архитектуры (Clean Architecture):

- **domain** - доменные модели и интерфейсы репозиториев
- **repository** - реализация репозиториев для работы с PostgreSQL
- **usecase** - бизнес-логика приложения
- **handler** - HTTP handlers для обработки запросов
- **config** - конфигурация приложения

## Требования

- Go 1.21+
- Docker и Docker Compose
- PostgreSQL 15+ (запускается через Docker Compose)

## Быстрый старт

### Windows

#### PowerShell (рекомендуется для Windows)

```powershell
# Запуск всего проекта одной командой (сборка, миграции, запуск)
.\build.ps1

# Или явно
.\build.ps1 up

# Просмотр логов
.\build.ps1 logs

# Остановка
.\build.ps1 down
```

#### Batch файл (альтернатива для Windows)

```cmd
# Запуск всего проекта одной командой
build.bat

# Или явно
build.bat up

# Просмотр логов
build.bat logs

# Остановка
build.bat down
```

### Linux/macOS

#### Make (рекомендуется)

```bash
# Запуск всего проекта одной командой (сборка, миграции, запуск)
make

# Или явно
make up

# Просмотр логов
make logs

# Остановка
make down
```

Сервис будет доступен на `http://localhost:8080`

**Примечание:** Миграции применяются автоматически при первом запуске PostgreSQL через Docker Compose.

### Альтернативный запуск через Docker Compose

```bash
# Запуск всех сервисов
docker-compose up -d

# Просмотр логов
docker-compose logs -f app

# Остановка
docker-compose down
```

### Локальный запуск

1. Убедитесь, что PostgreSQL запущен и доступен
2. Установите зависимости:
```bash
go mod download
```

3. Запустите миграции (они автоматически применятся при первом запуске PostgreSQL через Docker)

4. Запустите приложение:
```bash
make run
# или
go run ./cmd/server
```

## Конфигурация

Конфигурация задается через переменные окружения:

- `SERVER_PORT` - порт сервера (по умолчанию: 8080)
- `DB_HOST` - хост БД (по умолчанию: localhost)
- `DB_PORT` - порт БД (по умолчанию: 5432)
- `DB_USER` - пользователь БД (по умолчанию: postgres)
- `DB_PASSWORD` - пароль БД (по умолчанию: postgres)
- `DB_NAME` - имя БД (по умолчанию: avitotest)
- `ADMIN_TOKEN` - токен администратора (по умолчанию: admin-token)
- `USER_TOKEN` - токен пользователя (по умолчанию: user-token)

## API Endpoints

### Teams

- `POST /team/add` - Создать команду с участниками
- `GET /team/get?team_name=<name>` - Получить команду (требует авторизации)

### Users

- `POST /users/setIsActive` - Установить флаг активности пользователя (требует админский токен)
- `GET /users/getReview?user_id=<id>` - Получить PR'ы пользователя (требует авторизации)

### Pull Requests

- `POST /pullRequest/create` - Создать PR и назначить ревьюверов (требует админский токен)
- `POST /pullRequest/merge` - Пометить PR как MERGED (требует админский токен)
- `POST /pullRequest/reassign` - Переназначить ревьювера (требует админский токен)

### Health

- `GET /health` - Проверка здоровья сервиса

## Авторизация

Для защищенных эндпоинтов необходимо передавать токен в заголовке `Authorization`:

```
Authorization: Bearer <token>
```

Используйте `ADMIN_TOKEN` для админских операций и `USER_TOKEN` для пользовательских.

## Примеры использования

### Создание команды

```bash
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "backend",
    "members": [
      {"user_id": "u1", "username": "Alice", "is_active": true},
      {"user_id": "u2", "username": "Bob", "is_active": true}
    ]
  }'
```

### Создание PR

```bash
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer admin-token" \
  -d '{
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search",
    "author_id": "u1"
  }'
```

### Получение PR пользователя

```bash
curl -X GET "http://localhost:8080/users/getReview?user_id=u2" \
  -H "Authorization: Bearer user-token"
```

## Makefile команды

- `make` или `make up` - **Запуск всего проекта** (сборка, миграции, запуск) - основная команда
- `make down` - Остановка проекта
- `make logs` - Просмотр логов приложения
- `make restart` - Перезапуск сервисов
- `make clean-all` - Полная очистка (включая volumes)
- `make build` - Собрать приложение локально
- `make run` - Запустить приложение локально
- `make test` - Запустить тесты
- `make help` - Показать справку по командам

## Реализованные функции

✅ Все основные эндпоинты согласно OpenAPI спецификации
✅ Автоматическое назначение до 2 ревьюверов при создании PR
✅ Переназначение ревьюверов
✅ Идемпотентная операция merge
✅ Управление активностью пользователей
✅ Фильтрация неактивных пользователей при назначении
✅ Защита от изменения ревьюверов после merge
✅ Чистая архитектура с разделением слоев

## Принятые решения

1. **Аутентификация**: Реализована упрощенная аутентификация через токены в заголовке Authorization. В production следует использовать полноценную JWT аутентификацию.

2. **Миграции**: Миграции автоматически применяются при первом запуске PostgreSQL через Docker Compose (файлы из `migrations/` копируются в `/docker-entrypoint-initdb.d/`).

3. **Хранение ревьюверов**: Список ревьюверов хранится в виде JSONB массива в PostgreSQL для удобства работы с JSON операциями.

4. **Выбор ревьюверов**: Используется случайный выбор из доступных кандидатов с использованием `math/rand`.

5. **Обработка ошибок**: Все доменные ошибки оборачиваются в структурированный формат согласно OpenAPI спецификации.

## Структура проекта

```
.
├── cmd/
│   └── server/
│       └── main.go          # Точка входа приложения
├── internal/
│   ├── config/              # Конфигурация
│   ├── domain/              # Доменные модели и интерфейсы
│   ├── handler/             # HTTP handlers
│   ├── repository/          # Репозитории для работы с БД
│   └── usecase/             # Бизнес-логика
├── migrations/              # SQL миграции
├── docker-compose.yml      # Docker Compose конфигурация
├── Dockerfile              # Docker образ приложения
├── Makefile                # Команды для сборки и запуска
├── go.mod                  # Go зависимости
└── README.md               # Документация
```

