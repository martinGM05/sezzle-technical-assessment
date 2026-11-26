# AI-Assisted Development — Curated Technical Record

## Editorial disclosure and original archive

This is an **editorially curated record**, not a verbatim transcript or a claim of
uninterrupted TDD compliance. It consolidates repeated requests and omits personal
identifiers, communication preferences and repetitive operational instructions.
Quoted blocks are literal technical excerpts from user prompts; surrounding prose
is an editorial summary. Execution evidence is attributed to assistant notes, not
presented as user-authored prompts. Earlier failures and deviations remain below.

The complete pre-curation file is preserved byte-for-byte locally at
`.submission-private/PROMPTS.original.md` (85,100 bytes). Its SHA-256 is:

```text
76c47f65db5c9bae9116e0015caa65ef18ea444fbfd18b7860ed558704af9c59
```

That private archive and the latest curation request are Git-ignored and excluded
from submission because they retain unredacted material. Preserve this local copy;
a fresh clone will not contain it. Prior committed prompt history is unchanged,
not scrubbed or rewritten. Original section names below locate the corresponding
archive entries; a reviewer-facing copy can be separately shared after privacy
review if needed.

The archive itself disclosed earlier restoration of abbreviated prompts and a
combined entry, plus a duplicate represented by reference. Restoring wording did
not establish contemporaneous compliance. Neither the archive nor this edition
is asserted to be a complete capture of every interaction.

## 1. Assessment framing and engineering foundation

**Prompt excerpts — “Assessment brief and Phase 1 request”:**

> Prioritize correctness, simplicity, clarity, maintainability, and engineering judgment over unnecessary features or architectural complexity.

> Distinguish these assessment requirements from our proposed engineering improvements. Verify them against the original assignment when provided.

> Do not implement business functionality before its corresponding tests.

**Technical decisions:** React/TypeScript/Vite frontend; independent Go module
using the standard library; stateless synchronous REST; domain/HTTP separation.
No database, broker, gateway or distributed infrastructure. General-purpose
float64 arithmetic is not monetary arithmetic. Mandatory reported requirements
are distinguished from optional arithmetic, Docker and proposed engineering work.
The original assessment brief remains independently unavailable.

**Corrections — Phase 1 adjustments/review:** strict TypeScript had to be explicitly
enabled; Node 22/jsdom compatibility and Node type alignment were corrected.
Toolchain versions and pnpm lockfile were pinned; CI uses frozen installation,
independent jobs, minimal permissions and verified immutable Action SHAs. Empty
backend scaffolding was not given artificial tests or coverage. Setup preceded
business implementation.

## 2. Calculator API and recovery of TDD discipline

**Prompt excerpts — “Phase 2 recovery request”:**

> OpenAPI was not defined before implementation. Document this deviation and establish the contract before implementing further behavior.

> Only one TDD cycle has been reported. Continue with genuine RED-GREEN-REFACTOR iterations for all remaining behaviors.

> Do not fabricate retrospective TDD evidence.

**Prompt excerpt — “Recovery clarification”:**

> For existing behavior, add characterization tests and report their actual results. Do not artificially force RED failures.

**Decisions:** OpenAPI was established before subsequent new behavior. A single
calculation route uses operation-specific arity: two operands for add/subtract/
multiply/divide/power; one for square root/percentage. Percentage is operand/100;
finite zero-to-zero is 1. Reject non-finite operands/results and overflow; document
underflow and comparison tolerances. Strict JSON handling and RFC 9457 errors
remain separate from mathematical rules.

**Prompt excerpt — “Domain test matrix review correction”:**

> It contains only 15 scenarios and does not satisfy the agreed domain test matrix.

**Historical assistant execution evidence:** the original 15-case suite actually
failed, and earlier passing claims/reuse of an old 60% coverage measurement were
not valid evidence. Expanding the matrix produced 150 passing characterization
cases and seven failures. Unary arity and non-finite validation were corrected
in observed RED/GREEN increments; all 157 domain cases then passed. The later
refactor preserved behavior rather than retroactively inventing RED runs.

