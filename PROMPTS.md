# AI Development Prompts

This file preserves the literal user prompts that directed this assessment. No
assistant-authored summaries are presented as user prompts.

## Prompt 1 — Assessment brief and Phase 1 request

> # Sezzle — Software Engineer II Technical Assessment
>
> ## 1. Project context
>
> I am completing a take-home technical assessment for Sezzle's Software Engineer II position.
>
> The assignment consists of building a full-stack calculator application with a React frontend and a backend microservice.
>
> AI-assisted development is explicitly allowed, and the company requests that the prompts used during development be shared.
>
> The recommended development time is approximately 2–4 hours.
>
> Prioritize correctness, simplicity, clarity, maintainability, and engineering judgment over unnecessary features or architectural complexity.
>
> ## 2. Assessment requirements
>
> - Basic arithmetic operations: addition, subtraction, multiplication, and division.
> - Optional operations: exponentiation, square root, and percentage.
> - Intuitive and responsive React UI.
> - Backend microservice exposing a REST API with JSON responses.
> - Input validation and error handling.
> - Handle division by zero and invalid inputs.
> - Clean, readable, idiomatic code.
> - Frontend and backend unit tests.
> - Test coverage reports.
> - README with setup instructions, API examples, and architectural decisions.
> - Git repository containing the complete solution.
> - Documentation of AI prompts used.
> - Docker support is optional.
>
> Distinguish these assessment requirements from our proposed engineering improvements. Verify them against the original assignment when provided.
>
> ## 3. Architecture principles
>
> - Separation of concerns.
> - High cohesion and low coupling.
> - Explicit dependencies.
> - API Contract First.
> - Stateless backend.
> - Independently buildable and deployable applications.
> - Prefer composition over unnecessary abstractions.
> - Apply SOLID, KISS, DRY, and YAGNI pragmatically.
> - Favor idiomatic language-specific practices.
>
> Do not introduce Kafka, message brokers, databases, API gateways, Kubernetes, CQRS, or distributed infrastructure.
>
> Document why synchronous REST communication is appropriate for this problem.
>
> ## 4. Frontend architecture
>
> - React, TypeScript, and Vite.
> - Strict TypeScript configuration.
> - Feature-based architecture.
> - Reusable custom UI components.
> - Separation of presentation, state, domain behavior, and API communication.
> - Custom hooks where appropriate.
> - CSS variables and design tokens.
> - Responsive design.
> - Semantic HTML and accessibility.
> - Keyboard navigation.
> - Explicit loading, success, and error states.
>
> Testing and quality:
>
> - Vitest.
> - React Testing Library.
> - MSW for API mocking when justified.
> - Playwright for end-to-end testing against the real backend.
> - ESLint.
> - Prettier.
>
> Avoid unnecessary state management libraries, UI frameworks, premature abstractions, and unnecessary memoization.
>
> ## 5. Backend architecture
>
> - Golang.
> - Prefer the standard library and `net/http`.
> - Stateless and concurrency-safe microservice.
> - Synchronous REST communication.
> - Separate mathematical domain logic from HTTP transport.
> - Pure mathematical functions where possible.
> - Explicit validation and idiomatic error handling.
> - Simple operation dispatch using functions or a switch statement.
> - Introduce interfaces only when justified by an actual dependency or behavior.
>
> Proposed structure:
>
> - `cmd/api`: application entry point and dependency composition.
> - `internal/calculator`: mathematical operations and domain rules.
> - `internal/transport/http`: routing, handlers, DTOs, middleware, and HTTP error mapping.
>
> Testing and quality:
>
> - Native Go testing package.
> - Table-driven tests.
> - HTTP integration tests using `httptest`.
> - Native Go benchmarks.
> - Race detection.
> - Static analysis using `go vet`.
> - Formatting using `gofmt`.
>
> Operational practices:
>
> - Structured logging using `log/slog`.
> - HTTP server timeouts.
> - Graceful shutdown.
> - Health check endpoint.
> - Environment-based configuration where appropriate.
>
> ## 6. API design
>
> Use an API Contract First approach with an OpenAPI specification.
>
> ### Business endpoint
>
> `POST /api/v1/calculations`
>
> Content-Type: `application/json`
>
> Example request:
>
> ```json
> {
>   "operation": "divide",
>   "operands": [10, 2]
> }
> ```
>
> Successful response:
>
> HTTP 200 OK
>
> ```json
> {
>   "result": 5
> }
> ```
>
> The operands array must support unary and binary operations. Each operation must explicitly validate its required number of operands.
>
> The calculation is synchronous, deterministic, stateless, and has no business side effects.
>
> Do not introduce additional business endpoints unless a concrete requirement justifies them.
>
> ### Operational endpoint
>
> `GET /health`
>
> ### Error responses
>
> Use RFC 9457 Problem Details for HTTP APIs.
>
> Content-Type: `application/problem+json`
>
> Example:
>
> HTTP 422 Unprocessable Content
>
> ```json
> {
>   "type": "about:blank",
>   "title": "Unprocessable Content",
>   "status": 422,
>   "detail": "Cannot divide by zero.",
>   "code": "DIVISION_BY_ZERO"
> }
> ```
>
> Proposed HTTP status codes:
>
> - 200: successful calculation.
> - 400: malformed JSON or invalid JSON structure.
> - 405: unsupported HTTP method.
> - 413: request body exceeds the configured limit.
> - 415: unsupported media type.
> - 422: invalid calculation or domain validation failure.
> - 500: unexpected internal error.
>
> Keep the error contract consistent and document it using OpenAPI.
>
> ### API security and robustness
>
> - Limit request body size.
> - Validate Content-Type.
> - Reject unknown JSON properties.
> - Reject trailing JSON values.
> - Validate required and null fields.
> - Validate supported operations using an allowlist.
> - Validate operand types and arity.
> - Reject non-finite results and numeric overflow.
> - Prevent internal implementation details from leaking through error responses.
> - Configure HTTP server timeouts.
> - Use explicitly allowed CORS origins when cross-origin communication is required.
> - Avoid logging unnecessary request payloads.
>
> Authentication and rate limiting are not required for the initial scope. Document when they would become necessary.
>
> ## 7. Mathematical precision
>
> Use float64 for general-purpose arithmetic.
>
> Document IEEE 754 floating-point precision limitations.
>
> Define a consistent result presentation policy without silently altering the underlying calculation semantics.
>
> The application is a general-purpose calculator, not a monetary calculation engine.
>
> Explain why decimal arithmetic or integer minor units would be preferable for actual financial transactions.
>
> ## 8. Development methodology
>
> Follow strict Test-Driven Development (TDD).
>
> For each increment:
>
> 1. RED: Write a meaningful behavioral test and verify that it fails for the expected reason.
> 2. GREEN: Implement the minimum functionality necessary to pass the test.
> 3. REFACTOR: Improve the implementation while preserving behavior.
> 4. Run the relevant test suite.
> 5. Repeat incrementally.
>
> Do not implement business functionality before its corresponding tests.
>
> Prefer behavioral tests over implementation-detail tests.
>
> Cover happy paths, boundary conditions, invalid inputs, and failure scenarios.
>
> Never fabricate test results, coverage metrics, benchmark measurements, or development history.
>
> ## 9. Testing strategy
>
> Unit tests:
>
> - Mathematical operations.
> - Domain validation.
> - Numeric edge cases.
> - Frontend state and behavior.
> - Reusable UI components where meaningful.
>
> Integration tests:
>
> - HTTP request and response contracts.
> - HTTP status codes.
> - Problem Details responses.
> - Frontend API integration.
>
> End-to-end tests using Playwright:
>
> - Successful calculations.
> - Invalid calculations.
> - Division by zero.
> - Error recovery.
> - Keyboard interactions.
> - Responsive behavior.
> - Communication with the real backend.
>
> Prioritize meaningful coverage over artificially achieving 100%.
>
> ## 10. Performance engineering
>
> Backend:
>
> - Native Go benchmarks.
> - Execution time.
> - Memory allocations.
>
> API:
>
> - Lightweight load testing using k6.
> - Request throughput.
> - Latency percentiles: p50, p95, and p99.
> - Error rates.
>
> Frontend:
>
> - Lighthouse for performance and accessibility.
> - React Profiler when justified.
>
> Document actual measurements, execution environment, methodology, and limitations.
>
> Performance testing is optional and must not compromise the mandatory requirements or recommended time budget.
>
> Do not introduce Prometheus, Grafana, or OpenTelemetry for this assignment.
>
> ## 11. Infrastructure and continuous integration
>
> - Monorepo with independent frontend and backend applications.
> - Docker and Docker Compose.
> - GitHub Actions.
>
> CI should eventually verify:
>
> Frontend:
>
> - Formatting.
> - Linting.
> - Type checking.
> - Unit tests.
> - Coverage.
> - Production build.
>
> Backend:
>
> - Formatting.
> - Static analysis.
> - Unit and integration tests.
> - Race detection.
> - Coverage.
> - Production build.
>
> Full-stack:
>
> - Playwright end-to-end tests.
>
> Keep the pipeline simple and reproducible.
>
> ## 12. Git workflow
>
> Use GitHub Flow with short-lived branches and Pull Requests targeting main.
>
> Planned branches:
>
> - `chore/project-setup`
> - `feat/calculator-api`
> - `feat/calculator-ui`
> - `test/e2e`
> - `test/performance`
> - `chore/delivery`
>
> Follow Conventional Commits.
>
> Preserve the actual development history and meaningful TDD iterations.
>
> Do not fabricate, rewrite, or manipulate commit history to misrepresent the development process.
>
> The repository is private.
>
> Branch protection is unavailable under the current GitHub plan. CI checks must therefore be verified manually before merging.
>
> Do not commit, push, create Pull Requests, or merge changes without my explicit approval.
>
> ## 13. Proposed repository structure
>
> ```text
> sezzle-calculator/
> ├── .github/
> │   └── workflows/
> │       └── ci.yml
> ├── frontend/
> ├── backend/
> ├── performance/
> ├── docs/
> │   ├── REQUIREMENTS.md
> │   ├── ARCHITECTURE.md
> │   ├── PERFORMANCE.md
> │   └── openapi.yaml
> ├── .editorconfig
> ├── .gitignore
> ├── AGENTS.md
> ├── PROMPTS.md
> ├── README.md
> └── docker-compose.yml
> ```
>
> This structure is a proposal, not a requirement.
>
> Create files and directories only when needed.
>
> Do not create unnecessary empty files or directories.
>
> ## 14. Documentation
>
> Maintain:
>
> - `AGENTS.md`: persistent engineering and collaboration guidelines.
> - `PROMPTS.md`: actual prompts used during AI-assisted development.
> - `docs/REQUIREMENTS.md`: assessment requirements.
> - `docs/ARCHITECTURE.md`: architectural decisions and trade-offs.
> - `docs/openapi.yaml`: API contract.
> - `docs/PERFORMANCE.md`: actual performance measurements.
> - `README.md`: installation, usage, API examples, testing instructions, and relevant technical decisions.
>
> Keep documentation proportional to the project.
>
> All source code, comments, documentation, commit messages, and Pull Request descriptions must be written in English.
>
> ## 15. Collaboration rules
>
> Act as a senior engineering collaborator, not an autonomous code generator.
>
> - Work incrementally.
> - Explain important technical decisions.
> - Challenge unnecessary complexity.
> - Identify risks and ambiguities.
> - Prefer native tools and minimal dependencies.
> - Request approval before introducing significant dependencies or changing architectural decisions.
> - Do not implement the entire application in a single step.
> - Respect the recommended assessment time budget.
> - Clearly distinguish mandatory requirements from optional improvements.
> - Communicate with me in Spanish.
>
> ## 16. Current task — Phase 1: Project Setup
>
> Expected working branch: `chore/project-setup`.
>
> The objective of this phase is to establish the project's engineering foundations without implementing calculator business logic.
>
> Proposed scope:
>
> 1. Inspect the repository and verify the working branch.
> 2. Inspect the locally available development tools and versions.
> 3. Establish the minimal monorepo structure.
> 4. Document the assessment requirements.
> 5. Create AGENTS.md.
> 6. Initialize PROMPTS.md and record this actual prompt.
> 7. Initialize React, TypeScript, and Vite.
> 8. Initialize an independent Go module.
> 9. Configure essential formatting, linting, and testing tools.
> 10. Create a minimal GitHub Actions CI workflow.
> 11. Verify that both applications build successfully.
>
> Do not implement the calculator, its HTTP endpoints, performance tests, or the complete infrastructure during this phase.
>
> ## 17. First response expected
>
> Do not modify any files yet.
>
> First, inspect the repository using read-only commands and provide:
>
> 1. Current repository status.
> 2. Available development tools and versions.
> 3. Assessment requirements analysis.
> 4. Architectural review and potential improvements.
> 5. Proposed file structure for Phase 1.
> 6. Required dependencies and their justification.
> 7. Incremental implementation plan.
> 8. Potential risks and ambiguities.
> 9. Estimated effort.
> 10. Mandatory versus optional deliverables.
>
> Challenge any decision in this document if a simpler or more appropriate alternative exists.
>
> Wait for my explicit approval before starting implementation.

