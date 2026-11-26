# Phase 3.2 — Calculator UI
Branch: feat/calculator-ui. Implementation verified; user authorized the Phase 3 commit and push. No merge or next phase.
Runner: Node22.18.0 / pnpm12.9.1, Vitest/RTL/userEvent. Parent owns prompts/log; one sequential frontend writer.

## Specs
S1: "Implement the calculator UI using the existing design system and shared components."
S2: "Plain CSS, BEM, centralized design tokens, and mobile-first responsive design."
S3: "Support light, dark, and system themes."
S4: "Follow the approved mobile and desktop references closely."
S5: "Include all seven arithmetic operations, parentheses, decimal input, clear, backspace, and equals controls."
S6: "Use controlled components and typed callbacks."
S7: "Include visual states for loading, errors, disabled controls, and selected operations."
S8: "Use accessible labels and stable data-testid attributes."
S9: "Colocate component implementation, CSS, and tests under features/calculator/components/."
S10: "Prevent horizontal overflow and handle long expressions and results gracefully."
S11: "Define the behavioral test matrix before implementation."
S12: "Complete each component individually through genuine RED-GREEN-REFACTOR iterations."
S13: "Test rendering, interactions, callbacks, accessibility, disabled states, and edge cases."
S14: "Reuse existing shared components rather than duplicating their functionality."
S15: "Do not implement mathematical evaluation, expression parsing, API integration, or history persistence yet."
S16: "Replace the temporary design system showcase with the calculator interface."
S17: "Run existing frontend quality gates and verify the layout in a real browser."
S18: "Record this prompt in PROMPTS.md."
S19: "Do not modify the backend, commit, or push."

## Tasks
T1 S1-S14 Display→Key→Keypad→OperationSelector→Calculator sequential TDD — delegated, complete.
T2 S3,S6,S15-S16 App textual editing/integration — delegated, complete.
T3 S10,S17 browser and quality verification — parent, complete.
T4 S18-S19 prompts/scope/checks — parent, complete.

## Behavioral matrix
Display: empty/zero/long expression and result, loading, errors, accessible names.
Key: native button semantics, label, action payload, click/Enter/Space, disabled.
Keypad: digits/decimal/parentheses/seven ops/clear/backspace/equals; callbacks and disabled controls/focus order.
OperationSelector: seven operations, accessible selection, controlled updates, callbacks/disabled.
Calculator: composition and state propagation, busy/disabled/error, exact callbacks, long content.
App: theme preservation, append-only text input, clear/backspace including empty, selected operation, equals announces integration pending without fabricated result or simulated loading.


## Log
L1 — literal user request:

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

L2 — literal approval: "Yes, implement this plan"

Approved plan: mobile display/compact operations/four-column keypad; desktop display+keypad left and operations right. Display→Key→Keypad→OperationSelector→Calculator sequential RED/GREEN. App owns text editing only, selected operation and themes; equals announces integration pending. No fake results/loading. Delegation reason: long-session context backstop, one owner preserving sequential TDD. Forecast about900–1300 authored lines including tests; no commit slicing authorized. Preserve existing responsive body-minimum fix.

### L3 — implementation and verification
T1 complete: all five components sequential RED/GREEN, regression32/35/48/58/63.
T2 complete: App10RED/61pass againstshowcase→71GREEN. Existingthemes preserved; onlytextediting plus integrationnotice.
T3 complete: allfrontendgates pass Node22.18/pnpm12.9.1; parent71tests pass. Coverage100%11exercisedfiles, notCSSproof. RealChrome interactions7+2/equals/clear/theme;360viewportscroll360card16–344. CachedChromium80layoutcases320/360/768/1280×font16/20×light/dark×5states allfit; controls>=44px. Probe media-query rem expectation corrected, no fabricatedproductfailure.
T4 complete: literalprompts recorded, backend/images unchanged, nocommit/push. Native review assessment requires untracked scope selection; no verdict claimed.
Risk: none for Phase3.2 incremental UI; shared component contracts unchanged. Phase3.1 dependency change remains separate historical item5.

### L4 — final review applicability
Explicit intended-untracked assessment (frontend source and task docs, excluding reference images): native risk medium, review_due false, reason under_budget. No native review START/consent/verdict was required by this assessment. This is applicability evidence, not an approval receipt. No commits/push.

### L5 — approved visual refinement
User request recorded verbatim in PROMPTS.md under Calculator visual refinement. Existing component contracts preserved; CSS-only inline refinement. Real-browser baseline RED Equalsbottom924.19mobile/841.98desktop, GREEN746.98/686.78. Header/display/gaps compacted,58rem maxshell, operationradius12px, subtledecorativeborder.20cases light/dark widths320/360/390/1280/1440 font16/20 retain44pxtargets andzeroXoverflow. Small/enlarged layouts scrollnaturally. Sixqualitygates pass71tests. Independent renderer verified andparentviewed screenshots. No commits/push; no sharedbehavior change. Risk: item3 (shared CSS tokens; independent rendering and all regressions passed).


### Delivery authorization
Phase 3 accepted. Final diff audit found no blocking issues; all six frontend gates passed with98tests. User explicitly authorized commit and push with message `feat(ui): implement responsive calculator and history components`. Historical log restrictions above describe the earlier implementation stages. No merge or next phase authorized.