**HTTP evidence:** the expanded 94-case run initially had 76 passing
characterization cases and 18 failures in media types, strict structure, health
and routing. Separate body-limit tests exposed 400 instead of 413; panic recovery
and restrictive CORS also received actual failing tests before corrections.
Concurrent/parity tests passed as characterization. OpenAPI is stored as JSON
syntax, a YAML 1.2 subset, allowing dependency-free Go contract checks; that
checker is intentionally limited, not a full third-party OpenAPI validator.

Lifecycle work covered SIGINT/SIGTERM, bounded shutdown, nonzero startup failure
and expected ErrServerClosed. Same-origin proxy tests exposed Origin forwarding:
explicit ALLOWED_ORIGINS configuration fixed the deployment issue without a
wildcard or trusting forwarded headers. HEAD documentation was corrected.
Duplicate JSON member names retaining the last value remain a documented
limitation, not falsely described as rejected. Historical review details remain
in [docs/REVIEW.md](docs/REVIEW.md).

## 3. Design system, calculator and history UI

**Prompt excerpts — “Phase 3.1” and responsive investigation:**

> Use plain CSS exclusively.

> Keep shared components independent of calculator-specific business logic.

> Do not hide the issue using overflow-x: hidden.

**Decisions:** mobile-first CSS, BEM and centralized tokens; light/dark/system
independent of viewport size; accessible reusable Button, IconButton, Alert and
ThemeToggle; controlled calculator/history components with colocated tests.
user-event was explicitly approved for behavioral interactions. UI-only phases
did not evaluate expressions or display fabricated calculation history.

**Historical evidence:** components were implemented sequentially after observed
missing-module failures; those failures are not misrepresented as mathematical
behavior failures. Initial 100% coverage applied only to five exercised component
files, not CSS or real browser rendering.

**Responsive correction:** the reported 360px clipping first reproduced as pinch
zoom, while layout/scroll widths were both 360. A separate real defect occurred
with a 20px root font: body min-width:20rem forced 400px width. Removing that
minimum fixed 320/360px layouts with enlarged text while preserving desktop
geometry. Real browser probes, not jsdom coverage, established that correction.
Later compactness refinement preserved touch targets, scrolling and accessibility.
History UI used controlled callbacks and confirmation before clearing entries.

## 4. Backend expression engine

**Prompt excerpts — “Phase 4 expression engine”:**

> Reuse the existing calculator domain for mathematical operations.

> Support all seven operations, operator precedence, nested parentheses, unary signs, decimals, right-associative exponentiation, and postfix percentage.

> Define explicit grammar, precedence, and resource limits.

**Acceptance examples from that request:** `2 + 3 * 4` → 14, `(2 + 3) * 4` → 20,
`2 ^ 3 ^ 2` → 512, `sqrt(9 + 7)` → 4, `200 * 15%` → 30, `-2 ^ 2` → -4,
`(-2) ^ 2` → 4, and `2 ^ -3` → 0.125, plus syntax/domain/limit errors.

The new expression endpoint preserved the calculation endpoint and used a
standard-library tokenizer, parser and evaluator with backend-owned arithmetic.
**Evidence limitation:** the archive has no dedicated Phase 4 execution-evidence
section. Later approvals and passing suites do not reconstruct historical
RED/GREEN cycles. Current behavior is covered by the reproducible backend suite
and [documented grammar](docs/EXPRESSIONS.md).

## 5. Functional integration and the Enter regression

**Prompt excerpts — “Phase 5 functional integration”:**

> Implement a typed API client using native fetch.

> Handle loading, RFC 9457 errors, network failures, cancellation, and stale responses.

> Do not duplicate mathematical evaluation in the frontend.

**Decisions:** relative expression API requests, visual-symbol normalization only,
focused hooks, synchronous duplicate guard plus request identity/AbortController,
success-only history, guarded versioned localStorage and explicit clipboard
feedback. Invalid/unavailable storage retains usable in-memory interaction.

**Historical evidence:** responsibilities received tests before implementations.
Missing-module failures, actual behavioral failures and passing characterization
were reported separately. A cross-realm AbortError required a defensive name
check. A later stale-success case passed immediately and was characterization.
Real client/proxy/Go tests and manual Chrome checks complemented mocks.