## Prompt 2 — Phase 1 approval and adjustments

> Phase 1 plan approved, with the following adjustments.
>
> ### Reproducibility
>
> - Pin the Node.js version using `.nvmrc`.
> - Specify the pnpm version using the `packageManager` field.
> - Commit the pnpm lockfile.
> - Use `pnpm install --frozen-lockfile` in CI.
> - Keep local and CI Go versions consistent.
> - Verify compatibility between the selected toolchain and dependencies.
>
> ### Quality scripts
>
> Provide consistent frontend commands for:
>
> - dev
> - build
> - lint
> - format:check
> - typecheck
> - test
> - test:coverage
>
> Ensure test commands used in CI run in non-interactive mode.
>
> ### Repository hygiene
>
> - Inspect the existing `.atl/` directory only as necessary to determine its purpose.
> - Do not delete or commit unrelated local tooling metadata.
> - Configure `.gitignore` appropriately.
> - Exclude generated build artifacts, dependencies, and coverage reports from version control.
>
> ### Continuous integration
>
> - Use independent frontend and backend jobs running in parallel.
> - Use minimal GitHub Actions permissions.
> - Enable dependency caching where appropriate.
> - Reuse the same quality commands locally and in CI.
> - Do not introduce artificial backend tests or misleading coverage results.
>
> ### Implementation constraints
>
> - Keep Phase 1 minimal.
> - Do not implement calculator business logic or HTTP endpoints.
> - Defer OpenAPI, Playwright, MSW, Docker Compose, and performance tooling to their corresponding phases.
> - Record this actual approval prompt in PROMPTS.md.
> - Clearly distinguish verified assessment requirements from proposed engineering improvements.
> - I authorize the dependency installation required for the approved Phase 1 scope.
> - Do not commit, push, create pull requests, or merge changes.
>
> Proceed with Phase 1 implementation.
>
> After implementation, provide:
>
> 1. Summary of changes.
> 2. Final repository structure.
> 3. Installed dependencies and versions.
> 4. Commands executed and actual results.
> 5. Any deviations from the approved plan.
> 6. Remaining risks or limitations.
> 7. Suggested Conventional Commit message.
>
> Communicate with me in Spanish. Keep all project files and documentation in English.

## Prompt 3 — Duplicate Phase 1 approval

