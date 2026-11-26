import assert from "node:assert/strict";
import { test } from "node:test";

// Native Node smoke tests against a running Compose stack, never mocked fetch.
const base = process.env.APP_URL ?? "http://127.0.0.1:8088";
const origin = new URL(base).origin;
function request(path, options = {}) {
  return fetch(new URL(path, base), {
    signal: AbortSignal.timeout(5000),
    ...options,
  });
}
function post(path, value, headers = {}) {
  return request(path, {
    method: "POST",
    headers: { "Content-Type": "application/json", Origin: origin, ...headers },
    body: JSON.stringify(value),
  });
}

await test("published health verifies Nginx and the Go proxy", async () => {
  const frontend = await request("/healthz");
  assert.equal(frontend.status, 200);
  assert.equal(await frontend.text(), "ok\n");
  const backend = await request("/health");
  assert.equal(backend.status, 200);
  assert.deepEqual(await backend.json(), { status: "ok" });
});

await test("serves the built React entry point and its JS/CSS assets", async () => {
  const response = await request("/");
  assert.equal(response.status, 200);
  assert.match(response.headers.get("content-type"), /text\/html/);
  const html = await response.text();
  assert.match(html, /id="root"/);
  const assets = [...html.matchAll(/(?:src|href)="(\/assets\/[^"\s]+)"/g)];
  assert.ok(assets.length >= 2, "Expected built JS and CSS assets");
  for (const [, path] of assets) {
    const asset = await request(path);
    assert.equal(asset.status, 200, path);
    assert.ok((await asset.text()).length > 0, path);
  }
});

for (const { operation, operands, result } of [
  { operation: "add", operands: [2, 3], result: 5 },
  { operation: "subtract", operands: [2, 3], result: -1 },
  { operation: "multiply", operands: [2, 3], result: 6 },
  { operation: "divide", operands: [10, 4], result: 2.5 },
  { operation: "power", operands: [2, 3], result: 8 },
  { operation: "sqrt", operands: [9], result: 3 },
  { operation: "percentage", operands: [15], result: 0.15 },
]) {
  await test(`calculation proxy: ${operation}`, async () => {
    const response = await post("/api/v1/calculations", {
      operation,
      operands,
    });
    assert.equal(response.status, 200);
    assert.deepEqual(await response.json(), { result });
  });
}

for (const { expression, result } of [
  { expression: "2+3*4", result: 14 },
  { expression: "(2+3)*4", result: 20 },
  { expression: "2^3^2", result: 512 },
  { expression: "sqrt(9+7)", result: 4 },
  { expression: "200*15%", result: 30 },
  { expression: "-5+3", result: -2 },
  { expression: "-2^2", result: -4 },
  { expression: "(-2)^2", result: 4 },
  { expression: "2^-3", result: 0.125 },
  { expression: "10/4", result: 2.5 },
  { expression: ".5+.25", result: 0.75 },
]) {
  await test(`expression proxy: ${expression}`, async () => {
    const response = await post("/api/v1/expressions", { expression });
    assert.equal(response.status, 200);
    assert.equal(response.headers.get("access-control-allow-origin"), origin);
    assert.deepEqual(await response.json(), { result });
  });
}

for (const { expression, code } of [
  { expression: "10/(5-5)", code: "DIVISION_BY_ZERO" },
  { expression: "sqrt(-1)", code: "INVALID_OPERANDS" },
  { expression: "(2+3", code: "SYNTAX_ERROR" },
  { expression: "2+*3", code: "SYNTAX_ERROR" },
  { expression: "10^309", code: "INVALID_RESULT" },
]) {
  await test(`preserves Problem Details: ${expression}`, async () => {
    const response = await post("/api/v1/expressions", { expression });
    assert.equal(response.status, 422);
    assert.equal(
      response.headers.get("content-type"),
      "application/problem+json",
    );
    const problem = await response.json();
    assert.equal(problem.type, "about:blank");
    assert.equal(problem.title, "Unprocessable Content");
    assert.equal(problem.status, 422);
    assert.equal(problem.code, code);
    assert.ok(problem.detail.length > 0);
  });
}

await test("preserves restrictive origins and allowed preflight", async () => {
  const denied = await post(
    "/api/v1/expressions",
    { expression: "1+1" },
    { Origin: "https://untrusted.example" },
  );
  assert.equal(denied.status, 403);
  assert.equal((await denied.json()).code, "CORS_DENIED");
  assert.equal(denied.headers.get("access-control-allow-origin"), null);
  const allowed = await request("/api/v1/expressions", {
    method: "OPTIONS",
    headers: {
      Origin: origin,
      "Access-Control-Request-Method": "POST",
      "Access-Control-Request-Headers": "Content-Type",
    },
  });
  assert.equal(allowed.status, 204);
  assert.equal(allowed.headers.get("access-control-allow-origin"), origin);
});

await test("keeps HTTP constraints and backend body limits through the proxy", async () => {
  for (const { path, options, status } of [
    { path: "/api/v1/expressions", options: { method: "GET" }, status: 405 },
    { path: "/api/missing", options: {}, status: 404 },
    {
      path: "/api/v1/expressions",
      options: {
        method: "POST",
        headers: { "Content-Type": "text/plain" },
        body: "{}",
      },
      status: 415,
    },
    {
      path: "/api/v1/expressions",
      options: {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: "{",
      },
      status: 400,
    },
    {
      path: "/api/v1/expressions",
      options: {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ expression: "1".repeat(1024 * 1024) }),
      },
      status: 413,
    },
  ]) {
    const response = await request(path, options);
    assert.equal(response.status, status, `${path}: expected ${status}`);
    assert.equal(
      response.headers.get("content-type"),
      "application/problem+json",
    );
    assert.equal((await response.json()).status, status);
  }
  const head = await request("/health", { method: "HEAD" });
  assert.equal(head.status, 200);
  assert.equal(await head.text(), "");
});
