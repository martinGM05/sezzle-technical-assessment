# Backend review findings

The four-lens native backend review was approved and acknowledged on 2026-11-26
without blockers or a correction round. It produced **five WARNING findings**
representing three overlapping concerns. The entries below summarize all five;
locations refer to the reviewed revision, not necessarily current line numbers.
Approval covered the prior backend candidate, not subsequent delivery edits.

Review: `review-25465c4a7f1e810e`.
Reviewed target: `a3ff5dbbd95a81b2710f06b197fc140f69c82cda0553ea05769f58983e1fa816`.
The native reviewers inspected immutable patches; separate functional verification
ran the tests and real-process lifecycle probes. No native test run is implied.

## R1-001 — Duplicate JSON properties

- Lens: risk. Location: `backend/internal/transport/http/request.go:38-39`.
- Finding: decoding into a map overwrites earlier duplicate keys before validation.
  Conflicting operation values may produce the calculation selected by the last
  value rather than the value an intermediary inferred from the first.
- Status: **open, accepted non-blocking limitation**. This delivery preserves the
  existing decoder behavior. Rejecting duplicates would be a later hardening
  change with explicit contract documentation and failing tests first.

## R2-001 — HEAD documentation contradicts its bodyless description

- Lens: readability. Location: `docs/openapi.yaml:6`.
- Finding: the global description says HEAD is bodyless, but health HEAD 403/500
  declared JSON content. A generated client could attempt to decode an empty body.
  The reviewer recommended removing content, documenting the media-type header
  and checking that HEAD schemas contain no response content.
- Status: **resolved**. HEAD 200/403/500 have no content declarations; response
  descriptions state the actual Content-Type. Contract tests verify both facts.
- The original text was recovered from the retained native completion event;
  it was not reconstructed from the earlier truncated tool output.
- Standards clarification: OpenAPI 3.1 ignores a `Content-Type` entry in a
  response's `headers` map. Therefore the media type is described in prose,
  rather than adding an ignored header schema or a false body declaration.
  See the [OpenAPI Response Object](https://spec.openapis.org/oas/v3.1.0.html#response-object).

## R3-001 — HEAD error schema parity gap

- Lens: reliability. Location: `docs/openapi.yaml:55-79` in the frozen review.
- Finding: the parity helper returned after checking an empty HEAD body without
  rejecting response content declared in OpenAPI.
- Status: **resolved** by the same contract correction as R2-001. The parity
  helper now rejects content on the documented health HEAD operation, and a
  dedicated table checks all three response statuses. Wire behavior is unchanged.

## R3-002 — Duplicate members can hide invalid earlier values

- Lens: reliability. Location: `backend/internal/transport/http/request.go:37-38`.
- Finding: a later correctly typed operation/operands member can overwrite an
  earlier wrong-typed member before validation and produce a successful result.
- Status: **open, accepted non-blocking limitation**, overlapping R1-001. For
  example, duplicate `operation` members retain the last value. This is not
  claimed to be strict duplicate-key rejection or silently treated as resolved.

## R4-001 — Same-origin proxy requests can be denied

- Lens: resilience. Location: `backend/internal/transport/http/cors.go:36`.
- Finding: same-origin browser POSTs can carry Origin; a proxy can forward it,
  while the executable's empty allowlist rejected every nonempty Origin.
- Status: **resolved for explicitly configured deployment**. A real HTTPS proxy
  to the HTTP backend reproduced 403. Failing tests then drove minimal executable
  wiring of `ALLOWED_ORIGINS=https://calculator.example` into the existing handler.
  Origin checks remain exact and restrictive; empty configuration still denies
  nonempty Origin. Forwarded headers are not authorization.
- Tests cover successful public-origin requests, multiple configured origins,
  attacker denial, wildcard denial and origin-less clients. No production proxy
  or frontend implementation was introduced.