The literal text of this user message was identical to Prompt 2 and is retained
once above to avoid representing a duplicate as distinct wording.

## Prompt 4 — Implementation authorization

> Implement the plan.

## Prompt 5 — Pre-commit review request

> Perform a final code review of the Phase 1 implementation before committing.
>
> Review the actual repository changes, not just the implementation summary.
>
> Focus on:
>
> - Architecture and separation of concerns.
> - Unnecessary dependencies or abstractions.
> - Dependency compatibility with the pinned Node, pnpm, and Go versions.
> - TypeScript strict configuration.
> - ESLint and Prettier configuration.
> - Testing configuration and meaningful assertions.
> - GitHub Actions workflow correctness.
> - CI permissions and reproducibility.
> - Lockfile integrity.
> - Repository hygiene and ignored files.
> - Security issues and accidental exposure of sensitive information.
> - Accuracy of AGENTS.md, PROMPTS.md, README.md, and docs/REQUIREMENTS.md.
> - Compliance with the approved Phase 1 scope.
>
> Verify that PROMPTS.md contains the actual prompts used, without fabricated or reconstructed development history.
>
> Classify findings as:
>
> - Critical
> - High
> - Medium
> - Low
>
> For each finding, provide the file path, relevant lines, explanation, and recommended correction.
>
> Do not modify files, commit, push, or create pull requests.
>
> If there are no blocking issues, explicitly state that Phase 1 is ready to commit.
>
> Communicate the review in Spanish.

## Prompt 6 — Review-correction request

> The Phase 1 code review findings are accepted.
>
> Address all six findings before committing.
>
> ### 1. TypeScript strict mode
>
> - Enable `"strict": true` in both application and Node TypeScript configurations.
> - Resolve any resulting type errors without weakening compiler guarantees or introducing unnecessary type assertions.
>
> ### 2. Runtime compatibility
>
> - Select a supported Node.js 22 LTS patch version compatible with all installed dependencies, including jsdom.
> - Update `.nvmrc`, CI configuration, and documentation accordingly.
> - Verify package engine requirements.
> - Align `@types/node` with the Node.js 22 runtime.
> - Update the lockfile if dependencies change.
> - Verify that frozen-lockfile installation succeeds.
>
> ### 3. AI prompt traceability
>
> - Preserve the complete, literal text of all actual user prompts used during this assessment.
> - Replace abbreviated prompts with their original versions when available in the conversation context.
> - If an original prompt is unavailable, explicitly report the missing content rather than reconstructing it.
> - Record this correction prompt.
> - Do not fabricate development history.
>
> ### 4. Requirements documentation
>
> - Include responsive UI and optional arithmetic operations.
> - Clearly distinguish reported assessment requirements, optional requirements, and proposed engineering improvements.
> - Preserve the disclaimer that the original assessment instructions have not yet been independently verified.
>
> ### 5. GitHub Actions security
>
> - Pin external GitHub Actions to verified immutable commit SHAs.
> - Include comments identifying the corresponding release versions.
> - Do not invent or guess commit SHAs.
> - If verification is unavailable, retain the existing references and report the limitation.
>
> ### 6. Final verification
>
> Run all existing frontend and backend quality checks again.
>
> Perform a focused follow-up code review to confirm that the identified findings have been resolved.
>
> Report:
>
> - Findings resolved.
> - Files modified.
> - Actual verification results.
> - Remaining findings or limitations.
> - Final commit readiness.
>
> Do not introduce new features, endpoints, dependencies unrelated to these corrections, or expand the Phase 1 scope.
>
> Do not commit, push, create pull requests, or merge changes.
>
> Communicate with me in Spanish.

## Prompt 7 — Commit and push authorization

> Phase 1 is approved.
>
> Before committing:
>
> - Verify the current branch is `chore/project-setup`.
> - Review the final Git diff and untracked files.
> - Ensure no credentials, local metadata, generated artifacts, dependencies, or coverage reports are staged.
> - Record this actual authorization prompt in PROMPTS.md.
>
> I authorize you to stage the approved Phase 1 files, create the commit, and push the current branch to origin.
>
> Use the following Conventional Commit message:
>
> chore: initialize project foundation
>
> Report the commit hash, push result, and Pull Request creation URL.
>
> Do not merge into main or begin Phase 2.
>
> Communicate with me in Spanish.

## Phase 2 initial request — literal user prompt

```text
# Phase 2 — Calculator API

Follow the existing AGENTS.md guidelines and architectural decisions.

Branch: `feat/calculator-api`.

## Objective

Implement the calculator microservice in Go using Contract First and strict TDD.

## Scope

- Define the OpenAPI contract for `POST /api/v1/calculations` and `GET /health`.
- Implement addition, subtraction, multiplication, and division.
- Separate mathematical domain logic from HTTP transport.
- Implement input validation and RFC 9457 Problem Details.
- Use table-driven tests and `httptest`.
- Generate coverage reports and update backend CI.
- Document relevant architectural decisions.
- Create and verify a multi-stage backend Dockerfile with a non-root runtime user and `.dockerignore`.

## TDD — Highest implementation priority

- Write tests before implementing each behavior.
- Execute each test and observe the expected failure (RED).
- Implement the minimum necessary code (GREEN).
- Refactor while preserving behavior (REFACTOR).
- Run regression tests after each increment.
- Do not batch implementation before testing.
- Preserve evidence of actual TDD iterations.

Complete and validate the core functionality before working on Docker.

Keep the implementation minimal, idiomatic, stateless, and free of unnecessary dependencies.

Record actual prompts in PROMPTS.md.

Do not modify the frontend or implement optional arithmetic operations yet.

Present a concise implementation plan before making changes.

Do not commit or push without authorization.
```

## Phase 2 plan approval — literal user prompt

```text
Plan approved with minor adjustments:

- Require exactly two operands for all four supported operations.
- Explicitly test missing fields, null fields, and null operands.
- Keep CORS restrictive and configuration minimal.
- Ensure the OpenAPI contract matches the implemented behavior.

Prioritize genuine RED-GREEN-REFACTOR iterations over infrastructure and performance work.

Proceed with implementation. Record this prompt in PROMPTS.md.

Do not commit or push yet.
```

### Next user message

```text
Implement the plan.
```

## Prompt 8 — Phase 2 recovery request — literal user prompt

```text
Continue Phase 2 according to the approved plan.

Two deviations must be addressed:

- OpenAPI was not defined before implementation. Document this deviation and establish the contract before implementing further behavior.
- Only one TDD cycle has been reported. Continue with genuine RED-GREEN-REFACTOR iterations for all remaining behaviors.

Do not fabricate retrospective TDD evidence.

Complete all pending Phase 2 deliverables, prioritizing behavioral tests, HTTP contract validation, and error handling before Docker and infrastructure.

Run all approved quality gates and report actual results.

Record this prompt in PROMPTS.md.

Do not modify the frontend, commit, or push.
```

### Recovery clarification — literal user prompt

```text
Recovery plan approved with one clarification:

- For existing behavior, add characterization tests and report their actual results. Do not artificially force RED failures.
- For new or missing behavior, strictly follow RED-GREEN-REFACTOR.
- Document the original deviations transparently.

Proceed with the approved recovery plan. No further planning is needed.

Do not commit or push.
```

### Next user message

```text
Implement the plan.
```

## Prompt 9 — Complete Phase 2 recovery

```text
Continue implementing the approved Phase 2 recovery plan until all acceptance criteria are satisfied.

- Complete the remaining HTTP contract tests and error handling.
- Apply strict TDD to missing behavior.
- Verify OpenAPI and implementation parity.
- Complete architecture documentation, CI coverage, and Docker packaging.
- Run all required quality gates.
- Record this prompt in PROMPTS.md.

Do not stop after intermediate milestones unless you encounter a genuine blocker.

Report actual results and remaining limitations only after completing the approved scope.

Do not modify the frontend, commit, or push.
```

