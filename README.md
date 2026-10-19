# Микросервисы на Go: Inventory · Order · Payment

Pet-проект: домен интернет-магазина автозапчастей, разбитый на три независимых сервиса.
Каждый сервис собран по Clean Architecture, общается с соседями по gRPC, а наружу Order
отдаёт REST по OpenAPI-контракту, сгенерированному из спецификации.

```
Клиент ──HTTP──▶ Order ──gRPC──▶ Inventory (склад запчастей, MongoDB)
                   │
                   └──gRPC──▶ Payment (проведение оплаты)
                                PostgreSQL ◀── Order (хранение заказов, миграции)
```

## Стек

| Область | Технологии |
|---|---|
| Язык и сборка | Go 1.25, Go workspaces (`go.work`), Taskfile |
| Контракты | Protocol Buffers + gRPC (`buf`), OpenAPI + `ogen` |
| Хранилища | PostgreSQL 17 (+ миграции `goose`), MongoDB 7 |
| Инфраструктура | Docker, Docker Compose, GitHub Actions |
| Качество | `testify` (suite), `testcontainers-go`, mockery, golangci-lint, gofumpt + gci |

## Архитектура

Каждый сервис разделён на пять слоёв, зависимости идут строго внутрь:

| Слой | Ответственность |
|---|---|
| `internal/api/<service>/v1/` | Транспорт: gRPC-хендлеры или HTTP-хендлеры, маппинг запросов в сервисный слой |
| `internal/service/<domain>/` | Бизнес-логика, оркестрация клиентов и репозиториев |
| `internal/repository/<domain>/` | Доступ к данным: MongoDB (inventory), PostgreSQL (order) |
| `internal/model/` | Доменные типы и предопределённые ошибки |
| `internal/converter/` | Маппинг model ↔ proto и repo-model ↔ service-model |

Ключевые решения:

- **Зависимости только через интерфейсы** — в каждом слое есть проверка
  `var _ Interface = (*Implementation)(nil)`, глобальных переменных нет.
- **DI-контейнер** — `internal/app/di.go`: `serviceProvider` с ленивыми геттерами
  (`Config → MongoClient/DB → Repository → Service → API → GRPCServer`).
- **Жизненный цикл** — `internal/app/app.go`: старт сервера, обработка
  `SIGTERM`/`SIGINT` и graceful shutdown через `platform/pkg/closer`
  (listener → gRPC → БД → логгер). `cmd/main.go` у каждого сервиса — около 30 строк.
- **Health-check** — gRPC Health Checking Protocol, регистрируется в общем
  пакете `platform/pkg/grpc/health`.

### Роль сервисов

| Сервис | Транспорт | Хранилище | Что делает |
|---|---|---|---|
| `inventory` | gRPC `:50051` | MongoDB | Каталог запчастей: `GetPart`, `ListParts` с фильтром по UUID/названию/категории/стране/тегам. При старте засеивает коллекцию демо-деталями, если она пуста |
| `order` | HTTP `:8080` | PostgreSQL | Заказы: создание, получение, оплата, отмена. Считает стоимость через `inventory`, проводит платёж через `payment`, хранит статусы `PENDING_PAYMENT → PAID / CANCELLED`. Миграции применяются при старте |
| `payment` | gRPC `:50052` | — | Проведение оплаты, возвращает `transaction_uuid`. Пока заглушка без персистентности (см. Roadmap) |

## Структура репозитория

```
.
├── inventory/ | order/ | payment/   # сервисы: cmd/, internal/{api,service,repository,model,converter}, tests/integration
├── platform/                        # общий модуль: logger, closer, health, migrator, testcontainers
├── shared/                          # proto-контракты, OpenAPI-спека, сгенерированный код
│   ├── proto/{inventory,payment}/v1/*.proto
│   ├── api/order/v1/order.openapi.yaml
│   └── pkg/{proto,openapi}/         # результат task gen
├── deploy/
│   ├── compose/{core,inventory,order,payment}/docker-compose.yml
│   └── env/*.env.template           # конфиги + generate-env.sh
├── Taskfile.yml                     # все команды проекта
└── AGENTS.md                        # правила и стек для AI-ассистентов
```

## Быстрый старт

