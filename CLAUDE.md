# Microservices Project - Clean Architecture

Проект микросервисной архитектуры на Go с использованием Clean Architecture и gRPC.

## Структура проекта

Монорепозиторий с Go workspace (`go.work`), содержащий:
- `inventory` — сервис управления запчастями
- `order` — сервис заказов
- `payment` — сервис платежей
- `shared` — общие proto-файлы и генерируемый код

## Архитектура

Все сервисы следуют строгой пятислойной Clean Architecture:

### 1. API Layer (`internal/api/[service]/v1/`)
- Точка входа: gRPC или REST
- Handlers маппят запросы в service layer
- gRPC использует `*proto.Request`, REST — стандартные HTTP handlers

### 2. Service Layer (`internal/service/[domain]/`)
- Бизнес-логика
- Интерфейсы с проверкой `var _ Interface = (*Implementation)(nil)`
- Dependency injection через конструкторы

### 3. Repository Layer (`internal/repository/[domain]/`)
- In-memory хранилища с `sync.RWMutex`
- Изолированные интерфейсы

### 4. Model Layer (`internal/model/`)
- Доменные типы
- Предопределённые ошибки в `error.go`

### 5. Converter Layer (`internal/converter/`)
- Маппинг между слоями:
  - model ↔ proto
  - repo-model ↔ service-model

**Важно:**
- Строгие границы интерфейсов
- Никаких глобальных переменных
- Только dependency injection

## Инструменты и команды

Все команды выполняются через `task` (Taskfile). Бинарники автоматически скачиваются в `bin/`.

### Генерация кода
```bash
task gen  # buf для gRPC/proto + ogen для OpenAPI
```

### Форматирование и линтинг
```bash
task format  # gofumpt + gci
task lint    # golangci-lint
```

### Моки
```bash
task mockery:gen  # Обновить mockery-стабы
```

### Тестирование
```bash
task test-coverage  # Unit-тесты + проверка покрытия (минимум 40%)
task test-api       # E2E тесты через grpcurl и curl
```

## Тестирование

- Используется `github.com/stretchr/testify/suite`
- Моки в `internal/[layer]/mocks/`
- Структура тестов:
  - `suite_test.go` — setup test suite
  - `[action]_test.go` — конкретные тесты

**Требования:**
- Минимум 40% покрытия кода
- Unit-тесты для каждого слоя
- E2E тесты для API

## Управление зависимостями

**ВАЖНО:** Никогда не запускай `go get` из корня проекта!

Варианты:
1. `cd inventory && go get ...` (или `order`, `payment`)
2. `task deps:update` из корня

## Обработка ошибок

- Используй предопределённые ошибки из `internal/model/error.go`
- Не создавай локальные string-based ошибки
- Маппинг для gRPC:
  - Валидация → `codes.InvalidArgument`
- Маппинг для HTTP:
  - Валидация → `400 Bad Request`

## Workflow

1. Изменил код → `task format`
2. Проверь линтером → `task lint`
3. Запусти тесты → `task test-coverage`
4. Если менял proto → `task gen`
5. Если менял интерфейсы → `task mockery:gen`

## Связанная документация

- [Homework Instructions](HOW_TO_DO/hw.md) — контекст и цели проекта
- [Taskfile.yml](Taskfile.yml) — все доступные команды

## Принципы разработки

- **Чистая архитектура** — строгое разделение слоёв
- **Интерфейсы везде** — зависимости только через интерфейсы
- **Тестируемость** — каждый слой покрыт тестами
- **Генерация кода** — proto и OpenAPI генерируются автоматически
- **Никаких глобалок** — только DI через конструкторы