**Prompt excerpt — “Phase 5 keyboard regression”:**

> Expected behavior: pressing Enter must evaluate the expression and display 20, consistently with clicking equals.

**Observed correction:** retained focus on Clear caused native Enter activation
to clear the expression; focus on Digit 4 appended another 4. The keyboard hook
exempted all button Enter events. Before the fix, regression tests produced
**8 failures / 38 passes**. Removing the exemption and preventing default before
repeat suppression produced **46 focused passes**. Eight focus cases assert the
unchanged expression, result 20, exactly one request and one history entry.
Space and controls outside the calculator retain native behavior; propagation
is not stopped. Real Chrome/Go verification confirmed the fix. An incorrect test
working directory and unsupported RTL option were test/tooling errors, not
additional product RED cycles.

**Scope correction — “Roadmap scope reduction”:** Playwright, k6, Lighthouse and
additional performance testing were explicitly removed from remaining work as
optional improvements, not assessment requirements.

## 6. Final container delivery

**Prompt excerpts — “Phase 6 final delivery”:**

> Add a multi-stage frontend Dockerfile with a non-root Nginx runtime and .dockerignore.

> Allow the complete application to start with `docker compose up --build`.

> Generate fresh coverage reports and make them accessible to reviewers.

**Decisions:** verified digest-pinned frontend builder/Nginx runtime; unchanged
non-root Go image; frozen dependencies; same-origin proxy preserving Origin;
internal-only backend, loopback frontend and combined health verification.
No new application features or testing dependencies were required.

**Observed deployment RED/GREEN:** the initial internal-only frontend network
passed container health but did not publish host port 8088. A native Node
published-health regression failed both normally and outside the sandbox. Adding
a frontend web bridge alongside the internal network corrected it; the same
check and all **27 real HTTP smoke cases** passed. Backend stayed internal-only.
Environment socket/listener permission failures were not counted as product RED.

Manual Chrome verified seven operations, complex expressions, Enter, error
recovery, reuse/copy and history/theme reload persistence. A controlled backend
outage produced unhealthy readiness and a safe UI error; restart restored health
and evaluation. An already-unhealthy container briefly caused Compose --wait to
fail before its next successful probe; this was recorded, not concealed.

Historical Phase 6 evidence: **244 frontend tests**, **632 backend leaf cases**
(distinct from 684 parent-inclusive pass events), 99.66% frontend statement
coverage and 98.9% backend coverage. No Phase 6 remote CI, image signature audit,
vulnerability scan, non-arm64 execution or packaging-phase mobile emulation is
claimed. Manual destructive history deletion was not claimed; automated
application tests cover deletion and confirmation.

## 7. Submission preparation

**Prompt excerpt — final submission request:**

> Clearly disclose the editorial curation and preserve the original prompt archive.

> Preserve technical decisions, corrections, deviations, and genuine TDD evidence.

The current edition follows that request; it does not erase or repair the
historical deviations above. The README focuses on starting and evaluating the
application and reproducing coverage. The obsolete internal arithmetic footer
was already absent in the user-edited source before this preparation; a passing
characterization test verifies that condition without manufacturing RED evidence.
The existing author credit is preserved. Current verification results and report
locations are maintained in [README.md](README.md), separate from historical
milestone counts. No dependency or Git history changes are part of this curation.


**Submission checks (assistant execution notes):** all frontend gates passed with
245 tests and fresh coverage (99.66% statements, 97.95% branches, 100% functions/
lines). The rebuilt Compose application passed 27 HTTP checks. Reloading Chrome
confirmed the obsolete arithmetic footer is absent; an already-open page had
retained the previous bundle until reload. The pre-existing author credit was
not rewritten. Backend sources and their Phase 6 coverage remain unchanged.
A scan of 134 tracked/nonignored files found no credential patterns or accidentally
included generated/private artifacts. The auxiliary Go JSON counting log was
removed; reproducible HTML/profile/summary coverage reports were retained.
