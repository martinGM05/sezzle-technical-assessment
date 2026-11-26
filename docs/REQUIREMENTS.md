# Assessment Requirements

## Verification status

The original Sezzle assessment brief has not been independently provided for
comparison. The items below are recorded from the candidate's supplied planning
brief and must be verified against the original assignment when it becomes
available.

## Reported required deliverables

- A React user interface and a Go backend microservice exposing JSON REST
  responses.
- Addition, subtraction, multiplication, and division.
- Input validation, error handling, invalid-input handling, and division-by-zero
  handling.
- Clean, readable, idiomatic code.
- Frontend and backend unit tests with coverage reports.
- A README with setup instructions, API examples, and architectural decisions.
- A Git repository containing the complete solution and documentation of AI
  prompts used during development.

## Reported optional requirements

- Exponentiation, square root, and percentage operations.
- Docker support.

## Reported quality expectations

- An intuitive, responsive React user interface.
- Maintainable implementation choices proportionate to the assessment.

## Proposed engineering improvements

These are quality decisions, not independently verified assessment requirements:

- An OpenAPI contract before API implementation.
- RFC 9457 Problem Details responses.
- Strict JSON decoding, request limits, timeouts, and explicit CORS origins.
- React Testing Library and GitHub Actions.
- Docker Compose, explicitly approved for Phase 6 final delivery.
- MSW only if separately justified and approved; it was not introduced.
- Structured logging and graceful shutdown.

The project will adopt these only when they remain proportionate to the
assessment scope and time budget.

## Remaining roadmap and scope exclusions

Playwright, k6, Lighthouse and additional performance testing are optional
engineering improvements, not reported assessment requirements. They are removed
from the remaining implementation roadmap, rather than deferred deliverables.
No new browser-testing dependency, load-test suite, performance audit or benchmark
work is planned. Historical prompts and prior verification evidence remain intact.

The keyboard regression is complete using the existing Vitest, React Testing
Library and user-event setup, plus manual Chrome verification against the real Go
backend. Existing frontend/backend quality gates and coverage reporting remain
in scope; no additional tooling is needed to finish this correction.


## Final reported-requirement traceability

The supplied requirements are implemented and have the following evidence. The
original-assignment verification limitation above remains; no unseen requirement
is assumed or certified.

- React responsive UI: `frontend/src/features/calculator/components/`, shared UI
  and token/theme styles; existing component tests and approved visual references.
  Phase 6 verifies the unchanged UI against the containerized service in Chrome.
- Basic arithmetic and optional power/square-root/percentage: Go calculator domain
  and expression engine; 632 backend leaf cases plus real-proxy smoke coverage of
  all seven operations and precedence/parentheses/unary cases.
- REST/JSON validation and division-by-zero handling: `docs/openapi.yaml`, Go HTTP
  tests, native container tests and manual browser error recovery. Transport is
  separate from pure arithmetic; no frontend mathematical evaluator exists.
- Frontend and backend unit tests/coverage: 245 frontend cases, native Go suites,
  fresh HTML reports and documented reproducible commands in README. CI artifact
  uploads are configured for both services; Phase 6 remote execution is pending
  authorized publication, not claimed as already passed.
- Setup/API/architectural documentation: README, ARCHITECTURE and EXPRESSIONS.
- Git solution and AI prompt documentation: repository source and explicitly curated
  PROMPTS.md, with the unedited original preserved in a Git-ignored local archive. Final publication requires explicit authorization; local results do not establish a remote CI pass.
- Optional Docker support: two non-root multi-stage images plus Compose;
  `docker compose up --build` builds and starts the full application.

Expression support, persisted history, themes, Docker Compose, strict TDD,
OpenAPI, Problem Details and CI are approved engineering work, not a claim that
each was mandatory in the original assessment. Playwright, k6, Lighthouse and
additional performance testing remain excluded. No further feature work is
planned as part of final delivery.
