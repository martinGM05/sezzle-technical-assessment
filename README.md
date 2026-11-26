# Sezzle Calculator

React/TypeScript UI and a stateless Go REST API. Seven arithmetic operations,
compound expressions, keyboard input, light/dark/system themes and local history.
All mathematical parsing and evaluation happens in Go.

## Quick Start

**Docker is the recommended setup.** Install Docker Engine/Desktop with Docker
Compose v2; no local Node.js or Go installation is needed. Repository access is
required because the repository is private.

```bash
git clone --branch chore/final-delivery https://github.com/martinGM05/sezzle-technical-assessment.git
cd sezzle-technical-assessment
docker compose up --build
```

Open **http://localhost:8088** (or **http://127.0.0.1:8088**).
To shut down, press **Ctrl+C**, then remove the containers and networks:

```bash
docker compose down
```

### Docker health checks and configuration

For detached startup and health verification:

```bash
docker compose up --build -d --wait
docker compose ps
curl -fsS http://localhost:8088/healthz # Nginx
curl -fsS http://localhost:8088/health  # Go through the proxy
```

Only Nginx is published, on loopback; Go uses an internal network. Both runtimes
are non-root and read-only. The frontend health check covers Nginx and Go. After
an outage, `--wait` can fail while a previously unhealthy container awaits its
next probe; check `docker compose ps` and the health URLs for recovery.

For another port, change the exact allowed origins too:

```bash
APP_PORT=8090 ALLOWED_ORIGINS=http://localhost:8090,http://127.0.0.1:8090 \
  docker compose up --build -d --wait
```

Nginx preserves Origin and the API enforces the allowlist even for same-origin
POSTs. CORS is not authentication. Public deployment requires explicit TLS,
bind-address and origin configuration; this setup is local by default.

## Local development without Docker

Pinned tools: Node **22.18.0** (`.nvmrc`), pnpm **12.9.1** (`packageManager`),
Go **1.24.2** (`backend/go.mod`). Select Node with `nvm use` or your version manager.

```bash
npm install --global pnpm@12.9.1
cd frontend
pnpm install --frozen-lockfile
```

Start these in separate terminals, each from the repository root:

```bash
cd backend
ALLOWED_ORIGINS=http://127.0.0.1:5173 go run ./cmd/api
```

```bash
cd frontend
pnpm dev --host 127.0.0.1 --port 5173 --strictPort
```

Open http://127.0.0.1:5173. Vite proxies `/api` to Go on port 8080;
`API_PROXY_TARGET` overrides that target. Update `ALLOWED_ORIGINS` when changing
the frontend host/port. The browser has no hardcoded production API URL.
`pnpm preview` does not provide the development API proxy.

## Evaluate the application

- Try `(2+3)*4` → **20**, `2^3^2` → **512**, `sqrt(9+7)` → **4**,
  `200*15%` → **30**, and `10/0` → a recoverable error.
- Keyboard: digits, decimal point, parentheses and operators; `r` inserts `√(`.
  Enter/`=` evaluates, Backspace edits, Escape/Delete clears. Enter evaluates with
  any calculator button focused; Space activates it. Other controls remain native.
- Clear cancels pending requests; duplicate submissions and stale responses are
  guarded. Further input edits the expression; Clear starts over.
- History saves successful results only. Reuse, copy, delete and confirmed clear
  are available. Theme/history survive reload in browser localStorage, not Docker
  volumes. History retains 100 entries and displays 20 at a time.

Numbers use IEEE 754 `float64`, without decimal display rounding: `0.1+0.2` may
show `0.30000000000000004`. Percentage means division by 100; finite `0^0 = 1`.
Power is right-associative; `-2^2 = -4`, `(-2)^2 = 4`, `2^-3 = 0.125`.
Overflow/non-finite results are rejected; underflow follows IEEE 754. This is not
for monetary transactions: use decimal arithmetic or integer minor units for
those. [Grammar and limits](docs/EXPRESSIONS.md) document the 4096-byte ASCII and
64-level recursion bounds.

## API examples

Use port 8088 through Compose, or 8080 when running Go directly:

```bash
curl -i http://localhost:8088/api/v1/expressions \
  -H 'Content-Type: application/json' -d '{"expression":"(2+3)*4"}'
# 200: {"result":20}

curl -i http://localhost:8088/api/v1/calculations \
  -H 'Content-Type: application/json' -d '{"operation":"divide","operands":[10,2]}'
# 200: {"result":5}

curl -i http://localhost:8088/api/v1/expressions \
  -H 'Content-Type: application/json' -d '{"expression":"10/0"}'
# 422, application/problem+json:
# {"type":"about:blank","title":"Unprocessable Content","status":422,
#  "detail":"Cannot divide by zero.","code":"DIVISION_BY_ZERO"}
```

`add`, `subtract`, `multiply`, `divide`, `power` take two operands; `sqrt` and
`percentage` take one. [OpenAPI](docs/openapi.yaml) defines validation, status
codes and RFC 9457 errors. Go owns the one-MiB body limit. Proxy outages may return
Nginx 502/504 responses instead; the client displays a safe generic error.

## Tests and reproducible coverage

Frontend, from `frontend/`:

```bash
pnpm format:check
pnpm lint
pnpm typecheck
pnpm test
pnpm test:coverage
pnpm build
```

Backend, from `backend/`:

```bash
test -z "$(gofmt -l cmd internal)"
go vet ./...
mkdir -p coverage
go test -count=1 -race -covermode=atomic -coverprofile=coverage/coverage.out ./...
go tool cover -func=coverage/coverage.out | tee coverage/summary.txt
go tool cover -html=coverage/coverage.out -o coverage/index.html
go test -count=1 ./internal/transport/http -run 'Test(HTTPOpenAPI|OpenAPIHead|ContractChecker)'
go build ./...
```

Do not use `-short` for final checks: it skips real lifecycle subprocess tests.
With Compose running, execute native Node HTTP checks from the repository root:

```bash
node --test scripts/containers.test.mjs
# Override the target when needed: APP_URL=http://localhost:8090 node --test scripts/containers.test.mjs
```

Open generated HTML reports in a browser:
**`frontend/coverage/index.html`** and **`backend/coverage/index.html`**.
Go also generates `backend/coverage/summary.txt`. Reports are Git-ignored and
excluded from production images. Reviewers can reproduce them using the commands
above. Successful GitHub Actions coverage steps upload `frontend-coverage` and
`backend-coverage` artifacts, available from the workflow run.

### Verified results — 2026-11-26

- Frontend: 245 tests; formatting, lint, strict types and build pass. Coverage:
  99.66% statements, 97.95% branches, 100% functions/lines.
- Backend: 632 leaf cases, zero failures/skips; race detector, vet, contracts and
  build pass. Statement coverage: 98.9% overall (calculator/HTTP 100%, expression
  99.3%, executable 90.3%). Counts exclude parent subtest aggregation.
- Actual containers: 27 HTTP smoke cases; Chrome verified all seven operations,
  complex expressions, Enter, error recovery, history/copy/reuse and persistence.

Environment: macOS arm64, Docker 27.4.0 / Compose 2.31.0, Linux arm64 containers,
pinned local tools above. No performance claim or other-architecture verification.

## Architecture and limitations

React controlled components and focused hooks → native fetch → Nginx → Go HTTP
validation → expression parser → calculator domain. Synchronous REST fits small,
deterministic, side-effect-free calculations. Applications build independently;
no database, broker or frontend evaluator is needed. Container base images and CI
actions are pinned. [Architecture decisions](docs/ARCHITECTURE.md) explain details.

- History is plaintext and per-origin, with no cross-tab/device synchronization.
  Clipboard requires HTTPS/localhost support; failures show manual-copy guidance.
- Authentication, rate limiting and public TLS deployment are outside this local
  assessment scope. Digest pins do not prove vulnerability freshness or signatures.
- Duplicate JSON member names retain the last value; see [review notes](docs/REVIEW.md).

See [requirement traceability](docs/REQUIREMENTS.md) and the
[curated AI development record](PROMPTS.md) for assessment scope, technical
decisions and documented development corrections.