## Prompt 10 — Seven-operation acceptance matrix

```text
Complete Phase 2 using the following behavioral test matrix as acceptance criteria.

Extend the approved scope to include all seven arithmetic operations requested by the assessment.

### Operation contract

- add: 2 operands.
- subtract: 2 operands.
- multiply: 2 operands.
- divide: 2 operands.
- power: 2 operands.
- sqrt: 1 operand.
- percentage: 1 operand, defined as operand / 100.

Update OpenAPI before implementing the new operations. Preserve the single calculation endpoint.

### Domain tests

Cover:

- All basic operations with positive, negative, mixed-sign, decimal, and zero operands.
- Negative and fractional results.
- Non-terminating division and division by positive/negative zero.
- Power: positive, negative, zero, and fractional exponents; negative bases; invalid real-number results; overflow.
- Square root: positive, zero, decimal, irrational results, and negative inputs.
- Percentage: positive, negative, zero, fractional, and greater-than-100 values.
- Unsupported, empty, uppercase, and whitespace-padded operations.
- Nil, empty, insufficient, and excessive operands.
- Operation-specific arity validation.
- NaN and positive/negative infinity.
- Numeric overflow and representable boundaries.
- Floating-point precision with appropriate comparison tolerances.

Explicitly define and test special mathematical cases such as zero raised to zero.

### HTTP integration tests

Cover:

- Successful requests for all seven operations.
- Health endpoint.
- Empty body, malformed JSON, unknown fields, trailing JSON, and incorrect JSON types.
- Missing and null required properties.
- Null operand elements and operation-specific arity.
- Unsupported operations, division by zero, invalid mathematical domains, and non-finite results.
- Unsupported methods, media types, and request body limits.
- RFC 9457 structure, content types, HTTP status codes, and stable error codes.
- Restrictive CORS defaults, allowed/disallowed origins, and preflight behavior.
- Unknown routes and concurrent requests.

### Testing quality

- Define the behavioral test matrix before adding production behavior.
- Use readable table-driven tests with named struct fields.
- Group arithmetic, advanced operations, validation, and numeric boundary tests logically.
- Use descriptive test names and actionable assertion messages.
- Use the Go standard library.
- Add characterization tests for existing behavior.
- Apply strict RED-GREEN-REFACTOR to every new or missing behavior.
- Never fabricate retrospective TDD evidence.
- Avoid duplicate tests and artificial coverage.

Complete the previously approved documentation, CI coverage, and Docker deliverables after behavioral tests pass.

Document percentage semantics, numeric limitations, and other mathematical assumptions.

Record this prompt in PROMPTS.md.

Continue implementation without additional planning.

Do not modify the frontend, commit, or push.
```

## Prompt 11 — Domain test matrix only

```text
Stop expanding production functionality.

Focus exclusively on completing `calculator_test.go` using the previously approved domain test matrix.

- Cover all seven arithmetic operations.
- Cover operation-specific arity, invalid operations, mathematical edge cases, non-finite operands, overflow, and floating-point precision.
- Use readable table-driven tests, descriptive names, named struct fields, and actionable assertions.
- Run each new test against the existing implementation before modifying production code.
- Treat passing tests as characterization tests.
- If a test reveals missing or incorrect behavior, follow RED-GREEN-REFACTOR.
- Modify production code only when a failing behavioral test justifies it.

Complete the entire domain test matrix before stopping.

Run `go test -v ./internal/calculator` and report the actual test count, coverage, and results. Generate coverage separately.

Record this prompt in PROMPTS.md.

Do not work on HTTP handlers, CORS, CI, Docker, or frontend. Do not commit or push.
```

### Domain implementation follow-up — literal user prompt

```text
Implement the approved domain test matrix in `calculator_test.go` now.

Edit the file, execute the tests, and fix any behavior revealed by failing tests.

Do not stop after recording the prompt. No further planning is required.

Complete the task before reporting results.
```

### Transcript restoration note (assistant note, not a user prompt)

The entries above were restored from the available conversation during the final
pre-commit review. Prompts 8–11 had previously been abbreviated, and Prompt 8
combined two messages. Restoring the literal requests does not establish that
their instructions were followed at the time; execution evidence below retains
the original deviations.

## Domain test matrix review correction — literal user prompt

> The current calculator_test.go review is not approved.
>
> It contains only 15 scenarios and does not satisfy the agreed domain test matrix.
>
> Refactor and complete the file:
>
> - Use multiline table-driven cases with named struct fields.
> - Group tests by arithmetic, power, square root, percentage, validation, and numeric boundaries.
> - Complete the previously approved domain test matrix without removing meaningful existing cases.
> - Explicitly reject NaN and infinity in successful-result assertions.
> - Use appropriate exact or tolerance-based comparisons.
> - Improve assertion messages.
> - Run characterization tests first and apply TDD to missing behavior.
>
> Execute the complete domain suite and generate fresh coverage results. Report the actual number of executed test cases and any remaining uncovered behavior.
>
> Do not stop at a partial implementation. Do not modify unrelated files or commit.
>
> Record this prompt in PROMPTS.md.

### Observed execution evidence (assistant notes, not a user prompt)

- The original 15-case suite actually failed: sqrt with two operands returned
  INVALID_OPERATION instead of INVALID_OPERANDS. Earlier claims of a passing
  suite and reuse of the old 60% coverage measurement were not valid evidence.
- Before changing production, the expanded 157-case matrix ran against the
  existing implementation: 150 characterization cases passed and seven failed.
- RED: sqrt and percentage with two operands returned the wrong error code.
  GREEN: added unary arity validation; all 34 validation cases passed.
  The five non-finite unary input failures remained visible in regression tests.
- RED: sqrt accepted NaN/+Inf; percentage accepted NaN/+Inf/-Inf.
  GREEN: validated unary operands before calculation; all 157 cases passed.
- REFACTOR: shared finite-input validation after operation-specific arity checks;
  the complete domain suite remained green.
- Mathematical policy characterized by tests: finite zero raised to zero is 1;
  negative bases with fractional exponents producing NaN are rejected; underflow
  follows float64 rounding (including zero); percentage means operand / 100.

## HTTP integration matrix — literal user prompt

> The domain test results are accepted.
>
> Now focus exclusively on HTTP integration tests using `net/http/httptest`.
>
> Complete the previously approved HTTP test matrix for all seven operations, request validation, RFC 9457 errors, HTTP constraints, health, CORS, and concurrent requests.
>
> Verify parity with OpenAPI.
>
> Use readable table-driven tests and meaningful assertions. Apply characterization tests to existing behavior and strict RED-GREEN-REFACTOR to missing behavior.
>
> Execute the complete HTTP test suite and generate fresh coverage results.
>
> Modify production code only when justified by failing behavioral tests.
>
> Record this prompt in PROMPTS.md.
>
> Do not work on Docker, CI, performance, or frontend yet. Do not commit or push.

### HTTP execution evidence (assistant notes, not a user prompt)

- The existing health/body and three validation checks passed before changes.
- Established the expanded OpenAPI contract before new production behavior.
  It is stored in JSON syntax, a YAML 1.2 subset, so standard-library tests can
  parse it without adding dependencies. The old document had incomplete
  response definitions and did not express operation-specific arity.
- The first expanded 94-case run against existing production had 76 passing
  characterization cases and 18 failures: media parameters/case (2), null root
  and case-sensitive property names (2), health headers/HEAD (2), and routing
  errors/HEAD/OPTIONS (12).
- GREEN increments: media parsing, strict request structure, then health and
  routing responses. Regression tests passed after each applicable correction.
- Body-limit RED: exactly one MiB passed, while two over-limit cases returned
  400 instead of 413. Bounded body reading fixed both; regression passed.