Нужны Docker, Go 1.25+ и [Task](https://taskfile.dev). Бинарники инструментов
(`buf`, `ogen`, `mockery`, линтеры) Task скачает сам в `bin/`.

```bash
# 1. Сгенерировать .env'ы: deploy/env/.env.template -> deploy/env/.env -> deploy/compose/<service>/.env
task env:generate

# 2. Поднять всё: общая сеть + inventory+mongo, payment, order+postgres
task docker:compose:up

# 3. Остановить и удалить контейнеры вместе с томами
task docker:compose:down
```

Отдельно по сервисам: `task docker:compose:inventory:up`, `task docker:compose:payment:up`,
`task docker:compose:order:up` (сервисы зависят от сети, поэтому сначала `task docker:compose:core:up`).

Запуск без Docker (нужны локальные MongoDB и PostgreSQL). `.env` читается из текущей
директории, поэтому команды запускаем из корня репозитория:

```bash
cp .env.example .env       # указать свои MONGO_URI / DB_URI
go run ./inventory/cmd     # или ./order/cmd, ./payment/cmd
```

Для заказа нужны запущенные inventory и payment: order ходит в них по gRPC.

### Порты

| Что | Наружу |
|---|---|
| inventory gRPC | `50051` |
| payment gRPC | `50052` |
| order HTTP | `8080` |
| MongoDB | `27017` |
| PostgreSQL | `5433` (внутри docker-сети — `5432`) |

### Проверка API

```bash
# gRPC inventory (reflection включён)
task grpcurl:install
./bin/grpcurl -plaintext -d '{"filter":{}}' localhost:50051 inventory.v1.InventoryService/ListParts

# REST order
curl -X POST http://localhost:8080/api/v1/orders \
  -H 'Content-Type: application/json' \
  -d '{"user_uuid":"<uuid>","part_uuids":["<part-uuid>"]}'

curl http://localhost:8080/api/v1/orders/<order-uuid>
curl -X POST http://localhost:8080/api/v1/orders/<order-uuid>/pay
curl -X POST http://localhost:8080/api/v1/orders/<order-uuid>/cancel
```

Полный сценарий «получить деталь → создать заказ → оплатить → отменить» прогоняется одной
командой `task test-api` (grpcurl + curl) по уже запущенным сервисам.

## Тесты и качество

```bash
task test             # юнит-тесты всех модулей (testify/suite, моки из internal/*/mocks)
task test-coverage    # покрытие бизнес-логики (только service/ и repository/ пакеты) + сводный отчёт
task test-integration # интеграционные тесты на реальных БД через testcontainers (тег integration)
task test-api         # E2E-проверка API через grpcurl и curl
task format && task lint
```

Текущее покрытие бизнес-логики: **81.9%** суммарно (inventory 68.1% · order 89.5% · payment 100%).
Покрытие считается только по пакетам `internal/service/<domain>` и `internal/repository/<domain>` —
mocks, converters и model из расчёта исключены, чтобы цифра отражала логику, а не генерат.

Интеграционные тесты работают как чёрный ящик: поднимают MongoDB/PostgreSQL и само
приложение в контейнерах (`platform/pkg/testcontainers`), затем дёргают реальный API.

## Генерация кода

| Команда | Что делает |
|---|---|
| `task gen` | `buf` → gRPC-код из `shared/proto`, `ogen` → HTTP-сервер из OpenAPI |
| `task mockery:gen` | Обновляет моки по `.mockery.yaml` |
| `task deps:update` | Обновляет зависимости всех модулей workspace |

Изменил `.proto` → `task gen`; изменил интерфейс → `task mockery:gen`.

## CI

`.github/workflows/ci.yml`: `extract-vars` (версии инструментов берутся из `Taskfile.yml`,
чтобы не было расхождений) → `test` → `integration-test`. Рядом лежат reusable-workflow
для линта, юнит- и интеграционных тестов и покрытия; в `ci.yml` сейчас включены тесты
и интеграция.

## Разработка с AI-ассистентами

Код я пишу сам, ассистента подключаю по ходу. Чаще всего это разбор незнакомой библиотеки,
взгляд на решение со стороны («что я тут проглядел?»), поиск краевых случаев и механические
правки по репозиторию. Черновики тестов и моки иногда тоже отдаю ему — так быстрее, но всё
равно перечитываю.

Главный подводный камень: ассистент уверенно выдаёт неправильный код. Поэтому проверяю не
глазами, а инструментами — `task test`, `task test-integration` на Testcontainers, `task lint`,
а после изменений в контрактах ещё `task gen` и `task mockery:gen`. Пока тесты и линтер не
зелёные, код не готов, кто бы его ни написал.

Чтобы не пересказывать ассистенту архитектуру в каждой сессии, завёл `AGENTS.md`: пять слоёв
и их границы, без глобальных переменных, DI через конструкторы, где лежат предопределённые
ошибки. Файл читают любые ассистенты, так что я не привязан к одному инструменту.

## Roadmap

- `payment`: хранить транзакции в БД и сделать повторные оплаты идемпотентными.
- Отдельный сервис аутентификации + передача `user_id` через контекст (логгер уже
  поддерживает enrich по `trace_id` / `user_id`).
- Redis для кеша каталога запчастей.
- Kafka для событий заказа: создан, оплачен, отменён.
- Notification Service — уведомления о смене статуса заказа.
- Prometheus + OpenTelemetry: метрики и трейсинг, дашборд в Grafana.
- Контрактные тесты между сервисами в CI.
