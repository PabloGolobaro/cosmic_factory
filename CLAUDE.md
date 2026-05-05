# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

All top-level commands use [Task](https://taskfile.dev) (`task`). Run from the project root.

```bash
task setup            # install all dev tools (golangci-lint, gofumpt, gci, buf, ogen, mockery, goose)
task format           # gofumpt + gci import sorting
task lint             # golangci-lint across all modules
task gen              # regenerate all code (proto → Go, OpenAPI → Go, mocks)
task test             # unit tests with race detector (all modules)
task test:coverage    # coverage with 40% minimum threshold
task deps:update      # go work sync + go mod tidy for all modules
```

**Run a single test:**
```bash
cd order && go test -v -run TestName ./internal/service/order/...
```

**Database migrations:**
```bash
task migrate:order:up
task migrate:inventory:up
```

**Infrastructure (Docker Compose):**
```bash
task up-all           # network + inventory + order
task down-all         # tear down everything
```

`.env` files (`order.env`, `inventory.env`) are at the project root and contain `DB_URI`, `POSTGRES_*`, `MIGRATIONS_DIR`.

## Architecture

### Monorepo layout

Go workspace (`go.work`) with four modules: `order`, `inventory`, `payment`, `shared`.

```
order/      REST API service (:8080) — central service, depends on Inventory & Payment
inventory/  gRPC service (:50051) — part catalog with PostgreSQL
payment/    gRPC service (:50052) — stateless payment processing
shared/     proto definitions, generated stubs, OpenAPI spec, common gRPC interceptors
```

### Request flow

External HTTP → **Order** (ogen-generated handlers + chi router)  
Order → **Inventory** (gRPC, `localhost:50051`) — check/reserve parts  
Order → **Payment** (gRPC, `localhost:50052`) — process payment  

### Internal layer structure (per service)

```
cmd/main.go              wires everything together (DI by hand)
internal/api/v1/         HTTP/gRPC handlers
internal/service/        business logic (interface-based)
internal/repository/     PostgreSQL via pgx/v5 + pgxpool
internal/client/grpc/    outbound gRPC clients (order only)
internal/model/          domain types
internal/converter/      domain ↔ transport DTO conversion
```

### Key libraries & conventions

| Concern | Library |
|---|---|
| HTTP routing | `go-chi/chi/v5` |
| HTTP codegen | `ogen-go/ogen` from OpenAPI 3.0 YAML |
| gRPC | `google.golang.org/grpc` |
| Proto codegen | `buf` + `protoc-gen-go` / `protoc-gen-go-grpc` |
| Proto validation | `buf.build/go/protovalidate` (gRPC interceptor) |
| DB driver | `jackc/pgx/v5` + `pgxpool` |
| Transactions | `avito-tech/go-transaction-manager` |
| Migrations | `pressly/goose/v3` |
| Mocks | `vektra/mockery/v3` (testify+Expecter, generated into `mocks/` subdir) |
| Integration tests | `testcontainers/testcontainers-go` |
| Logging | stdlib `log/slog` |

### Code generation

- **Proto → Go**: edit `shared/proto/`, run `task gen` → updates `shared/pkg/proto/`
- **OpenAPI → Go**: edit `shared/api/order/v1/*.yaml`, run `task gen` → updates `shared/pkg/openapi/order/v1/`
- **Mocks**: edit interfaces, run `task gen` → `mockery` reads `.mockery.yaml`

Never edit generated files under `pkg/proto/`, `pkg/openapi/`, or `mocks/` directly.

### Linting constraints

- Forbidden imports: `io/ioutil`, `log` (use `slog`), `satori/uuid` (use `google/uuid`), `math/rand` (use v2), `github.com/pkg/errors`, old `google.golang.org/protobuf` import paths.
- Max function length: 100 lines; max cyclomatic complexity: 20.
- `perfsprint` enforcer: use `strconv` / direct string ops instead of `fmt.Sprintf` where no format verb is needed.

### Environment

Services load their `.env` via `godotenv.Load("./../../<service>.env")` — path is relative to `<service>/cmd/` (the working directory when running the binary). Both `.env` files live at the project root.
