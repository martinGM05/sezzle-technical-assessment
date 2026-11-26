import { afterEach, describe, expect, it, vi } from 'vitest'
import { CalculationApiError, evaluateExpression } from './client'

function stubFetch() {
  const fetch = vi.fn<typeof globalThis.fetch>()
  vi.stubGlobal('fetch', fetch)
  return fetch
}

function jsonResponse(
  body: string,
  status = 200,
  contentType = 'application/json',
) {
  return new Response(body, {
    status,
    headers: { 'Content-Type': contentType },
  })
}

const domainProblem = {
  type: 'about:blank',
  title: 'Unprocessable Content',
  status: 422,
  detail: 'Cannot divide by zero.',
  code: 'DIVISION_BY_ZERO',
}

afterEach(() => vi.unstubAllGlobals())

describe('evaluateExpression', () => {
  // Boundary matrix: request and finite success, validated domain errors, untrusted
  // HTTP responses, transport failures, and cancellation before/during fetch or body reads.
  it('posts normalized JSON to the same-origin expression endpoint and returns zero', async () => {
    const fetch = stubFetch().mockResolvedValue(
      jsonResponse('{"result":0}', 200, 'application/json; charset=utf-8'),
    )
    const { signal } = new AbortController()

    await expect(evaluateExpression('√(9)×(−2)÷2+3', signal)).resolves.toBe(0)
    expect(fetch).toHaveBeenCalledExactlyOnceWith('/api/v1/expressions', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ expression: 'sqrt(9)*(-2)/2+3' }),
      signal,
    })
  })

  it('exposes a documented RFC 9457 domain problem', async () => {
    stubFetch().mockResolvedValue(
      jsonResponse(
        JSON.stringify(domainProblem),
        422,
        'application/problem+json',
      ),
    )

    const result = evaluateExpression('1/0')
    await expect(result).rejects.toBeInstanceOf(CalculationApiError)
    await expect(result).rejects.toMatchObject({
      code: domainProblem.code,
      message: domainProblem.detail,
      status: 422,
    })
  })

  it.each([
    ['malformed JSON', '{"result":'],
    ['out-of-range number', '{"result":1e999}'],
    ['numeric string', '{"result":"3"}'],
    ['missing result', '{}'],
  ])('rejects a 200 response with %s', async (_scenario, body) => {
    stubFetch().mockResolvedValue(jsonResponse(body))

    await expect(evaluateExpression('1+2')).rejects.toMatchObject({
      code: 'INVALID_RESPONSE',
      status: 200,
    })
  })

  it.each([
    ['a non-JSON 200 response', 200, 'text/html'],
    ['an undocumented 2xx response', 201, 'application/json'],
  ])(
    'rejects %s even if the body contains a result',
    async (_scenario, status, contentType) => {
      stubFetch().mockResolvedValue(
        jsonResponse('{"result":3}', status, contentType),
      )

      await expect(evaluateExpression('1+2')).rejects.toMatchObject({
        code: 'INVALID_RESPONSE',
        status,
      })
    },
  )

  it.each([
    ['a mismatched status', { status: 500 }],
    ['a blank detail', { detail: ' ' }],
    ['an unknown error code', { code: 'UNKNOWN_CODE' }],
    ['a missing title', { title: undefined }],
  ])('does not expose a domain problem with %s', async (_scenario, fields) => {
    stubFetch().mockResolvedValue(
      jsonResponse(
        JSON.stringify({
          ...domainProblem,
          detail: 'Untrusted details.',
          ...fields,
        }),
        422,
        'application/problem+json',
      ),
    )

    await expect(evaluateExpression('1/0')).rejects.toMatchObject({
      code: 'INVALID_RESPONSE',
      status: 422,
      message: 'The calculator received an invalid response.',
    })
  })

  it('does not expose proxy or server errors', async () => {
    const fetch = stubFetch()
      .mockResolvedValueOnce(
        jsonResponse('Untrusted proxy error.', 502, 'text/html'),
      )
      .mockResolvedValueOnce(
        jsonResponse(
          JSON.stringify({
            ...domainProblem,
            status: 500,
            detail: 'Untrusted server details.',
          }),
          500,
          'application/problem+json',
        ),
      )

    for (const status of [502, 500]) {
      await expect(evaluateExpression('1+2')).rejects.toMatchObject({
        code: 'SERVER_ERROR',
        message: 'The calculator service could not complete the request.',
        status,
      })
    }
    expect(fetch).toHaveBeenCalledTimes(2)
  })

  it('reports fetch rejection using a local network error message', async () => {
    stubFetch().mockRejectedValue(new TypeError('Untrusted network error.'))

    await expect(evaluateExpression('1+2')).rejects.toMatchObject({
      code: 'NETWORK_ERROR',
      status: null,
      message: 'Could not reach the calculator service. Please try again.',
    })
  })

  it('preserves an AbortError rejected by fetch', async () => {
    const abortError = new DOMException('Canceled.', 'AbortError')
    stubFetch().mockRejectedValue(abortError)

    await expect(evaluateExpression('1+2')).rejects.toBe(abortError)
  })

  it('preserves an AbortError from reading the response body', async () => {
    const abortError = new DOMException('Canceled.', 'AbortError')
    const body = new ReadableStream({
      start: (stream) => stream.error(abortError),
    })
    stubFetch().mockResolvedValue(
      new Response(body, { headers: { 'Content-Type': 'application/json' } }),
    )

    await expect(evaluateExpression('1+2')).rejects.toBe(abortError)
  })

  it('skips fetch for a pre-aborted signal, even with a custom reason', async () => {
    const fetch = stubFetch()
    const signal = AbortSignal.abort('Canceled.')

    await expect(evaluateExpression('1+2', signal)).rejects.toMatchObject({
      name: 'AbortError',
    })
    expect(fetch).not.toHaveBeenCalled()
  })
})
