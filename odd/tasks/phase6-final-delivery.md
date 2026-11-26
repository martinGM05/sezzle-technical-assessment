# Phase 6 — Final delivery
Branch: chore/final-delivery. Baseline c2a60f58d370b19ef7a29cc51ac7f91164f87e77. No commits/push. Native reviews disabled.
Runner: existing pnpm/Go gates, Docker Compose smoke checks and manual Chrome. Inline implementation; independent delivery/security verification. Approximately 350 new authored lines forecast, excluding logs.

## Specs
S1: Add a multi-stage frontend Dockerfile with a non-root Nginx runtime and .dockerignore.
S2: Pin container base images to verified versions or digests.
S3: Configure Nginx to serve React and proxy API requests to the existing Go container.
S4: Add Docker Compose for both services with appropriate networking, origin configuration, and health verification.
S5: Allow the complete application to start with `docker compose up --build`.
S6: Verify the real containerized application, all seven operations, complex expressions, keyboard input, history, persistence, and error handling.
S7: Run existing frontend and backend quality gates.
S8: Generate fresh coverage reports and make them accessible to reviewers.
S9: Finalize README with setup instructions, Docker usage, architecture, API examples, testing, coverage results, and known limitations.
S10: Ensure all original assessment requirements are satisfied.
S11: Record actual prompts and verification results in PROMPTS.md.
S12: "Keep the implementation minimal. Use TDD for any behavioral fixes."
S13: "Do not introduce Playwright, k6, Lighthouse, additional performance tooling, or new features."
S14: "Do not commit or push." Git history must remain unchanged; subsequent submission authorization in L4 permits editorial curation with an intact original archive.

S15: Submission curation, intact original archive, coverage accessibility, concise README, obsolete-footer verification and temporary-file hygiene, as quoted in L4.

## Tasks
T1 S1-S5,S12-S14: inline frontend packaging and two-service Compose — complete; no commit.
T2 S6-S8,S12-S14: existing frontend/backend quality gates, coverage and actual containers/browser verification; independent verifier for item5 delivery and item2 container isolation — complete; no commit.
T3 S9-S11,S14: inline final README, requirement traceability, prompt/results — complete; no commit.

T4 S15: inline submission preparation with read-only archive mapping — complete; no commit.

## Log
### L1 — Actual request
# Phase 6 — Final Delivery

Branch: `chore/final-delivery`.

Follow AGENTS.md and existing project decisions.

Complete the final deployment and documentation.

- Add a multi-stage frontend Dockerfile with a non-root Nginx runtime and .dockerignore.
- Pin container base images to verified versions or digests.
- Configure Nginx to serve React and proxy API requests to the existing Go container.
- Add Docker Compose for both services with appropriate networking, origin configuration, and health verification.
- Allow the complete application to start with `docker compose up --build`.
- Verify the real containerized application, all seven operations, complex expressions, keyboard input, history, persistence, and error handling.
- Run existing frontend and backend quality gates.
- Generate fresh coverage reports and make them accessible to reviewers.
- Finalize README with setup instructions, Docker usage, architecture, API examples, testing, coverage results, and known limitations.
- Ensure all original assessment requirements are satisfied.
- Record actual prompts and verification results in PROMPTS.md.

Keep the implementation minimal. Use TDD for any behavioral fixes.

Do not introduce Playwright, k6, Lighthouse, additional performance tooling, or new features.

Do not rewrite Git history or curate historical prompts yet.

Complete the scope without additional planning.

Do not commit or push.

### L2 — Actual baseline authorization
Approved. Verify that c2a60f5 contains the merged Phase 5 changes and update chore/final-delivery using fast-forward only.

Preserve the existing Git history and working tree.

Continue Phase 6 immediately after verifying the correct baseline. No additional planning is required.

Do not commit or push.

Verified clean branch initially at Phase4 merge893c640; Phase5 integration, persistence and Enterfix present in approvedc2a60f5 and origin/feat/calculator-integration. Authorized git merge --ff-only succeeded to existingcommit; no newcommit/historyrewrite. Local origin/main still points to893c640, so no claim of remotely verified Phase5 PR merge. Docker27.4.0/Compose2.31.0 available with sandbox approval. Original assessment brief remains unavailable independently; final traceability will distinguish reported requirements from engineering improvements.


### L3 — Final evidence
T1: verified digest-pinned Node/Nginx frontend build plus unchanged distroless Go,
non-root/read-only/capability restrictions, explicit origins and two-service
Compose. Native Node smoke test reproduced host-publication failure RED before
adding frontend's web network. Backend stays internal. GREEN: 27 cases passed.
T2: all existing gates pass: 244 frontend tests, 632 Go leaf cases. Coverage:
99.66% frontend statements / 97.95% branches / 100% functions and lines; Go 98.9%.
Fresh ignored HTML reports generated; frontend CI artifact configured. Parent
contract spot check and independent functional recheck passed. Chrome verified
seven operations, complex expression 20, Enter focus, error recovery, copy/reuse
and history/theme persistence. Controlled outage yielded 504, unhealthy readiness
and safe UI error; restart recovered health and result 7. Final 27 smoke cases pass.
T3: final README covers Docker/local startup, API, architecture, testing, coverage
access and limitations. Reported requirements have evidence; original brief
unavailable remains explicit. PROMPTS contains actual requests and results;
prior bytes are unchanged. No application/backend/dependency changes or excluded
performance tools. No commit or push.
Risk: items 2 (container isolation) and 5 (packaging/CI), verified statically and
against actual runtime behavior. Compose remains available on port 8088 for
manual review. No other local servers or browser data were removed.


### L4 — Submission preparation (actual user prompt)
Prepare the repository for final assessment submission.

- Curate PROMPTS.md for reviewer readability, omitting non-technical communication preferences, personal information, and repetitive operational instructions.
- Clearly disclose the editorial curation and preserve the original prompt archive.
- Preserve technical decisions, corrections, deviations, and genuine TDD evidence.
- Ensure coverage reports are reproducible and accessible to reviewers.
- Verify that the README is concise, accurate, and focused on running and evaluating the application.
- Verify that the "Evaluated by Go · IEEE 754 arithmetic" text is no longer displayed in the UI.
- Remove temporary files and check for accidental secrets or generated artifacts.

Do not expand the scope, introduce dependencies, or rewrite Git history.

Report final submission readiness.

Do not commit or push yet.


Original 85,100-byte prompt archive preserved unchanged under ignored
.submission-private/, SHA-256 recorded in the curated PROMPTS. Every quoted
excerpt was mechanically checked against the original/request. Full archive was
mapped read-only before curation; technical deviations and missing historical
Phase4 execution-evidence disclosure remain. PROMPTS reduced from 2160 to243
lines before adding current check results; README reduced from250 to191 lines
before the final privacy note. No reconstructed original prompts or new RED
claims. The archive is private/local, not shipped in Git or a submission ZIP.

The user had already replaced the arithmetic footer with an author credit.
Preserved that content, normalized its formatting and added a passing absence
characterization test. All frontend gates passed with245tests; Docker rebuild
and27HTTPsmoke cases passed. ReloadedChrome visibly showed no obsolete footer.
Coverage reports remain reproducible and ignored; removed only the generated
backend/coverage/tests.json counting log. Scanned134tracked/nonignoredfiles:
no credential-pattern or generated/private artifact findings. Git diff check
passes; nothing staged, no commit/push or history rewriting.
