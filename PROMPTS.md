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
