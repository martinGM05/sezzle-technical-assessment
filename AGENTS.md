# Engineering Guidelines

## Scope and delivery

- Keep the solution simple, maintainable, and proportional to the assessment.
- Implement work incrementally and verify observed behavior before reporting it.
- Do not introduce infrastructure that the calculator does not need.
- Do not commit, push, open pull requests, or merge without explicit approval.

## Architecture

- Keep frontend and backend independently buildable.
- Prefer explicit dependencies, high cohesion, low coupling, and composition.
- Keep calculator domain logic separate from HTTP transport.
- Use the Go standard library unless a real requirement justifies a dependency.
- Treat API contracts as a public boundary; document changes before implementing them.

## Quality and testing

- Use test-driven development for business behavior: observe RED, implement GREEN,
  then refactor.
- Prefer behavioral tests over implementation-detail tests.
- Use strict TypeScript and run formatting, linting, type checking, tests, and
  builds before delivery.
- Use `gofmt`, `go vet`, `go test`, race detection, and coverage for Go work.
- Never fabricate test results, coverage, benchmarks, or development history.

## Documentation and language

- Keep source code, comments, documentation, commit messages, and PR descriptions
  in English.
- Record AI prompts used for implementation in `PROMPTS.md`.
- Separate verified assessment requirements from proposed engineering improvements.
