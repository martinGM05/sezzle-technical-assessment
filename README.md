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

## Backend commands

```bash
cd backend
gofmt -l .
go vet ./...
go test -race ./...
go build ./...
```

## Current scope

Phase 1 establishes reproducible tooling, a minimal application shell, quality
checks, and continuous integration. Calculator behavior, HTTP endpoints, OpenAPI,
Docker, end-to-end testing, and performance tooling are intentionally deferred.

See [`docs/REQUIREMENTS.md`](docs/REQUIREMENTS.md) for the distinction between
reported assessment requirements and proposed engineering improvements.
