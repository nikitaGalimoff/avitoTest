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

- Docker и Docker Compose
- Go 1.25.4

Сервис будет доступен на `http://localhost:8080`

**Примечание:** Миграции применяются автоматически при первом запуске PostgreSQL через Docker Compose.

### Запуск через make

```bash
# Запуск всех сервисов
make run


# Остановка
make stop
```
## Конфигурация

Конфигурация задается через переменные окружения:

- `SERVER_PORT` - порт сервера (по умолчанию: 8080)
- `DB_HOST` - хост БД (по умолчанию: postgres)
- `DB_PORT` - порт БД (по умолчанию: 5432)
- `DB_USER` - пользователь БД (по умолчанию: postgres)
- `DB_PASSWORD` - пароль БД (по умолчанию: postgres)
- `DB_NAME` - имя БД (по умолчанию: avitotest)

В проекте используются значения по умолчанию, но можно добавить .env файл в проект и конфигурация будет задаваться в нем

## API Endpoints

### Teams

- `POST /team/add` - Создать команду с участниками
- `GET /team/get?team_name=<name>` - Получить команду

### Users

- `POST /users/setIsActive` - Установить флаг активности пользователя 
- `GET /users/getReview?user_id=<id>` - Получить PR'ы пользователя 

### Pull Requests

- `POST /pullRequest/create` - Создать PR и назначить ревьюверов 
- `POST /pullRequest/merge` - Пометить PR как MERGED 
- `POST /pullRequest/reassign` - Переназначить ревьювера 

### Health

- `GET /health` - Проверка здоровья сервиса


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
  -d '{
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search",
    "author_id": "u1"
  }'
```

### Получение PR пользователя

```bash
curl -X GET "http://localhost:8080/users/getReview?user_id=u2" \
```



## Реализованные функции

✅ Все основные эндпоинты согласно OpenAPI спецификации
✅ Автоматическое назначение до 2 ревьюверов при создании PR
✅ Переназначение ревьюверов
✅ Идемпотентная операция merge
✅ Управление активностью пользователей
✅ Защита от изменения ревьюверов после merge
✅ Чистая архитектура с разделением слоев
✅ Логи для http запросов
✅ Нагрузочное тестирование(результаты в файле report.html)
✅ e2e тестирование(make e2e)


## Принятые решения

1. **Миграции**: Миграции автоматически применяются при первом запуске PostgreSQL через Docker Compose (файлы из `migrations/` копируются в `/docker-entrypoint-initdb.d/`).

2. **Хранение ревьюверов**: Список ревьюверов хранится в виде JSONB массива в PostgreSQL для удобства работы с JSON операциями.

3. **Выбор ревьюверов**: Используется случайный выбор из доступных кандидатов с использованием `math/rand`.

4. **Обработка ошибок**: Все доменные ошибки оборачиваются в структурированный формат согласно OpenAPI спецификации.


