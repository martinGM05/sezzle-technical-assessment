# Architecture Decisions

## Contract-first recovery

The first API implementation preceded OpenAPI. This deviation is recorded in `PROMPTS.md`; `docs/openapi.yaml` is now the source of truth and handler tests cover its documented routes and responses.

## Service design

The service is stateless and uses synchronous REST because calculations are deterministic, immediate, and side-effect free. Pure calculator rules live in `internal/calculator`; HTTP decoding and Problem Details live in `internal/transport/http`.

## Precision

The API uses `float64` and rejects non-finite inputs/results. IEEE 754 rounding makes it unsuitable for monetary amounts; financial systems should use decimal arithmetic or integer minor units.

## Server lifecycle

SIGINT and SIGTERM initiate graceful shutdown: stop accepting connections and
allow active requests up to five seconds to finish. The shutdown context is
independent of the canceled signal context. If the deadline expires, active
connections are closed and the process exits with status 1. A normal shutdown
exits with status 0; `http.ErrServerClosed` is expected, not an operational error.
Bind/serve failures also exit with status 1. The lifecycle waits for the serving
goroutine before returning and uses only the standard library.

Tests use ephemeral loopback listeners, an active request that can be released
or canceled, and isolated subprocesses for real SIGINT/SIGTERM and startup
failure checks. They do not send signals to the test runner. `go test -short`
skips the subprocess tests; the full backend suite includes them.

## Same-origin reverse proxy and CORS

The reported behavior is reproducible: browsers can send `Origin` even for a
same-origin POST. A proxy serving UI and API under one public origin may forward
that header. The previous executable never populated the handler allowlist, so
such a calculation returned `403 CORS_DENIED` even though the browser considered
it same-origin.

Set `ALLOWED_ORIGINS=https://calculator.example` to the exact public origin when
running behind that deployment. Multiple explicit origins may be comma-separated;
surrounding whitespace is trimmed. Configuration is copied once at startup.
Empty configuration denies every nonempty Origin; invalid entries, `*`, and the
opaque `null` origin remain denied. Requests without Origin remain available.
This configuration is also necessary when TLS terminates at the proxy and the
backend connection is plain HTTP. Neither `Host` nor forwarded headers grant
origin access. CORS is not authentication.

A Go integration test exercises a real HTTPS reverse proxy, the real application
constructor and an HTTP backend. It checks allowed and attacker origins, empty
configuration and origin-less clients. This is deployment-level verification;
no frontend or production proxy configuration was added in Phase 2.

## API and numerical policy

A single synchronous `POST /api/v1/calculations` handles seven operations. Add,
subtract, multiply, divide and power require two operands; square root and
percentage require one. Percentage means operand / 100, and finite 0^0 is 1.
Underflow follows IEEE 754 rounding; overflow and non-finite results are rejected.
The backend returns the float64 result without presentation rounding. A future UI
may format it, but must not change calculation semantics or describe float64 as
appropriate for financial transactions.

The HTTP layer limits the complete body to one MiB, including trailing whitespace.
It rejects unsupported media types, malformed/trailing JSON, unknown properties,
incorrect types and missing/null required inputs. Error codes are stable Problem
Details extensions. HEAD has no response body, including errors. JSON duplicate
member names currently retain the last value; this known limitation is recorded
in [the review findings](REVIEW.md), not silently claimed to be rejected.

Authentication and rate limiting are outside this assessment's initial scope.
They would need reconsideration for user-specific data, privileged operations,
quotas or an Internet-facing deployment requiring abuse protection. Server
timeouts and body limits are not substitutes for those controls.

## Packaging and reproducibility

The Docker build uses Go 1.24.2, matching `backend/go.mod` and CI, and produces a
static binary with `CGO_ENABLED=0`. Only production Go sources and the module
file enter the build context. The runtime is distroless static Debian 13, numeric
UID/GID 65532, with an exec-form entrypoint and SIGTERM stop signal. It contains
no shell, curl or Go toolchain; health checks run outside the image.

Both image indexes were resolved with live `docker buildx imagetools inspect` on
2026-11-26 and are pinned in `backend/Dockerfile`:

- `golang:1.24.2-alpine`: `sha256:7772cb5322baa875edd74705556d08f0eeca7b9c4b5367754ce3f2f00041ccee`
- `gcr.io/distroless/static-debian13:nonroot`: `sha256:8ce6dd87e9823c5b48dcfb2c8604714ccb351ce59607327bcc4957448da279b5`

Debian 13 replaces the old unpinned Debian 12 draft because it is listed as
maintained in the [official distroless documentation](https://github.com/GoogleContainerTools/distroless).
Digest verification fixes the downloaded content identity; it is not signature
verification or a vulnerability audit. Pins must be refreshed deliberately when
updating dependencies. The inherited Go toolchain was not upgraded in this phase.

Local validation used Docker Engine 27.4.0 on linux/arm64. It checked actual UID
65532, health, all seven calculations, allowed/denied origins, idle shutdown,
active-request draining and forced closure after the five-second deadline.
Checks ran with a read-only filesystem, all capabilities dropped and
no-new-privileges. Other architectures were not executed locally.

## CI coverage

Frontend and backend jobs remain independent, with `contents: read` permissions.
The backend runs formatting, vet, uncached race-enabled tests with atomic coverage,
coverage reporting and a production build. The coverage profile, function summary
and HTML report are uploaded as `backend-coverage`; the actual function summary
is also written to the job summary. No fabricated metric or arbitrary threshold
is enforced. The standard-library module has no go.sum, so setup-go hashes go.mod
for its cache key.

External actions use immutable commit SHAs. The new upload action was verified
against the [official v4.6.2 tag](https://api.github.com/repos/actions/upload-artifact/git/ref/tags/v4.6.2):
`ea165f8d65b6e75b540449e92b4886f43607fa02`. Existing backend checkout/setup-go pins
were verified as v5.1.0/v6.5.0. Local step execution validates commands and reports,
but cannot prove GitHub runner execution or artifact upload; no push was authorized.