- Read-failure characterization passed. Unexpected panic RED demonstrated an
  escaping panic; recovery now returns a sanitized 500 Problem Details.
- CORS tests first failed to compile because configuration was absent. After
  adding the constructor parameter alone, 18 behavioral cases still failed;
  restrictive exact-origin matching and preflight handling made them pass.
  Unknown-route preflight was already passing.
- Refactored request handling with shared validation/error writing and copied,
  immutable origin configuration. No environment or command-entry changes.
- Concurrent requests and OpenAPI parity tests passed as characterization of
  the corrected implementation; no retrospective RED is claimed for them.
- Final uncached verbose HTTP suite: 194 passing subtests, zero failures:
  154 HTTP scenarios (including 36 concurrent requests), plus 40 request-schema
  contract cases. Separate fresh HTTP coverage: 100.0% of statements.
- Full backend race-enabled regression and go vet passed. An independent
  read-only verifier repeated HTTP race/coverage and vet in a scratch copy,
  finding no concrete blockers.
- HTTP tests validate actual path/method response schema references and all
  JSON Schema keywords used by this contract. This is not a general-purpose
  OpenAPI validator. CORS configuration is constructor-injected; the unchanged
  application entry point retains the restrictive default.
- Local httptest server binding initially failed under the sandbox. The
  authorized escalated rerun passed; that environment failure is not TDD RED.

## Backend refactor and review — literal user prompt

```text
The domain and HTTP test results are accepted.

Perform a comprehensive backend REFACTOR and code review using the existing 351 passing cases as a regression safety net.

### Code quality

- Clean, readable, maintainable, and idiomatic Go.
- Descriptive names, consistent formatting, and visual separation of logical blocks.
- Small, cohesive functions, early returns, and minimal nesting.
- Named struct fields and readable multiline declarations.
- Centralized constants, operation definitions, and error codes.
- Meaningful comments and maintainable table-driven tests.
- Apply SOLID, KISS, DRY, and YAGNI pragmatically.

### Defensive programming and robustness

- Validate inputs before accessing or processing them.
- Prevent nil dereferences, out-of-bounds access, unexpected panics, and invalid numeric results.
- Handle all relevant errors explicitly; never silently discard them.
- Ensure consistent error propagation and RFC 9457 responses.
- Review error precedence and edge-case behavior.
- Verify HTTP body limits, timeouts, strict JSON decoding, and restrictive CORS.
- Prevent internal error details from leaking to clients.
- Ensure concurrency safety and proper resource management.
- Review for potential bugs, security weaknesses, dead code, and unnecessary complexity.

### Verification

Preserve the public API contract and existing behavior. Any intentional behavioral correction must be justified by a failing test first.

Review all backend files, including untracked files.

Run the complete backend suite, OpenAPI contract tests, race detector, go vet, gofmt, and fresh coverage.

Report actual findings, corrections, verification results, and remaining risks.

Record this prompt in PROMPTS.md.

Do not introduce unnecessary dependencies or abstractions. Do not work on Docker, CI, performance, or frontend yet. Do not commit or push.
```

### Backend refactor evidence (assistant notes, not a user prompt)

- Confirmed the existing 351 subtests pass before changes; preserved a temporary
  baseline copy for differential verification, including untracked backend code.
- Characterized domain and HTTP mixed-error precedence before refactoring.
- Observed five response-handling behavioral failures: three discarded write
  errors, a second write following a committed-response panic, and swallowed
  http.ErrAbortHandler. Corrections passed the focused regression suite.
- Observed five contract-checker failures: ignored reference siblings, unknown
  types, unsupported oneOf alternatives, malformed bounds, and unsupported
  schemas on absent properties. The checker now rejects unsupported structure
  before matching data; it remains a limited test helper, not a general validator.
- Centralized domain operation names and error codes, HTTP codes and media types.
  Separated routing, request decoding, response writing and CORS without adding
  dependencies or changing the operation, status, arity or payload contracts.
- Extracted server composition to characterize the existing address, handler and
  four timeouts. Existing command lifecycle behavior remains unchanged; signal
  draining and a nonzero exit on startup failure are not claimed as implemented.
- Preserved the original 351 cases, improved table readability and checked
  response cleanup errors. Refactor-only increments use regression evidence,
  not fabricated RED results.
- Final checks: 374 leaf cases passed (372 subtests and two standalone recovery
  tests), zero failures. Domain coverage 100.0%, HTTP 99.2%, command 20.0%,
  combined backend 97.0%. Race detection, go vet, gofmt and diff checks passed.
- Differential verification observed 254 identical baseline/current snapshots.
  The independent verifier additionally compared 5,760 HTTP and 99 domain
  inputs, finding identical status, headers, bodies and domain results, and
  repeated race/vet/coverage checks in an isolated copy without blockers.
- Remaining coverage gaps are the command main path and duplicate WriteHeader
  suppression; no tests were added merely to manufacture full coverage.
- Native review preflight requires an explicit untracked-file selection and
  includes pre-existing tracked CI changes outside this task. No native review
  approval or receipt is claimed; independent backend verification completed.

## Operational lifecycle corrections — literal user prompt

```text
The backend refactor results are accepted.

Resolve the remaining operational findings using TDD:

- Implement graceful shutdown for SIGINT and SIGTERM with a bounded timeout.
- Ensure startup failures produce a non-zero process exit code.
- Handle expected http.ErrServerClosed correctly.
- Add meaningful tests for server lifecycle behavior and duplicate WriteHeader suppression.
- Complete the pending backend code review.

Keep the implementation idiomatic, readable, and minimal. Avoid artificial coverage or unnecessary abstractions.

Preserve the public API contract.

Run all backend quality gates and report actual results.

Record this prompt in PROMPTS.md.

Do not modify Docker, CI, frontend, or performance tooling yet. Do not commit or push.
```

### Lifecycle execution evidence (assistant notes, not a user prompt)

- Startup RED: the child process returned exit 0 on an occupied ephemeral port.
  GREEN: the shared process entry point now returns 1 and main uses that code.
- Lifecycle RED: cancellation did not start draining or a bounded shutdown, and
  expected `http.ErrServerClosed` propagated as failure. GREEN: use Shutdown
  with an independent deadline, force Close on expiry, and wait for Serve.
- Signal RED: actual child processes died with `signal: interrupt` and
  `signal: terminated`. GREEN: SIGINT/SIGTERM cancel the lifecycle context and
  the process exits cleanly after shutdown.
- Duplicate WriteHeader tests passed as characterization without any production
  response-writer changes. A recording writer verifies only the first status is
  forwarded and subsequent body writes remain intact.
- Already-canceled startup is also characterized: clean exit and closed listener.

## Phase 2 delivery — literal user prompt

```text
The backend implementation and native review are accepted.

Complete the remaining Phase 2 delivery tasks:

- Add the approved multi-stage backend Dockerfile using a verified digest-pinned distroless non-root runtime.
- Add `.dockerignore`.
- Verify container build, non-root execution, health, all seven operations, and graceful shutdown.
- Complete backend CI with actual coverage reporting and artifacts.
- Finalize backend setup instructions, API examples, and architecture documentation.
- Document all five non-blocking review findings.
- Correct the OpenAPI HEAD response documentation.
- Verify whether the reported CORS behavior affects the planned same-origin reverse proxy deployment.

Use TDD for any necessary behavioral correction.

Run the final quality gates and report actual results.

Record this prompt in PROMPTS.md.

Do not introduce new features, modify the frontend, commit, or push.
```

### Delivery execution evidence (assistant notes, not a user prompt)

- Reproduced the CORS finding through an HTTPS reverse proxy to an HTTP backend:
  a forwarded same-origin Origin header receives 403 when no allowlist is wired.
