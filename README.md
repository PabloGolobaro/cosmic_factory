# Cosmic Factory

![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/PabloGolobaro/989a2757665014f43115badfbc137219/raw/coverage.json)
![CI](https://github.com/PabloGolobaro/cosmic_factory/actions/workflows/ci.yml/badge.svg)

Учебный проект в рамках курса по микросервисам на Go. Симулирует фабрику сборки космических кораблей: клиент выбирает детали, оформляет заказ, оплачивает — корабль уходит в сборку и доставляется.

## Архитектура

Go workspace-монорепо из семи модулей:

```
order/      REST API (:8080)  — оркестрирует весь флоу заказа
inventory/  gRPC (:50051)     — каталог деталей, резервирование и списание
payment/    gRPC (:50052)     — обработка платежей
iam/        gRPC (:50053)     — аутентификация, сессии (PostgreSQL + Redis)
assembly/   Kafka consumer    — сборка корабля по событию OrderPaid
platform/   общие утилиты    — Kafka, closer, middleware
shared/     proto + OpenAPI   — контракты между сервисами
```

**Поток запроса:**

```
HTTP Client → Order (ogen/chi) → Inventory (gRPC) — резервирование деталей
                               → Payment  (gRPC) — списание средств
                               → Kafka           — событие OrderPaid
                                    ↓
                              Assembly (consumer) → Kafka — ShipAssembled
                                    ↓
                              Order (consumer)    → Inventory.CommitParts
```

## Технологии

| Слой | Инструмент |
|---|---|
| HTTP codegen | [ogen](https://github.com/ogen-go/ogen) из OpenAPI 3.0 |
| gRPC | buf + protoc-gen-go |
| База данных | PostgreSQL · pgx/v5 · goose |
| Транзакции | avito-tech/go-transaction-manager |
| Кэш / сессии | Redis |
| Очередь | Kafka (Sarama) — Redpanda в тестах |
| Аутентификация | JWT-less сессии через IAM-сервис |
| Тесты | testcontainers, testify, mockery |
| Линтер | golangci-lint v2 |
| Task runner | [Task](https://taskfile.dev) |

## Быстрый старт

```bash
# 1. Инфраструктура
task deploy:all:up

# 2. Миграции
task migrate:order:up
task migrate:inventory:up
task migrate:iam:up

# 3. Сервисы
task run:order       # :8080
task run:inventory   # :50051
task run:payment     # :50052
```

Swagger UI: [http://localhost:8080/swagger-ui.html](http://localhost:8080/swagger-ui.html)

## Тесты

```bash
task test:unit      # unit с race-детектором
task test:api       # интеграционные (testcontainers)
task test:e2e       # сквозные (Postgres + Redpanda)
task test:coverage  # покрытие с порогом 40%
```

## Best practices курса

- Layered architecture: `api → service → repository`
- Dependency injection вручную (без фреймворков)
- Code generation: OpenAPI → handlers, proto → stubs, mockery → моки
- `SELECT FOR UPDATE` для конкурентных операций
- Транзакционный outbox: Pay + PublishOrderPaid в одной транзакции
- Структурированное логирование через `log/slog`
- CI: build → lint → unit → api → e2e → coverage (GitHub Actions)
