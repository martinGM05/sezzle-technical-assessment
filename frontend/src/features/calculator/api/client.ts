import { normalizeExpression } from '../normalizeExpression'

const expressionErrorCodes = new Set([
  'INVALID_EXPRESSION',
  'SYNTAX_ERROR',
  'EXPRESSION_LIMIT',
  'INVALID_OPERANDS',
  'DIVISION_BY_ZERO',
  'INVALID_RESULT',
])

export class CalculationApiError extends Error {
  readonly code: string
  readonly status: number | null

  constructor(message: string, code: string, status: number | null) {
    super(message)
    this.name = 'CalculationApiError'
    this.code = code
    this.status = status
  }
}

export async function evaluateExpression(
  expression: string,
  signal?: AbortSignal,
): Promise<number> {
  throwIfAborted(signal)

  let response: Response
  try {
    response = await fetch('/api/v1/expressions', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ expression: normalizeExpression(expression) }),
      signal,
    })
  } catch (error) {
    throwIfAborted(signal, error)
    throw new CalculationApiError(
      'Could not reach the calculator service. Please try again.',
      'NETWORK_ERROR',
      null,
    )
  }

  throwIfAborted(signal)
  const { status } = response
  if (status !== 200 && status !== 422) {
    if (response.ok) throw invalidResponse(status)
    throw new CalculationApiError(
      'The calculator service could not complete the request.',
      'SERVER_ERROR',
      status,
    )
  }

  const contentType = response.headers
    .get('Content-Type')
    ?.split(';', 1)[0]
    .trim()
    .toLowerCase()
  if (
    contentType !==
    (status === 200 ? 'application/json' : 'application/problem+json')
  ) {
    throw invalidResponse(status)
  }

  const payload = await readJson(response, signal)
  throwIfAborted(signal)
  if (status === 422) {
    if (!isExpressionProblem(payload)) throw invalidResponse(status)
    throw new CalculationApiError(payload.detail, payload.code, status)
  }

  if (
    !isRecord(payload) ||
    typeof payload.result !== 'number' ||
    !Number.isFinite(payload.result)
  ) {
    throw invalidResponse(status)
  }
  return payload.result
}

async function readJson(
  response: Response,
  signal?: AbortSignal,
): Promise<unknown> {
  try {
    return await response.json()
  } catch (error) {
    throwIfAborted(signal, error)
    throw invalidResponse(response.status)
  }
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}

function isExpressionProblem(
  payload: unknown,
): payload is { detail: string; code: string } {
  return (
    isRecord(payload) &&
    payload.type === 'about:blank' &&
    typeof payload.title === 'string' &&
    payload.title.trim() !== '' &&
    payload.status === 422 &&
    typeof payload.detail === 'string' &&
    payload.detail.trim() !== '' &&
    typeof payload.code === 'string' &&
    expressionErrorCodes.has(payload.code)
  )
}

function invalidResponse(status: number): CalculationApiError {
  return new CalculationApiError(
    'The calculator received an invalid response.',
    'INVALID_RESPONSE',
    status,
  )
}

function throwIfAborted(signal?: AbortSignal, error?: unknown): void {
  if (isRecord(error) && error.name === 'AbortError') throw error
  if (!signal?.aborted) return

  const reason: unknown = signal.reason
  if (isRecord(reason) && reason.name === 'AbortError') throw reason
  throw new DOMException('The calculation was canceled.', 'AbortError')
}
