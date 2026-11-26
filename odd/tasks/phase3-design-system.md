# Phase 3.1 — Design system
Branch: feat/calculator-ui. Implementation verified; user authorized the Phase 3 commit and push. No merge or next phase.
Runner: pinned Node 22.18.0 / pnpm 12.9.1, Vitest and React Testing Library.

## Specs
S1: "Use plain CSS exclusively." "Support light, dark, and system themes."
S2: "Each component must colocate its implementation, styles, and tests."
S3: "Complete each component individually before proceeding to the next."
S4: "Never fabricate TDD evidence or test results."
S5: "Do not modify the backend." "Do not modify or delete the reference images."
S6: "Approved. Add @testing-library/user-event as a development dependency."
S7: "Verify compatibility with our pinned toolchain and update the lockfile."
S8: "Prefer userEvent for behavioral interactions and semantic queries combined with stable data-testid attributes."
S9: "Keep both themes independent of viewport size."
S10: "Treat the reference images as visual guidance only, preserving our approved mathematical semantics."

## Tasks
T1 S6-S7 dependency and frozen install — delegated, complete, commit authorized.
T2 S1-S4,S8 Button → IconButton → Alert → ThemeToggle TDD — delegated, complete, commit authorized.
T3 S1,S9 shell and themes — delegated, complete, commit authorized.
T4 S4-S5 all gates and independent verification — complete, commit authorized.


## Log
L1: "Implement the plan."
The approved implementation plan and full literal Phase 3.1 request are retained in the conversation; PROMPTS.md will hold the literal requests. Parent owns prompts and task log. Delegated implementation reason: long-session context backstop; one worker owns frontend and completes components sequentially, not parallel component implementation. Estimated 700–1000 authored lines including tests; no commit slicing because commits are forbidden.

### L2 — observed implementation and checks
T1 complete: exact user-event14.6.13; pinned runtime and frozen installation passed with escalation.
T2 complete: missing-module RED for each component; sequential GREEN totals9/12/18/24. App two missing-control RED cases →26 tests GREEN.
T3 complete: plain CSS tokens/BEM, light-dark()+color-scheme; both themes any viewport, controlled radio group, system default, no persistence.
T4 partial: format/lint/typecheck/test/coverage/build pass; parent26tests pass. Coverage100% of exercised5files only, no CSS/browser proof. Browser tool has no browsers; visual/responsive/zoom/reduced-motion checks pending. Independent audit pending.
Risk: item5 (test dependency and lockfile change). No commits or push.

### L3 — independent verification
Independent scratch-copy checks: Prettier/ESLint/TypeScript/Vitest26tests/coverage/build passed under Node22.18.0. No reproducible blockers. Dependency-only lockfile additions verified. Text contrast minima6.42:1 light/6.74:1 dark; border minima3.57:1/3.41:1. Reference hashes unchanged and backend untouched. Native review preflight requested untracked inventory selection; no review START/consent/approval occurred. Visual checks remain pending because CUA lists no browsers. Implementation complete; acceptance pending browser proof.


### Delivery authorization
Phase 3 accepted. Final diff audit found no blocking issues; all six frontend gates passed with98tests. User explicitly authorized commit and push with message `feat(ui): implement responsive calculator and history components`. Historical log restrictions above describe the earlier implementation stages. No merge or next phase authorized.