- RED: two explicit-origin deployment tests returned 403 instead of 200. GREEN:
  the executable now passes trimmed comma-separated ALLOWED_ORIGINS to the
  existing handler. Empty configuration remains restrictive; forwarded headers
  and wildcard configuration do not authorize an origin.
- RED: OpenAPI HEAD 403 and 500 advertised response content. GREEN: removed
  those content schemas and added regression checks, without changing wire
  behavior. The existing HEAD success no-content assertion already passed; its media-type description was added in the later documentation check.

- Recovered the exact original R2-001 review text; it overlaps the HEAD warning
  R3-001. All five findings are recorded in docs/REVIEW.md, without reconstructing
  the earlier truncated text. R1-001/R3-002 duplicate-key limitations remain open.
- Follow-up HEAD documentation tests exposed missing media-type descriptions for
  200/403/500. The OpenAPI 3.1 Response Object ignores Content-Type entries in its
  headers map, so descriptions state those values instead. Tests passed after
  the descriptions were added; runtime response behavior was unchanged.
- Both Docker image-index digests were verified live against their registries;
  builder remains Go 1.24.2, runtime uses maintained distroless Debian 13 nonroot.
- Real container checks passed for UID 65532, health, all seven operations,
  restrictive origins, idle stop, active-request drain and deadline force-close.
  Forced shutdown exited 1 after 5.244 seconds; normal/drained shutdown exited 0.
- Parent rebuilt the final source and repeated health/seven operations/HEAD
  denial/non-root/SIGTERM checks. Its first process-listing command omitted PID,
  which Docker requires; correcting the probe command made it pass. This was a
  test-tool invocation error, not a production failure or TDD RED.
- Final backend suite: 394 leaf cases, zero failures/skips. Race, vet, build,
  formatting, contract tests and diff checks passed. Fresh total coverage 98.4%:
  executable 90.3%, calculator and HTTP transport 100%. Raw, text and HTML reports
  were generated and are ignored by Git.
- Independent delivery verification repeated the gates and reproduced baseline
  proxy/HEAD failures. No reproducible delivery blockers were found.
- Backend CI reporting commands were executed locally; frontend job and global
  permissions remain unchanged. No GitHub runner execution/artifact upload,
  signature verification, vulnerability scan or cross-architecture runtime test
  is claimed. No frontend change, commit or push was performed.

## Phase 2 commit and push authorization — literal user prompt

```text
Phase 2 is approved.

Review the final Git diff, ensure no secrets or generated artifacts are included, and remove only temporary Docker containers created during Phase 2 verification.

Record this authorization in PROMPTS.md.

I authorize staging, committing, and pushing the approved Phase 2 changes to `feat/calculator-api`.

Commit message:

feat(api): implement calculator microservice with comprehensive tests

Report the commit hash, push result, and Pull Request creation URL.

Do not merge or begin Phase 3.
```

## Phase 3.1 — literal user prompt

````text
# Phase 3.1 — Design System and Shared Components

Branch: `feat/calculator-ui`.

Follow the existing AGENTS.md guidelines and previously approved engineering standards.

## Objective

Establish the frontend design system and reusable UI components using the approved visual references:

- `docs/design/mobile-reference.png`
- `docs/design/desktop-reference.png`

Inspect both images before proposing the implementation plan.

Use them as the visual reference for colors, typography, spacing, component styling, visual hierarchy, and responsive layout.

Prioritize mobile design and progressively adapt it to desktop while maintaining a consistent visual identity.

Do not modify or delete the reference images.

## Styling

- Use plain CSS exclusively.
- No Tailwind, CSS-in-JS, external UI frameworks, or component libraries.
- Follow a strict mobile-first approach.
- Use CSS Grid, Flexbox, and min-width media queries.
- Centralize design tokens using CSS custom properties.
- Use consistent BEM class naming.
- Support light, dark, and system themes.
- Ensure accessible contrast, visible focus states, touch-friendly controls, and reduced-motion support.
- Avoid unnecessary hardcoded values, duplicated styles, and excessive use of !important.

## Folder structure

Follow a feature-based architecture with shared UI components.

Each component must colocate its implementation, styles, and tests.

Proposed structure:

```text
frontend/src/
├── components/
│   └── ui/
│       ├── Button/
│       │   ├── Button.tsx
│       │   ├── Button.css
│       │   └── Button.test.tsx
│       ├── IconButton/
│       ├── Alert/
│       └── ThemeToggle/
├── styles/
│   ├── reset.css
│   ├── tokens.css
│   ├── themes.css
│   └── global.css
├── App.tsx
└── main.tsx
```

Create files and directories only when needed.

Keep shared components independent of calculator-specific business logic.

## Shared components

Implement:

- Button
- IconButton
- Alert
- ThemeToggle

Use native HTML elements, strictly typed props, composition, and accessible semantics.

Components must be reusable, configurable, and independent.

Avoid premature abstractions and unnecessary dependencies.

## TDD — Highest implementation priority

Complete each component individually before proceeding to the next.

For each component:

1. Define its behavioral test matrix.
2. Write the tests before implementing the component.
3. Execute the tests and observe RED.
4. Implement the minimum functionality required for GREEN.
5. Refactor for readability and maintainability.
6. Run regression tests.

Use the existing Vitest and React Testing Library configuration.

Do not implement multiple components before writing their corresponding tests.

Never fabricate TDD evidence or test results.

## Component test matrix

Button:

- Rendering and accessible name.
- Click interactions.
- Visual variants.
- Disabled behavior.
- Keyboard interaction.
- Native HTML attributes.

IconButton:

- Accessible label.
- Icon rendering.
- Click interactions.
- Disabled behavior.
- Keyboard interaction.

Alert:

- Message rendering.
- Supported variants.
- Accessible status or alert semantics.

ThemeToggle:

- Light, dark, and system theme selection.
- Active theme representation.
- User interactions.
- Keyboard accessibility.
- Accessible labels.

Expand the matrix where meaningful edge cases are identified.

## Test identifiers

Use stable and descriptive `data-testid` attributes where appropriate.

Examples:

- `theme-toggle`
- `theme-option-light`
- `theme-option-dark`
- `theme-option-system`

Prefer semantic queries such as getByRole when they adequately represent user behavior.

Do not couple tests to CSS classes or internal implementation details.

## Code quality

- Strict TypeScript.
- Avoid any and unnecessary type assertions.
- Descriptive naming and consistent formatting.
- Small, cohesive components.
- Defensive programming and explicit error handling.
- Meaningful comments for non-obvious decisions.
- Accessible and maintainable component APIs.
- No unnecessary dependencies or abstractions.

## Verification

Run the existing frontend quality gates:

- Formatting.
- Linting.
- Type checking.
- Unit tests.
- Coverage.
- Production build.

Report only actual observed results.

## Scope boundaries

Do not implement the calculator business logic, API integration, expression parsing, or history persistence yet.

The future calculator will support seven arithmetic operations, parentheses, complex expressions, calculation history, and keyboard input.

Expression parsing and mathematical evaluation will belong to the Go backend.

Do not modify the backend.

Record this actual prompt in PROMPTS.md.

First, inspect the repository and both reference images.

Present a concise implementation plan containing:

- Design token proposal.
- Component structure and responsibilities.
- Behavioral test matrix.
- Implementation sequence following TDD.
- Any relevant risks or ambiguities.

Wait for approval before modifying files.

Do not commit or push.

Communicate with me in Spanish. Keep source code, comments, tests, and project documentation in English.
````

## Phase 3.1 dependency authorization — literal user prompt

```text
Approved. Add @testing-library/user-event as a development dependency.

Verify compatibility with our pinned toolchain and update the lockfile.

Prefer userEvent for behavioral interactions and semantic queries combined with stable data-testid attributes.

Keep both themes independent of viewport size.

Treat the reference images as visual guidance only, preserving our approved mathematical semantics.

Continue with the concise implementation plan and component test matrix. Record this authorization in PROMPTS.md.
```

