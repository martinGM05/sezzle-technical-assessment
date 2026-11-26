# Phase 5 — Functional integration
Branch: feat/calculator-integration. Single-commit delivery authorized in L6. Native reviews remain disabled.
Runner: Node 22.18.0 / pnpm 12.9.1, Vitest/RTL/userEvent; real Go backend. Approximately 1,200 authored lines forecast; no publication authorized.

## Specs
S1: "Implement a typed API client using native fetch." "Use POST /api/v1/expressions for expression evaluation." "Normalize visual mathematical symbols into the supported ASCII syntax."
S2: "Implement calculator state and interactions through focused custom hooks." "Support all seven operations, parentheses, decimals, clear, backspace, and equals." "Support physical keyboard input."
S3: "Handle loading, RFC 9457 errors, network failures, cancellation, and stale responses." "Prevent duplicate submissions."
S4: "Connect history actions: reuse, copy, delete, and clear." "Persist successful calculations and theme preferences using localStorage." "Handle invalid or unavailable localStorage safely."
S5: "Configure local API communication without hardcoding production URLs."
S6: "Define behavioral tests before implementing each responsibility." "Use genuine RED-GREEN-REFACTOR iterations and the existing testing tools."
S7: "Preserve the approved UI, CSS architecture, accessibility, and stable data-testid attributes." "Do not duplicate mathematical evaluation in the frontend."
S8: "Do not introduce unnecessary dependencies or modify the backend without identifying a concrete defect." "Run all frontend quality gates and verify integration against the real Go backend." "Record actual prompts and test results in PROMPTS.md." "Do not commit or push."

## Tasks
T1 S1,S3,S6,S7: delegated API/normalization behavior with its tests — complete; no commit.
T2 S4,S6,S7: delegated storage/history/theme behavior with its tests — complete; no commit.
T3 S2-S7: inline calculator state, keyboard, clipboard, UI integration and local proxy — complete; no commit.
T4 S8: inline quality gates, real-backend/browser verification, docs — complete; no commit.

## Log
### L1 — Actual user request
# Phase 5 — Functional Integration

Branch: `feat/calculator-integration`.

Follow AGENTS.md, the existing component architecture, and the OpenAPI contract.

## Scope

Connect the existing calculator UI to the Go backend.

- Implement a typed API client using native fetch.
- Use POST /api/v1/expressions for expression evaluation.
- Normalize visual mathematical symbols into the supported ASCII syntax.
- Implement calculator state and interactions through focused custom hooks.
- Support all seven operations, parentheses, decimals, clear, backspace, and equals.
- Support physical keyboard input.
- Handle loading, RFC 9457 errors, network failures, cancellation, and stale responses.
- Prevent duplicate submissions.
- Connect history actions: reuse, copy, delete, and clear.
- Persist successful calculations and theme preferences using localStorage.
- Handle invalid or unavailable localStorage safely.
- Configure local API communication without hardcoding production URLs.

## TDD

Define behavioral tests before implementing each responsibility.

Cover API communication, symbol normalization, calculator interactions, keyboard input, error recovery, asynchronous state, history, and persistence.

Use genuine RED-GREEN-REFACTOR iterations and the existing testing tools.

## Constraints

Preserve the approved UI, CSS architecture, accessibility, and stable data-testid attributes.

Do not duplicate mathematical evaluation in the frontend.

Do not introduce unnecessary dependencies or modify the backend without identifying a concrete defect.

Run all frontend quality gates and verify integration against the real Go backend.

Record actual prompts and test results in PROMPTS.md.

Proceed directly with implementation without additional planning.

Do not commit or push.

### L2 — Evidence, route and responsibility boundaries
Baseline clean feat/calculator-integration at 893c6404febe4e74069e14e7a00f4be0eef185cf; 98 existing tests pass. CodeGraph CLI/index unavailable; direct source inspection used. Two independent delegated behavior units have disjoint surfaces and sufficient work: native API/normalization and storage/history/theme. Parent owns calculator hook, keyboard, clipboard, App, local proxy, docs/log and integration. No concurrent writes to shared files. No backend changes planned.

Internal seams agreed before writes: evaluateExpression(expression: string, signal?: AbortSignal): Promise<number> uses relative /api/v1/expressions; normalizeExpression(string): string maps only visual tokens, no parsing/evaluation. useHistory exposes entries, add(expression: string, result: string), remove(id), clear(), storageError. useTheme exposes theme, setTheme, storageError and applies data-theme. History key sezzle.history.v1 and theme key sezzle.theme.v1; bounded validated versioned JSON and safe failures. Hooks receive only confirmed server successes, no fixtures in application.

