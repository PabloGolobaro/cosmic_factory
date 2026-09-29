# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

All top-level commands use [Task](https://taskfile.dev) (`task`). Run from the project root.

```bash
task setup            # install all dev tools (golangci-lint, gofumpt, gci, buf, ogen, mockery, goose)
task format           # gofumpt + gci import sorting
task lint             # golangci-lint (order/inventory/payment/shared only — lint iam/platform explicitly)
task gen:all          # regenerate all code (proto → Go, OpenAPI → Go)
task test:unit        # unit tests with race detector (all modules)
task test:api         # API tests (order, iam, inventory) via testcontainers — Docker required
task test:e2e         # order e2e with real Kafka (Redpanda)
task run:all          # run services locally (or run:<svc>, cwd = <svc>/)
task test:coverage    # coverage with 40% minimum threshold
task test:coverage:html  # generate HTML coverage report
task test:mocks:gen   # regenerate mocks
task deps:update      # go work sync + go mod tidy for all modules
```

**Run a single test:**
```bash
cd order && go test -v -run TestName ./internal/service/order/...
```

**Single API test:** `go test -tags=apitest -run TestName ./<svc>/tests/...`

**Database migrations:**
```bash
task migrate:order:up
task migrate:inventory:up
task migrate:iam:up     # files in migrations/<svc>/ (goose)
```

**Infrastructure (Docker Compose):**
```bash
task deploy:all:up    # core (Kafka) + order (3 instances) + inventory + iam
task deploy:all:down  # tear down everything
task deploy:core:up   # Kafka only
task deploy:order:up  # order service
```

`.env` files (`order.env`, `inventory.env`, `iam.env`, `core.env`) are at the project root and contain `DB_URI`, `POSTGRES_*`, `MIGRATIONS_DIR`.

## Architecture

### Monorepo layout

Go workspace (`go.work`, gitignored) with modules: `order`, `inventory`, `payment`, `iam`, `assembly`, `platform`, `shared`.

```
order/      REST API service (:8080) — central service, depends on Inventory & Payment
inventory/  gRPC service (:50051) — part catalog with PostgreSQL
payment/    gRPC service (:50052) — stateless payment processing
iam/        gRPC service (:50053) — users (PostgreSQL), sessions (Redis), authorization (embedded OPA)
assembly/   Kafka consumer — ship assembly
platform/   shared infra libs: auth context, authz actions, closer, logger, tracing, kafka, ratelimit
shared/     proto definitions, generated stubs, OpenAPI spec, common gRPC interceptors
```

### Request flow

External HTTP → **Order** (ogen-generated handlers + chi router)  
Order → **Inventory** (gRPC, `localhost:50051`) — check/reserve parts  
Order → **Payment** (gRPC, `localhost:50052`) — process payment  
Order/Inventory → **IAM** (gRPC, `localhost:50053`) — session check (`Whoami`) and authorization (`Authorize`)  
Order ⇄ **Assembly** via Kafka — order paid → ship assembled  

### Auth & authorization

- HTTP: `Authorization: Bearer <session_uuid>` → Order middleware calls `IAM.Whoami`; gRPC: metadata `session-uuid` (forwarded by `SessionForwarder`)
- AuthZ: IAM is the PDP — `AuthService.Authorize(session, action, owner_uuid)`; policy in `iam/internal/authz/policy/authz.rego`; action names in `platform/pkg/authz`
- Order enforces in the service layer (owner known only after loading the order); inter-service calls are not authorized
- Roles: `client` (own orders), `manager` (read/cancel any). Seeded users `testuser` / `testmanager`, password `password123`

### Internal layer structure (per service)

```
cmd/main.go              entry point: loads .env + YAML config, starts internal/app
internal/app/di.go       lazy DI container (composition root); pkg/app — wiring reused by API tests
internal/api/v1/         HTTP/gRPC handlers
internal/service/        business logic (interface-based)
internal/repository/     PostgreSQL via pgx/v5 + pgxpool
internal/client/grpc/    outbound gRPC clients
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

- **Proto → Go**: edit `shared/proto/`, run `task gen:all` → updates `shared/pkg/proto/`
- **OpenAPI → Go**: edit `shared/api/order/v1/*.yaml`, run `task gen:all` → updates `shared/pkg/openapi/order/v1/`
- **Mocks**: edit interfaces, run `task test:mocks:gen` → `mockery` reads `.mockery.yaml`

Never edit generated files under `pkg/proto/`, `pkg/openapi/`, or `mocks/` directly.

### Linting constraints

- Forbidden imports: `io/ioutil`, `log` (use `slog`), `satori/uuid` (use `google/uuid`), `math/rand` (use v2), `github.com/pkg/errors`, old `google.golang.org/protobuf` import paths.
- Max function length: 100 lines; max cyclomatic complexity: 20.
- `perfsprint` enforcer: use `strconv` / direct string ops instead of `fmt.Sprintf` where no format verb is needed.

### Environment

Config: YAML `<svc>/config.<env>.yaml` (`-config` flag > `CONFIG_PATH` > `config.local.yaml`), env vars override. `.env` is loaded via `godotenv.Load("./../<svc>.env")` — run binaries from `<svc>/` (as `task run:<svc>` does).

### Gotchas

- Local Go is 1.27, but mockery v3.7 and golangci-lint can't read its export data (`export data version 4 ...`): run them with `GOTOOLCHAIN=go1.26.0` (e.g. `GOTOOLCHAIN=go1.26.0 bin/mockery`)
- `task format` runs `go fix`, which may touch unrelated files — review the diff; `task gen:all` may regenerate `shared/pkg/proto/buf/validate/validate.pb.go` without cause — revert it
- Adding a column: follow expand → backfill → contract as separate migrations (see `migrations/order/20260517*`)
- Rego tests run from Go (`iam/internal/authz/authz_test.go` via `opa/v1/tester`), no opa CLI needed