## Phase 3.1 implementation authorization — literal user prompt

```text
Implement the plan.
```

### Phase 3.1 execution evidence (assistant notes, not user prompts)

- Installed exact user-event 14.6.13 using Node 22.18.0 and pnpm 12.9.1;
  frozen-lockfile installation passed with sandbox escalation. Initial sandbox
  installation failed on store permissions/network access; escalated restore succeeded.
- Sequential RED: each component test suite failed because its implementation
  module did not exist. GREEN regression totals: Button 9, IconButton 12,
  Alert 18, ThemeToggle 24. Each cycle completed before the next component.
- App integration RED: two new tests could not find theme/preview controls;
  existing shell test passed. GREEN: 26 tests passed across five files.
- Formatting, lint, typecheck, coverage and production build passed. Coverage
  reported 100% of the five exercised App/component files (15 statements,
  13 branches, 10 functions); this does not measure CSS or browser rendering.
- Parent repeated all 26 tests successfully. Real-browser checks remain pending:
  the available computer-use tool reports no browsers. No visual, responsive,
  zoom or operating-system preference validation is claimed from jsdom tests.
- Both themes use CSS color-scheme and light-dark(), requiring modern browser
  support. System preference is resolved by CSS, not a JavaScript listener.

## Phase 3.1 responsive investigation — literal user prompt

```text
Manual browser verification revealed a responsive layout defect.

At a 360x800 mobile viewport, the design system card overflows horizontally and its right edge is clipped.

- Identify and fix the root cause.
- Review container widths, box sizing, padding, intrinsic sizing, and Grid/Flexbox constraints.
- Ensure the layout fits viewports starting at 320px.
- Preserve consistent horizontal spacing.
- Do not hide the issue using overflow-x: hidden.
- Preserve the existing desktop layout and visual identity.
- Run existing frontend quality gates.

Record this prompt in PROMPTS.md.

Do not implement calculator functionality, commit, or push.
```

### Responsive investigation evidence (assistant notes)

- Live Chrome at360x800 had layout width360 and scrollWidth360, but visualViewport
  width324/scale1.111111. Toggling device emulation off/on restored scale1 and
  the complete card edge without any CSS change. This was pinch-zoom clipping.
- Separate actual RED in Chromium: root font20px made body min-width20rem equal
  400px. At320/360 widths scrollWidth was400 and the card ended at380px.
- Removed only the unnecessary body minimum. Scratch browser GREEN at widths
  320/360/1280 with root fonts16/20: zero descendants outside the viewport,
  symmetric gutters, unchanged desktop geometry. At320/font20 the card spans
  20–300; at360/font20 it spans20–340. Desktop16px remains128–1152.
- Live Chrome after the source fix confirmed viewport360/scrollWidth360 with
  root font20px and card20–340. Temporary font override was removed afterward.
- Existing frontend format/lint/typecheck/test/coverage/build passed under
  Node22.18.0 and pnpm12.9.1;26 tests passed. Coverage100% of five exercised
  component files does not cover CSS. Browser probes, not jsdom, established
  the layout correction. No overflow hiding or calculator behavior was added.


## Phase 3.2 — literal user prompt

```text
# Phase 3.2 — Calculator UI

Phase 3.1 is accepted.

Follow AGENTS.md and the approved visual references in docs/design/.

Implement the calculator UI using the existing design system and shared components.

### Components

- CalculatorDisplay
- CalculatorKey
- CalculatorKeypad
- OperationSelector
- Calculator

### Requirements

- Plain CSS, BEM, centralized design tokens, and mobile-first responsive design.
- Support light, dark, and system themes.
- Follow the approved mobile and desktop references closely.
- Include all seven arithmetic operations, parentheses, decimal input, clear, backspace, and equals controls.
- Use controlled components and typed callbacks.
- Include visual states for loading, errors, disabled controls, and selected operations.
- Use accessible labels and stable data-testid attributes.
- Colocate component implementation, CSS, and tests under features/calculator/components/.
- Prevent horizontal overflow and handle long expressions and results gracefully.

### TDD

Define the behavioral test matrix before implementation.

Complete each component individually through genuine RED-GREEN-REFACTOR iterations.

Test rendering, interactions, callbacks, accessibility, disabled states, and edge cases.

Reuse existing shared components rather than duplicating their functionality.

Do not implement mathematical evaluation, expression parsing, API integration, or history persistence yet.

Replace the temporary design system showcase with the calculator interface.

Run existing frontend quality gates and verify the layout in a real browser.

Record this prompt in PROMPTS.md.

Present a concise implementation plan before making changes.

Do not modify the backend, commit, or push.
```

## Phase 3.2 approval — literal user prompt

```text
Yes, implement this plan
```

### Phase 3.2 execution evidence (assistant notes)

- Defined the component behavioral matrix before implementation in
  odd/tasks/phase3-calculator-ui.md. Components completed sequentially.
- Each component first produced a missing-implementation import RED; complete
  GREEN regression counts were Display32, Key35, Keypad48, Selector58,
  Calculator63. No retrospective failures were manufactured.
- App tests against the old showcase produced10 failures/61 passes; the theme
  characterization passed. After textual-editing integration all71 tests passed.
- Node22.18.0/pnpm12.9.1: format/lint/typecheck/test/coverage/build passed.
  Coverage100% of11 exercised files:46statements/37branches/22functions/43lines.
  This does not establish CSS/browser coverage. Parent repeated71tests successfully.
- Live Chrome verified7+2 input, selected Addition, equals integration notice,
  unchanged result placeholder, clear, and theme switching. At360x800 scale1,
  viewport/scrollWidth360 and card16–344; both edges visible.
- Cached Chromium rendered80 cases:320/360/768/1280 widths,16/20px root font,
  light/dark and empty/long/loading/error/disabled states. Long fixtures included
  200-character expressions and100-digit results. No document/descendant overflow;
  calculator controls at least44x44; symmetric gutters and desktop sidebar.
- The browser probe initially assumed media-query rem scaled with an explicit
  CSS root-font override. Correcting the probe expectation resolved10 false
  layout-order failures, without production changes. These were not product RED.
- Rendered fixtures use actual SSR component markup and built CSS; they do not
  substitute for the separately observed hydrated Chrome interactions.
- No math evaluation, parser, API calls, history persistence, backend changes,
  new dependencies, commits or pushes were introduced.

## Calculator visual refinement — literal user prompt

```text
Refine the current calculator UI based on manual browser feedback.

- Reduce excessive vertical spacing, especially in the header and calculator display.
- Make the mobile layout more compact while preserving accessible touch targets.
- Replace oval operation buttons with rounded rectangular pill-style controls.
- Reduce the desktop calculator max-width and improve its overall proportions.
- Soften excessive borders while maintaining accessible contrast.
- Ensure the complete keypad, including equals, is visible without scrolling on typical mobile and desktop viewports whenever space permits.
- Preserve scrolling and accessibility for smaller viewports and enlarged text.
- Maintain the approved visual references, existing design tokens, plain CSS, BEM, and mobile-first architecture.

Do not hide overflow or use fixed viewport heights that clip content.

Preserve existing component contracts and behavior.

Verify both themes in a real browser and run all frontend quality gates.

Record this prompt in PROMPTS.md.

Do not add business logic, dependencies, commit, or push.
```

### Visual refinement evidence (assistant notes)

- Captured real Chromium baseline before CSS edits. At normal text size, Equals
  ended at924.19px on360x800/390x844 and841.98px on1280x800: those height-fit
  assertions failed. Mobile operation borders used50% radius (ovals).