Interaction policy: edits clear prior result/error and invalidate in-flight work; equals while pending is ignored synchronously. Clear/reuse/unmount cancels and stale responses never update state/history. Repeated equals after success is allowed as a new explicit calculation. Square-root control inserts √(, backend owns syntax validation. Keyboard handles digits/operators/parentheses/decimal, Enter/=, Backspace, Escape/Delete clear, and r for sqrt; native control activation/shortcuts/editable fields are not hijacked. Relative fetch and configurable development proxy keep production routing same-origin.


### L3 — Completed implementation and verification
T1 and T2 completed through the disjoint delegated responsibility boundaries.
T3 completed inline with observed RED/GREEN for calculator state, keyboard,
clipboard, App composition, clear-during-loading and proxy configuration.
The cross-realm AbortError assertion exposed an instanceof limitation; the hook
now recognizes the error name defensively. One later stale-success regression
passed immediately and was recorded honestly as characterization.

T4 complete: all six frontend gates pass with 235 tests in 24 files. Coverage:
99.66% statements, 98% branches, 100% functions/lines. Actual test evidence is in
PROMPTS.md. Real production client/proxy/Go verification passed 12 successful
expressions and five Problem Details responses. Chrome verified calculations,
error recovery, clipboard and reload persistence in light/dark themes.
Independent isolated functional probes passed without product findings. No native
review was started; the user-disabled switch remains respected.

README and architecture now document integrated startup, restrictive origin
configuration, keyboard semantics, cancellation, plaintext versioned local
storage, the 100-entry history bound and secure-context clipboard limitations.
No new dependency, backend or CSS change; no commit or push. Existing Phase 2
metrics are explicitly labeled historical rather than presented as current.

Risk: items 1 (history deletion/persistence), 4 (asynchronous response ordering)
and 5 (local proxy configuration); behavioral and isolated regression checks plus
real-backend verification passed. Remaining limitations: no cross-tab history
synchronization, mobile verification or automated full-stack browser suite in
this phase; no claims about remote CI execution.


### L4 — Manual keyboard regression correction
User reported Enter clearing `(2+3)×4` after calculator-button interaction; the
complete literal correction prompt and actual results are in PROMPTS.md under
"Phase 5 keyboard regression". This intentionally corrects L2's earlier policy:
inside the calculator Enter always evaluates, while Space activates the focused
button. External controls and editable fields retain native keyboard behavior.
Browser reproduction confirmed retained focus plus native Enter click as cause.
Observed RED: 8 failures/38 passes; GREEN: 46 focused passes. Full gates pass
with 244 tests; coverage 99.66% statements, 97.95% branches, 100% functions/lines.
Real Chrome/backend returned20 with Clear,Digit4,Backspace focus. Independent
read-only focused verification passed; no native review, commit or push.
Risk: item3 (shared keyboard hook); focused interaction/accessibility regression
checks and real-backend verification passed. Prior Phase5 work remains intact.


### L5 — Roadmap scope reduction (actual user prompt)
Remove Playwright, k6, Lighthouse, and additional performance testing from the remaining implementation roadmap.

These are optional engineering improvements, not assessment requirements.

Complete the current keyboard regression fix using the existing testing tools and manual browser verification.

Update the roadmap accordingly without altering historical prompts.

Do not introduce additional dependencies or expand the project scope.

Applied to README.md and docs/REQUIREMENTS.md: these optional tools and additional
performance work are removed from remaining deliverables, not deferred. Existing
quality gates remain. L4's completed keyboard regression and observed 244-test
verification stand; no source, dependency or CI change was needed for this update.
Historical prompts are preserved byte-for-byte; this request is appended only.


### L6 — Delivery authorization (actual user prompt)
Phase 5 is approved.

Verify the final diff and personal Git identity. Record this authorization in PROMPTS.md.

I authorize committing and pushing the approved changes to feat/calculator-integration.

Commit message:

feat(calculator): integrate expression API and persistent history

Report the commit hash, push result, and PR creation URL.

Do not merge or begin the next phase.

Final delivery checks: effective Git author and committer both use
mtn.gm05@gmail.com. Scope is the approved Phase 5 integration, Enter correction
and roadmap reduction. Earlier no-publication instructions describe their
original implementation stages; this authorization permits one coherent commit
with the exact requested message and push to origin/feat/calculator-integration.
No PR creation or merge is included. Generated artifacts and local metadata stay
ignored. The commit containing this log is the delivery boundary for T1–T4.
