# Phase 2 backend delivery

Branch: `feat/calculator-api`. Route: parent integration with parallel Docker and CI units.
Delivery: local changes only; no commits, push, frontend or new business features. Runner: Go testing and local Docker.

## Specs

- S1: "Add the approved multi-stage backend Dockerfile using a verified digest-pinned distroless non-root runtime." "Add `.dockerignore`."
- S2: "Verify container build, non-root execution, health, all seven operations, and graceful shutdown."
- S3: "Complete backend CI with actual coverage reporting and artifacts."
- S4: "Finalize backend setup instructions, API examples, and architecture documentation." "Document all five non-blocking review findings."
- S5: "Correct the OpenAPI HEAD response documentation."
- S6: "Verify whether the reported CORS behavior affects the planned same-origin reverse proxy deployment." "Use TDD for any necessary behavioral correction."
- S7: "Run the final quality gates and report actual results." "Record this prompt in PROMPTS.md."
- S8: "Do not introduce new features, modify the frontend, commit, or push."

## Tasks

- T1 (S1-S2, S8): delegated Docker unit; complete (digest/build/non-root/API/shutdown verified); no commit authorized.
- T2 (S3, S8): delegated CI unit; complete (local steps and reports verified; remote run not authorized); no commit authorized.
- T3 (S5-S6, S8): inline contract/proxy characterization and justified corrections; complete (observed RED then GREEN).
- T4 (S4, S7-S8): inline documentation, integration and final verification; complete; no commit authorized.

## Log

### L1 — original user request

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

### L2 — inspection and boundaries

- Existing untracked Dockerfile/runtime tag and backend artifact action are not digest/SHA pinned. Existing frontend CI must remain byte-for-byte unchanged.
- Docker Desktop was stopped; launched locally with explicit tool approval.
- Native review was approved and acknowledged. Four warnings have retained claims; R2-001 (`docs/openapi.yaml:6`) was truncated. Recover its original evidence if available; never reconstruct it as an original finding.
- Parallel writer reason: Docker packaging/real container verification and CI coverage/artifact configuration are independent multi-minute units with disjoint edit surfaces. Parent owns Go, OpenAPI, docs, prompts, task state and git.
- Forecast: approximately 350-500 authored lines plus verification tests/documentation. No commit slicing is performed because the user expressly forbids commits and push.

### L3 — observed progress

- T1: verified builder/runtime image-index digests; image built from final CORS source. Actual UID 65532, health and seven operations passed. SIGTERM idle/draining exit 0; deadline forced close exit 1 after 5.244 seconds. No signature/vulnerability/cross-platform claim.
- T2: backend-only workflow now generates profile/text/HTML and job summary, pins upload-artifact v4.6.2 by verified SHA. Local equivalent checks pass; no GitHub run/upload authorized. Frontend job unchanged.
- T3: observed two proxy RED failures for configured origins, then GREEN after executable ALLOWED_ORIGINS wiring. Default and hostile origins remain denied. Two HEAD content-schema failures corrected; media type documented in descriptions because OpenAPI ignores Content-Type response-header entries.
- T4: recovered exact R2-001 native claim from original completion event. All five findings documented; duplicate-key concerns remain open, HEAD and configured-proxy concerns resolved. Final integrated checks pending.

### L4 — final verification

- Parent final suite:394 leaf cases/390 subtests,0failures/0skips. Race/vet/build/gofmt/contract/diff checks pass. Coverage98.4% total,cmd90.3%,domain/HTTP100%.
- Parent rebuild succeeded and runtime probe passed health/seven operations/HEAD denial, actual UID65532 and SIGTERM exit0. Initial docker top probe lacked PID column; corrected probe passed. Worker separately verified active drain and bounded deadline exit1.
- Independent verifier: no reproducible blockers, full gates repeated, baseline proxy and HEAD documentation failures reproduced. Remote CI/artifact upload not authorized; no image signature/vulnerability/cross-architecture checks claimed.
- Frontend source and frontend CI prefix unchanged. No commits or pushes. Two duplicate JSON-member warnings remain accepted followups; three overlapping HEAD/CORS warnings resolved.
