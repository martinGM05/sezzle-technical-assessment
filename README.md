# Sezzle Technical Assessment

A full-stack calculator being built with React, TypeScript, and Go.

## Prerequisites

- Node.js 22.18.0 (see `.nvmrc`)
- pnpm 12.9.1
- Go 1.24.2

## Setup

```bash
cd frontend
pnpm install --frozen-lockfile
```

## Frontend commands

```bash
pnpm dev
pnpm build
pnpm lint
pnpm format:check
pnpm typecheck
pnpm test
pnpm test:coverage
```

## Backend setup and execution

The backend is an independent Go module with standard-library dependencies only.
From the repository root:

```bash
cd backend
go build -o bin/api ./cmd/api
./bin/api
```

The server listens on port 8080. `Ctrl+C` (SIGINT) or SIGTERM initiates a graceful
shutdown with a five-second deadline. Expiry closes active connections and exits
with status 1; a clean shutdown exits with status 0. Startup failures also exit 1.

For the planned same-origin reverse proxy deployment, configure the public origin:

```bash
ALLOWED_ORIGINS=https://calculator.example ./bin/api
```

The proxy may terminate HTTPS and forward HTTP to the backend. Explicit origins
are required even when browser UI and API share one public origin: browsers can
send Origin on same-origin POSTs. Comma-separated origins are supported; wildcard
and invalid origins are not allowed. Empty configuration denies nonempty Origin
headers. Forwarded headers never grant access. See [architecture](docs/ARCHITECTURE.md).

## API examples

```bash
curl -i http://localhost:8080/health
curl -I http://localhost:8080/health
curl -i http://localhost:8080/api/v1/calculations \
  -H 'Content-Type: application/json' \
  -d '{"operation":"divide","operands":[10,2]}'
# 200: {"result":5}

curl -i http://localhost:8080/api/v1/calculations \
  -H 'Content-Type: application/json' \
  -d '{"operation":"divide","operands":[10,0]}'
# 422, application/problem+json:
# {"type":"about:blank","title":"Unprocessable Content","status":422,
#  "detail":"Cannot divide by zero.","code":"DIVISION_BY_ZERO"}
```

Supported operation names are case-sensitive:

- `add`, `subtract`, `multiply`, `divide`, `power`: exactly two operands.
- `sqrt`: exactly one non-negative operand.
- `percentage`: exactly one operand, divided by 100.

Examples: power `[2,3]` → 8; sqrt `[9]` → 3; percentage `[25]` → 0.25.
Finite 0^0 is 1. Results use float64 without display rounding; IEEE 754 precision
limitations apply. This is not a monetary calculation engine.

[OpenAPI](docs/openapi.yaml) defines request/response schemas, arities, stable
error codes and status codes: 400, 403, 404, 405, 413, 415, 422 and 500. Request
bodies are limited to one MiB. HEAD never returns a body, including for errors.

## Expression API

```bash
curl -i http://localhost:8080/api/v1/expressions \
  -H 'Content-Type: application/json' \
  -d '{"expression":"200 * 15% + sqrt(9 + 7)"}'
# 200: {"result":34}
```

Supports `+`, `-`, `*`, `/`, `^`, `sqrt(...)`, postfix `%`, unary signs and nested
parentheses. Power is right-associative: `2^3^2 = 512`; `-2^2 = -4` and `2^-3 =
0.125`. Percent always divides by 100, not contextual financial percentage.
Expressions are limited to 4096 ASCII bytes and 64 nested recursive constructs.
See [grammar, errors and mathematical assumptions](docs/EXPRESSIONS.md).
The frontend is not connected to this endpoint in Phase 4.

## Backend quality checks

From `backend/`:

```bash
test -z "$(gofmt -l cmd internal)"
go vet ./...
go test -count=1 -race ./...
go test -count=1 ./internal/transport/http -run 'Test(HTTPOpenAPI|OpenAPIHead|ContractChecker)'
go build ./...
```

The full suite includes subprocess tests for real SIGINT/SIGTERM and startup
failures. `go test -short ./...` skips those subprocess tests and is not the full
quality gate. Coverage commands and artifacts are described below.

## Backend coverage and CI

Run the same coverage commands as CI from `backend/`:

```bash
mkdir -p coverage
go test -count=1 -race -covermode=atomic -coverprofile=coverage/coverage.out ./...
go tool cover -func=coverage/coverage.out | tee coverage/summary.txt
go tool cover -html=coverage/coverage.out -o coverage/index.html
```

Generated reports are ignored by Git. CI publishes the measured function/total
coverage in its job summary and uploads profile, text and HTML reports as
`backend-coverage`. Backend and frontend jobs run independently. External actions
are SHA-pinned; the backend Go version comes from `backend/go.mod`. Full CI must
still be observed on GitHub after an authorized push; local checks are not a
remote CI result.

## Container

From the repository root:

```bash
docker build -t sezzle-calculator-api:local backend
docker run --rm --name sezzle-calculator-api \
  -p 127.0.0.1:8080:8080 \
  -e ALLOWED_ORIGINS=https://calculator.example \
  sezzle-calculator-api:local
# In another terminal:
curl -i http://localhost:8080/health
docker stop --time 10 sezzle-calculator-api
```

The multi-stage image contains a static Go binary and a digest-pinned distroless
non-root runtime. It has no shell or curl; check health externally through
`GET /health`, not with `docker exec sh`. Allow more than five seconds before an
external supervisor sends SIGKILL. Docker Compose is not required for the backend.

## Current scope

Phase 4 adds the bounded Go expression engine and its HTTP endpoint.

Phase 2 provides the Go calculator API, contract, behavioral tests, lifecycle,
backend container and quality pipeline. The frontend now includes the design
system, controlled calculator UI and collapsible history UI. Calculator controls
edit text only; equals reports that backend integration is pending. History starts
empty and no sample entries are injected. History components expose reuse, copy,
delete and confirmed-clear callbacks; copy is disabled without a supplied handler.
Dates are displayed in the user's locale and time zone. History shows 20 entries
at a time, with a button to reveal more.

Frontend evaluation, parsing, API/clipboard integration, history persistence,
full-stack end-to-end tests and performance work remain outside the delivered UI
scope. No database, broker, authentication or distributed infrastructure is needed.

See [requirements](docs/REQUIREMENTS.md), [architecture](docs/ARCHITECTURE.md) and
[non-blocking review findings](docs/REVIEW.md). Assessment requirements remain
reported rather than independently verified against the original assignment.

## Delivery verification

On 2026-11-26, the final backend suite passed **394 cases** with no failures or
skips. Fresh coverage was **98.4% overall** (executable 90.3%, calculator and HTTP
transport 100%). Race detection, vet, build, formatting and contract checks passed.
The command main wrapper remains outside direct test coverage; no artificial
coverage target was used.

The linux/arm64 Docker image passed actual non-root UID 65532, health, all seven
operations, CORS and SIGTERM checks. Active requests drained successfully; a
stalled request was forcibly closed after the shutdown deadline with exit 1.
Verification did not include image-signature verification, vulnerability scanning,
other runtime architectures or a remote GitHub Actions run.
