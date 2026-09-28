# AI Coding Agent Instructions

This project follows a strict Clean Architecture pattern using Go modules (`go.work`). It relies heavily on `task` (Taskfile) for building, testing, and formatting.

## Related Documentation
* [Homework Instructions (`HOW_TO_DO/hw.md`)](HOW_TO_DO/hw.md) - For context on the project's goal.
* [Taskfile (`Taskfile.yml`)](Taskfile.yml) - Definitions for all build, test, and generation commands.

## Architecture Guidelines
This project enforces a five-layer structure across microservices (`inventory`, `order`, `payment`):

### 1. API Layer (`internal/api/[service]/v1/`)
Exposes gRPC or REST. Handlers map requests to the service layer.
* gRPC uses `*proto.Request` and REST uses standard HTTP handlers.

### 2. Service Layer (`internal/service/[domain]/`)
Business logic layer. Implementations live in sub-packages (e.g., `service/part/`, `service/payment/`).

### 3. Repository Layer (`internal/repository/[domain]/`)
Data access. **Not all repositories are in-memory** — inventory uses MongoDB, order uses Postgres, only payment uses in-memory with `sync.RWMutex`.

### 4. Model Layer (`internal/model/`)
Business domain types and predefined errors in `error.go`. Map validation to `codes.InvalidArgument` (gRPC) or `400 Bad Request` (HTTP).

### 5. Converter Layer (`internal/converter/`)
Distinct mappers (model <-> proto, repo-model <-> service-model).

*Strict interface boundaries* exist (`var _ Interface = (*Implementation)(nil)`) and dependency injection is enforced via constructors instead of global variables.

---

## App Layer (`internal/app/`)

Every service has an `app` package as the single entry point for initialization and lifecycle:

* **`di.go`** — Dependency injection container. A `serviceProvider` struct with lazy-init getters:
  ```
  Config() → MongoClient() → Collection() → DBCollection() → PartRepository() → InventoryService() → API() → GRPCServer()
  ```
  Order of initialization matters — each getter calls the previous ones.
* **`app.go`** — Application lifecycle. `NewApp()` initializes all deps, `Run(ctx)` starts gRPC server, handles `SIGTERM`/`SIGINT`, and calls `Closer().CloseAll()` for graceful shutdown (`listener → gRPC → DB → logger`).
* **`cmd/main.go`** — Must be ~15-20 lines: load `.env` → `signal.NotifyContext` → `app.NewApp` → `app.Run`.

Health check is registered via `platform/pkg/grpc/health` inside the gRPC server init.

---

## Config Layer (`internal/config/`)

Follows 12-factor app principles — configuration from environment variables:

* **`interfaces.go`** — Defines `[Service]Config` interface (e.g., `InventoryConfig`) with sub-interfaces: `GRPCConfig`, `LoggerConfig`, `MongoConfig` (or `PostgresConfig` for order).
* **`env/`** — Env providers implementing the interfaces. Each reads from `os.Getenv` with defaults and validation.
* **`config.go`** — `New()` constructor that assembles all sub-configs.
* **`mocks/`** — Auto-generated mocks via `task mockery:gen`. Config mocks are used in service/repo tests to avoid real env vars.

---

## Platform Library (`/platform`)

Separate Go module shared by all services:

| Component | Path | Purpose |
|-----------|------|---------|
| Logger | `platform/pkg/logger` | Unified logging on top of `go.uber.org/zap` |
| Closer | `platform/pkg/closer` | Graceful resource cleanup (`AddNamed`, `CloseAll`) |
| Health | `platform/pkg/grpc/health` | gRPC Health Check Service (`SERVING`/`NOT_SERVING`) |
| Testcontainers | `platform/pkg/testcontainers` | Wrappers for MongoDB, networks, app containers |
| Migrator | `platform/pkg/migrator` | Postgres migrations (used by `order`) |

---

## Build and Tooling Commands
Run all tasks using the `task` CLI tool. It pulls localized binaries to `bin/` automatically.

| Command | What it does |
|---------|-------------|
| `task gen` | Code generation: `buf` for gRPC/protobufs, `ogen` for OpenAPIs |
| `task format` | Formatting: `gofumpt` + `gci` |
| `task lint` | Linting: `golangci-lint` |
| `task mockery:gen` | Update mockery stubs for all modules |
| `task test-coverage` | Unit tests + 40% coverage check |
| `task test-api` | E2E smoke tests via `grpcurl` and `curl` |
| `task test-integration` | Integration tests with Testcontainers (build tag `integration`) |
| `task env:generate` | Generate `.env` files from `deploy/env/*.env.template` |

**Important:** Never use `go get` from root. Use `cd <module> && go get ...` or `task deps:update`.

---

## Testing Guidelines

### Unit Tests
* Use `github.com/stretchr/testify/suite` for test suites.
* Every module has mocks in `internal/[layer]/mocks/`.
* Suite setup: `suite_test.go`, assertions: `[action]_test.go`.
* Config env tests use `t.Setenv` for isolation.

### Integration Tests (`tests/integration/`)
* Build tag: `//go:build integration`
* **Inventory**: e2e test against real MongoDB via testcontainers (API as black box).
* **Order**: repo-level test against real Postgres via testcontainers.
* Structure:
  ```
  tests/integration/
  ├── constants.go
  ├── setup.go
  ├── teardown.go
  ├── test_environment.go
  ├── suite_test.go
  └── inventory_test.go (or order_test.go)
  ```
* Run: `task test-integration`

---

## Environment Management

### Templates (`deploy/env/`)
* `.env.template` — master template with defaults for all services
* `[service].env.template` — per-service templates (payment, inventory, order)
* `generate-env.sh` — generates `[service]/.env` files from templates

### Docker Compose (`deploy/compose/`)
* Each service has its own `docker-compose.yml` with `env_file: .env`
* Use `task up-inventory`, `task down-inventory` etc. for container management

### Workflow
1. Edit `deploy/env/.env.template` or `[service].env.template`
2. Run `task env:generate` → creates `deploy/compose/[service]/.env`
3. Service reads env via `internal/config/env/` providers

---

## Environment & Common Issues
* Never use `go get` from the root layer. Change directories into a specific module (`cd inventory`, `cd order`, etc.) before managing Go dependencies. Or prefer `task deps:update` from root.
* Check predefined errors in `internal/model/error.go` before defining locally string-based ones.
* After changing proto files → `task gen`. After changing interfaces → `task mockery:gen`.
* After all code changes → `task format && task lint`.
