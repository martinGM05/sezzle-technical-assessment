# Phase 4 — Expression engine
Branch: feat/expression-engine. Implementation verified; commit and push authorized. No merge or next phase. Runner: Go 1.24.2 native testing.
Route: inline, sequential contract-first TDD. Native reviews remain disabled by user. Forecast: approximately 1,000 authored lines; single approved publication commit.

## Specs
S1: "Add `POST /api/v1/expressions`." "Update OpenAPI before implementing the endpoint."
S2: "Implement a lightweight tokenizer, parser, and evaluator under `internal/expression`." "Reuse the existing calculator domain for mathematical operations."
S3: "Support all seven operations, operator precedence, nested parentheses, unary signs, decimals, right-associative exponentiation, and postfix percentage."
S4: "Preserve the existing calculations endpoint." "Use the Go standard library only."
S5: "Define and implement behavioral tests before production code." "Follow genuine RED-GREEN-REFACTOR iterations."
S6: "Also test empty expressions, invalid tokens, excessive nesting, expression length limits, non-finite results, and HTTP contract behavior."
S7: "Define explicit grammar, precedence, and resource limits." "Use existing RFC 9457 error handling."
S8: "Run existing backend quality gates and verify OpenAPI parity." "Document mathematical assumptions and record the actual prompt in PROMPTS.md."
S9: "Do not modify the frontend, commit, or push."

## Tasks
T1 S1,S7,S8: inline contract and grammar before production — complete; commit authorized.
T2 S2-S6: inline sequential expression behavior TDD — complete; commit authorized.
T3 S1,S4-S8: inline HTTP integration and parity TDD — complete; commit authorized.
T4 S4,S8,S9: inline quality gates and documentation — complete; commit authorized.

## Log
### L1 — Actual user request
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

### L2 — Initial evidence and design
Clean feat/expression-engine at c148ebd4d2f5edae77781dd981d36c846e835600. CodeGraph CLI/index unavailable; inspected local sources directly. Baseline go test -count=1 ./... failed only on sandbox loopback restrictions; same command outside sandbox passed all three packages.
Define a bounded recursive-descent grammar, build a small syntax tree before evaluating through calculator.Calculate, keep evaluation state per request. Length limit 4096 ASCII bytes; recursive grammar nesting 64 (parentheses, sqrt, signs, exponent RHS), linear operator chains bounded by input length. No extra dependency. Native reviews remain disabled. Shared decoding refactor will preserve old endpoint behavior under existing tests.

### L3 — Observed TDD, literals and basic grammar
- Literal tests first: RED (package had no production files), GREEN: eight literal cases pass after minimal ParseFloat implementation.
- Added arithmetic and lexical/syntax matrix: RED with 10 arithmetic failures and 10 lexical/domain failures. Existing passing invalid cases were characterization, not manufactured failures.
- GREEN: tokenizer, syntax tree, sum/product parser and calculator delegation pass all 37 leaf cases. Refactored responsibilities into tokenizer.go/parser.go/expression.go; reran the package successfully.

### L4 — Observed TDD, recursive and advanced grammar
- Grouping/sign/resource tests RED: unsupported parentheses/signs and absent limits; GREEN after recursive unary/grouping parsing and 4096-byte / 64-level guards. Exact-length and flat-chain cases already passed and were characterization.
- Power tests RED: 11 success cases rejected and four errors misclassified; GREEN after right-associative power with unary RHS. Missing exponent was already a syntax error.
- Square-root tests RED: eight success cases rejected, domain and nesting codes wrong; GREEN after sqrt token/primary support. Invalid call forms were characterized.
- Percentage tests RED: 12 valid cases rejected; GREEN after postfix parsing above power precedence. Invalid forms were characterized.
- All package regression runs passed after each increment. Added numeric exact-bit boundaries, error precedence and deterministic finite-result fuzz seeds: all passed on first execution (characterization, no invented RED).

### L5 — Observed HTTP TDD and contract verification
- OpenAPI was written before production. Extended existing endpoint-count expectation for the authorized third path.
- Request-schema tests RED: valid expressions rejected by the local checker's unsupported maxLength keyword. Added standard JSON Schema maxLength handling; all OpenAPI tests GREEN.
- HTTP success/validation tests RED: all 40 cases returned 404 before route implementation. GREEN after endpoint, expression DTO decoding, and extraction of the existing bounded object decoder. Full backend regression passed outside sandbox immediately after integration.
- Added HTTP method/media/body limits, read failures, duplicate/escaped fields, restrictive CORS/preflight and 32 concurrent requests. All passed on their first execution (characterization of reused HTTP transport behavior).
- No modifications to calculator arithmetic, CORS implementation, lifecycle, frontend, dependencies or CI. Risk: item 3 (shared HTTP decoder extraction must preserve existing calculation responses).


### L6 — Final observed verification
Go 1.24.2 on darwin/arm64. Native reviews remained disabled; no review lifecycle was invoked.
- `go test -count=1 -race -covermode=atomic -coverprofile=coverage/coverage.out -json ./...`: PASS, 632 leaf executions (160 calculator, 133 expression including nine fuzz seeds, 321 HTTP, 18 lifecycle). Parent test containers are excluded from that count.
- Coverage: 98.9% total statements; calculator 100%, expression 99.3%, HTTP 100%, cmd/api 90.3%. Fresh ignored profile, summary and HTML under backend/coverage/.
- `go vet ./...`, `go build ./...`, `gofmt` check and `git diff --check`: PASS.
- Explicit OpenAPI/HEAD/contract-checker suite: PASS. Local schema-subset validation is not a full third-party OpenAPI validator.
- `go test ./internal/expression -run '^$' -fuzz '^FuzzEvaluate$' -fuzztime=5s -parallel=2`: PASS, 118755 executions, no failure found. This is bounded fuzz evidence, not proof for all possible inputs.
- Independent functional verifier (risk item 3, not native review) compared baseline c148ebd4d2f5edae77781dd981d36c846e835600 with candidate in isolated scratch: 53 old-route scenarios identical in status, headers and exact body; 28 additional expression HTTP boundary cases passed. Its three-package suite passed and the parent reran that exact command successfully.
- Sandbox denied loopback binding during initial tests; unrestricted approved reruns passed. No test was skipped to hide that restriction.
- No frontend changes, dependency additions, commit or push. Original shared decoder semantics (including last-value duplicate JSON keys) preserved. ASCII-only expressions and float64 limitations are intentional and documented.

### L7 — Publication authorization (literal user prompt)

Phase 4 is approved.

Verify the effective Git author and committer email is mtn.gm05@gmail.com.

Review the final diff for secrets, unrelated changes, and generated artifacts.

Record this authorization in PROMPTS.md.

I authorize committing and pushing the approved changes to feat/expression-engine.

Commit message:

feat(api): add mathematical expression evaluation

Report the commit hash, push result, and PR creation URL.

Do not merge or begin the next phase.

Delivery checks: effective author and committer email verified as mtn.gm05@gmail.com; final diff contains only approved Phase 4 files. No detectable credential patterns, local metadata, generated artifacts, or frontend changes. Earlier no-commit restrictions above are historical implementation-stage instructions.
