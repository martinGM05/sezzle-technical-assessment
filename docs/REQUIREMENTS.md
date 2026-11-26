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
- MSW and Docker Compose only if separately justified and approved.
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