- Changed only CSS/tokens: reduced cumulative spacing and display typography,
  calculator-specific max width58rem,48px minimum keypad height,12px operation
  radius, and a separate subtle decorative-border token. Interactive border,
  focus and text colors remain unchanged; component contracts are untouched.
- Repeated20 browser cases (five viewport sizes, two font sizes, two themes).
  Equals now ends at746.98px on360/390 mobile and686.78px ondesktop; typical
  viewports fit the complete keypad. Desktop card width1024→896px.
- All20 cases: no horizontal overflow, symmetric gutters, calculator controls
  at least44x44px. Small320x568/enlarged text retain natural vertical scrolling.
  No overflow hiding or fixed viewport height was introduced.
- Rendered sampled contrast: text minimum6.42:1 light/7.99:1 dark; interactive
  border-to-exterior4.02:1/4.29:1. Decorative borders are intentionally softer.
- Inspected mobile and desktop screenshots in both themes; parent independently
  inspected mobile-light and desktop-dark. Probes used actual component SSR and
  built CSS in Chromium, not jsdom layout approximations.
- Node22.18.0/pnpm12.9.1: format/lint/typecheck/test/coverage/build passed;
 71tests, coverage100% of11 exercised files (not CSS). Backend and reference
  images unchanged. No business logic, dependency, commit or push changes.


## Phase 3.3 — literal user prompt

```text
# Phase 3.3 — Calculation History UI

Phase 3.2 is accepted.

Follow AGENTS.md, the existing design system, and approved visual references.

Implement the history UI under `features/history/`.

### Scope

- HistoryPanel
- HistoryItem
- HistoryEmptyState
- HistoryToggle

Display expressions, results, and timestamps.

Provide controlled callbacks for reusing expressions, copying results, deleting individual entries, and clearing history.

Use an accessible confirmation interaction before clearing all entries.

Support a responsive, collapsible history panel without compromising the calculator layout.

### Quality

- Plain CSS, BEM, centralized tokens, and mobile-first design.
- Reuse existing shared components.
- Colocate implementation, CSS, and tests.
- Use strict TypeScript, accessible semantics, and stable data-testid attributes.
- Handle empty history, long expressions, long results, and large lists gracefully.

### TDD

Define the behavioral test matrix first.

Complete each component individually through observed RED-GREEN-REFACTOR iterations.

Test rendering, callbacks, empty states, keyboard accessibility, confirmation behavior, and controlled state.

Use fixture data only in tests. Do not display fabricated calculation history in the application.

Do not implement persistence, clipboard integration, mathematical evaluation, or API integration yet.

Run frontend quality gates and real browser verification.

Record this prompt in PROMPTS.md.

Present a concise plan before modifying files.

Do not modify the backend, add dependencies, commit, or push.
```

## Phase 3.3 approval — literal user prompt

```text
Yes, implement this plan
```

### Phase 3.3 execution evidence (assistant notes)

- Defined the behavioral matrix before implementation in
  odd/tasks/phase3-history-ui.md. Completed EmptyState, Item, Toggle and Panel
  individually; each initial RED was an unresolved implementation import before
  test cases executed, not a fabricated behavioral assertion failure.
- GREEN regression totals72/79/83/96; Panel13focused tests passed. App integration
  produced2 actual missing-toggle assertion failures, then13App/98total passed.
- Node22.18.0/pnpm12.9.1: formatting, lint, typecheck,98tests/14files,
  coverage and production build passed. Coverage:92.92% statements,
  98.57% branches,89.36% functions,94.56% lines. Future App callbacks remain
  uncovered because production history is genuinely empty; no artificial
  fixture injection was added to raise coverage. Parent repeated98tests.
- Live Chrome: opening history shows0entries and disabled Clear history;
  native Space collapses the panel and keeps toggle focus. Initial tool key
  spelling Space was unsupported; lowercase space worked (tool issue only).
-60 real Chromium render cases covered320/360/1280 widths,16/20px text,
  light/dark, collapsed App/open empty/long entry/1000entries/invalid timestamp.
  No horizontal overflow, controls>=44x44,1000entries initially render20,
  invalid dates show Time unavailable. Calculator geometry stays unchanged.
-11 hydrated React fixture assertions passed in Chromium:20→40 entries,
  Cancel initial focus, no premature clear, Escape and Cancel preserve data,
  confirmation clears once, empty controlled state restores toggle focus,
  collapse removes content controls and reopening stays empty.
- Hydrated fixture setup needed a library-mode process.env.NODE_ENV definition
  and virtual-time settling instead of animation frames. These were test-harness
  issues, not product RED. Fixtures existed only in isolated test scratch files.
- Production App supplies no copy callback, so copying remains unavailable;
  no clipboard API or fake success feedback. Feature callback behavior is tested.
- README current scope updated; no backend/dependency/reference changes,
  mathematical evaluation, persistence, API integration, commit or push.

## Phase 3 commit and push authorization — literal user prompt

```text
Phase 3 UI implementation is accepted.

Review the final diff for accidental changes, secrets, generated artifacts, and code readability. Run the existing frontend quality gates.

Include the approved design references under docs/design/.

Record this authorization in PROMPTS.md.

I authorize committing and pushing the approved changes to feat/calculator-ui.

Commit message:

feat(ui): implement responsive calculator and history components

Report the commit hash, push result, and PR creation URL.

Do not merge or begin the next phase.
```


## Phase 4 expression engine — literal user prompt

```text
# Phase 4 — Expression Engine

Branch: `feat/expression-engine`.

Follow AGENTS.md and existing backend architectural standards.

Implement mathematical expression evaluation in Go.

## Scope

- Add `POST /api/v1/expressions`.
- Update OpenAPI before implementing the endpoint.
- Implement a lightweight tokenizer, parser, and evaluator under `internal/expression`.
- Reuse the existing calculator domain for mathematical operations.
- Support all seven operations, operator precedence, nested parentheses, unary signs, decimals, right-associative exponentiation, and postfix percentage.
- Preserve the existing calculations endpoint.
- Use the Go standard library only.

## TDD

Define and implement behavioral tests before production code.

Cover at least:

| Expression | Expected |
|---|---|
| 2 + 3 * 4 | 14 |
| (2 + 3) * 4 | 20 |
| 2 ^ 3 ^ 2 | 512 |
| sqrt(9 + 7) | 4 |
| 200 * 15% | 30 |
| -5 + 3 | -2 |
| -2 ^ 2 | -4 |
| (-2) ^ 2 | 4 |
| 2 ^ -3 | 0.125 |
| 10 / (5 - 5) | Division by zero |
| sqrt(-1) | Domain error |
| (2 + 3 | Syntax error |
| 2 + * 3 | Syntax error |

Also test empty expressions, invalid tokens, excessive nesting, expression length limits, non-finite results, and HTTP contract behavior.

Follow genuine RED-GREEN-REFACTOR iterations.

## Quality

Keep the implementation readable, defensive, deterministic, stateless, and concurrency-safe.

Define explicit grammar, precedence, and resource limits.

Use existing RFC 9457 error handling.

Do not use dynamic code evaluation or unnecessary abstractions.

Run existing backend quality gates and verify OpenAPI parity.

Document mathematical assumptions and record the actual prompt in PROMPTS.md.

Complete the scope without additional planning or intermediate progress-only reports unless a genuine blocker occurs.

Do not modify the frontend, commit, or push.
```

## Phase 4 commit and push authorization — literal user prompt

```text
Phase 4 is approved.

Verify the effective Git author and committer email is mtn.gm05@gmail.com.

Review the final diff for secrets, unrelated changes, and generated artifacts.

Record this authorization in PROMPTS.md.

I authorize committing and pushing the approved changes to feat/expression-engine.

Commit message:

feat(api): add mathematical expression evaluation

Report the commit hash, push result, and PR creation URL.

Do not merge or begin the next phase.
```
