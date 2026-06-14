# AI Coding Agent Instructions

This project follows a strict Clean Architecture pattern using Go modules (`go.work`). It relies heavily on `task` (Taskfile) for building, testing, and formatting.

## Related Documentation
* [Homework Instructions (`HOW_TO_DO/hw.md`)](HOW_TO_DO/hw.md) - For context on the project's goal.
* [Taskfile (`Taskfile.yml`)](Taskfile.yml) - Definitions for all build, test, and generation commands.

## Architecture Guidelines
This project enforces a five-layer structure across microservices (e.g., `inventory`, `order`, `payment`):
1. **API Layer (`internal/api/[service]/v1/`)**: Exposes gRPC or REST. Handlers map requests to the service layer.
   * *Protocol*: gRPC uses `*proto.Request` and REST uses standard HTTP handlers.
2. **Service Layer (`internal/service/[domain]/`)**: Business logic layer.
3. **Repository Layer (`internal/repository/[domain]/`)**: In-memory stores configured with mutexes `sync.RWMutex`.
4. **Model Layer (`internal/model/`)**: Business domain types and predefined errors.
5. **Converter Layer (`internal/converter/`)**: Distinct mappers (e.g., model <-> proto, repo-model <-> service-model). 

*Strict interface boundaries* exist (`var _ Interface = (*Implementation)(nil)`) and dependency injection is enforced via constructors instead of global variables.

## Build and Tooling Commands
Run all tasks using the `task` CLI tool. It pulls localized binaries to `bin/` automatically.
* **Code Generation**: Run `task gen`. Uses `buf` for gRPC/protobufs and `ogen` for OpenAPIs.
* **Linting & Formatting**: Run `task format` (gofumpt & gci) then `task lint` (golangci-lint).
* **Mocks**: Run `task mockery:gen` to update `mockery` stubs.

## Testing Guidelines
* Run `task test-coverage` to run unit tests and check against the 40% coverage requirement.
* Run `task test-api` to execute broader e2e scenarios against the APIs via `grpcurl` and curls.
* **Libraries Pattern**: `github.com/stretchr/testify/suite` runs test suites.
* Every module requires its own mocks inside `internal/[layer]/mocks/`. Write testing suites in `suite_test.go` and logic assertions in `[action]_test.go`.

## Environment & Common Issues
* Never use `go get` from the root layer. Change directories into a specific module (`cd inventory`, `cd order`, etc.) before managing Go dependencies. Or prefer `task deps:update` from root.
* Check predefined errors in `internal/model/error.go` before defining locally string-based ones. Map validation to `codes.InvalidArgument` inside gRPC handlers, and use `400 Bad Request` mapping for HTTP.
