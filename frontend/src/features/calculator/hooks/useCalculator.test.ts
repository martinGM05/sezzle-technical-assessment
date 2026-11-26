import { act, cleanup, renderHook } from '@testing-library/react'
import { StrictMode } from 'react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { CalculationApiError, evaluateExpression } from '../api/client'
import { useCalculator } from './useCalculator'

vi.mock('../api/client', async (importOriginal) => {
  const original = await importOriginal<typeof import('../api/client')>()
  return { ...original, evaluateExpression: vi.fn() }
})
const evaluate = vi.mocked(evaluateExpression)
afterEach(cleanup)
beforeEach(() => vi.clearAllMocks())

function deferred<T>() {
  let resolve!: (value: T) => void
  let reject!: (reason: unknown) => void
  const promise = new Promise<T>((yes, no) => {
    resolve = yes
    reject = no
  })
  return { promise, resolve, reject }
}

// Matrix: text editing/operations, blank submit, success/error recovery,
// synchronous duplicate lock, abort/clear/reuse/edit/unmount, out-of-order results,
// StrictMode cleanup, exact server result presentation and successful history only.
describe('useCalculator', () => {
  it('edits controlled text and selected operations without local evaluation', () => {
    const { result } = renderHook(() => useCalculator(vi.fn()))
    expect(result.current.expression).toBe('')
    expect(result.current.result).toBeNull()
    act(() => {
      result.current.handleAction({ type: 'input', value: '(' })
      result.current.handleAction({ type: 'input', value: '2' })
      result.current.handleAction({ type: 'input', value: '.' })
      result.current.handleAction({ type: 'input', value: '5' })
      result.current.handleAction({ type: 'input', value: ')' })
      result.current.handleAction({ type: 'operation', operation: 'multiply' })
      result.current.handleAction({ type: 'input', value: '3' })
    })
    expect(result.current.expression).toBe('(2.5)×3')
    expect(result.current.selectedOperation).toBe('multiply')
    expect(evaluate).not.toHaveBeenCalled()
    act(() => result.current.handleAction({ type: 'backspace' }))
    expect(result.current.expression).toBe('(2.5)×')
    act(() => result.current.handleAction({ type: 'clear' }))
    act(() => result.current.handleAction({ type: 'backspace' }))
    expect(result.current.expression).toBe('')
    expect(result.current.selectedOperation).toBeNull()
  })
  it.each([
    { operation: 'add' as const, token: '+' },
    { operation: 'subtract' as const, token: '−' },
    { operation: 'multiply' as const, token: '×' },
    { operation: 'divide' as const, token: '÷' },
    { operation: 'power' as const, token: '^' },
    { operation: 'sqrt' as const, token: '√(' },
    { operation: 'percentage' as const, token: '%' },
  ])('inserts $operation syntax', ({ operation, token }) => {
    const { result } = renderHook(() => useCalculator(vi.fn()))
    act(() => result.current.handleAction({ type: 'operation', operation }))
    expect(result.current.expression).toBe(token)
    expect(result.current.selectedOperation).toBe(operation)
  })
  it('rejects empty equals without submitting or creating history', () => {
    const onSuccess = vi.fn()
    const { result } = renderHook(() => useCalculator(onSuccess))
    act(() => result.current.handleAction({ type: 'equals' }))
    expect(result.current.error).toBe('Enter an expression.')
    expect(evaluate).not.toHaveBeenCalled()
    expect(onSuccess).not.toHaveBeenCalled()
  })
  it('locks duplicates synchronously and records exactly one confirmed success', async () => {
    const request = deferred<number>()
    evaluate.mockReturnValue(request.promise)
    const onSuccess = vi.fn()
    const { result } = renderHook(() => useCalculator(onSuccess))
    act(() => {
      result.current.handleAction({ type: 'input', value: '2' })
      result.current.handleAction({ type: 'equals' })
      result.current.handleAction({ type: 'equals' })
    })
    expect(evaluate).toHaveBeenCalledTimes(1)
    expect(evaluate).toHaveBeenCalledWith('2', expect.any(AbortSignal))
    expect(result.current.loading).toBe(true)
    expect(onSuccess).not.toHaveBeenCalled()
    await act(async () => request.resolve(2))
    expect(result.current.loading).toBe(false)
    expect(result.current.result).toBe('2')
    expect(onSuccess).toHaveBeenCalledExactlyOnceWith('2', '2')
  })
  it.each([
    { value: 0, shown: '0' },
    { value: -0, shown: '-0' },
    { value: 0.30000000000000004, shown: '0.30000000000000004' },
    { value: 1e30, shown: '1e+30' },
  ])('does not round the server result $shown', async ({ value, shown }) => {
    evaluate.mockResolvedValue(value)
    const { result } = renderHook(() => useCalculator(vi.fn()))
    act(() => result.current.reuse('0.1+0.2'))
    await act(async () => result.current.handleAction({ type: 'equals' }))
    expect(result.current.result).toBe(shown)
  })
  it('shows safe API errors then recovers without storing failed calculations', async () => {
    evaluate
      .mockRejectedValueOnce(
        new CalculationApiError(
          'Cannot divide by zero.',
          'DIVISION_BY_ZERO',
          422,
        ),
      )
      .mockResolvedValueOnce(5)
    const onSuccess = vi.fn()
    const { result } = renderHook(() => useCalculator(onSuccess))
    act(() => result.current.reuse('10÷0'))
    await act(async () => result.current.handleAction({ type: 'equals' }))
    expect(result.current.error).toBe('Cannot divide by zero.')
    expect(result.current.result).toBeNull()
    expect(onSuccess).not.toHaveBeenCalled()
    act(() => {
      result.current.handleAction({ type: 'backspace' })
      result.current.handleAction({ type: 'input', value: '2' })
    })
    expect(result.current.error).toBeNull()
    await act(async () => result.current.handleAction({ type: 'equals' }))
    expect(result.current.result).toBe('5')
    expect(onSuccess).toHaveBeenCalledExactlyOnceWith('10÷2', '5')
  })
  it.each(['clear', 'edit', 'reuse'] as const)(
    'cancels on %s and ignores an obsolete success',
    async (action) => {
      const first = deferred<number>()
      evaluate.mockReturnValue(first.promise)
      const onSuccess = vi.fn()
      const { result } = renderHook(() => useCalculator(onSuccess))
      act(() => result.current.reuse('1+1'))
      act(() => result.current.handleAction({ type: 'equals' }))
      const signal = evaluate.mock.calls[0]?.[1]
      act(() => {
        if (action === 'clear') result.current.handleAction({ type: 'clear' })
        else if (action === 'reuse') result.current.reuse('9')
        else result.current.handleAction({ type: 'input', value: '2' })
      })
      expect(signal?.aborted).toBe(true)
      expect(result.current.loading).toBe(false)
      await act(async () => first.resolve(2))
      expect(result.current.result).toBeNull()
      expect(onSuccess).not.toHaveBeenCalled()
    },
  )
  it('keeps the newer result when an older request rejects late', async () => {
    const first = deferred<number>()
    const second = deferred<number>()
    evaluate
      .mockReturnValueOnce(first.promise)
      .mockReturnValueOnce(second.promise)
    const onSuccess = vi.fn()
    const { result } = renderHook(() => useCalculator(onSuccess))
    act(() => result.current.reuse('1'))
    act(() => result.current.handleAction({ type: 'equals' }))
    act(() => result.current.reuse('2'))
    act(() => result.current.handleAction({ type: 'equals' }))
    await act(async () => second.resolve(2))
    await act(async () => first.reject(new Error('old network failure')))
    expect(result.current.result).toBe('2')
    expect(result.current.error).toBeNull()
    expect(onSuccess).toHaveBeenCalledExactlyOnceWith('2', '2')
  })
  it('does not unlock a newer pending request when an obsolete success arrives', async () => {
    const first = deferred<number>()
    const second = deferred<number>()
    evaluate
      .mockReturnValueOnce(first.promise)
      .mockReturnValueOnce(second.promise)
    const onSuccess = vi.fn()
    const { result } = renderHook(() => useCalculator(onSuccess))
    act(() => result.current.reuse('1'))
    act(() => result.current.handleAction({ type: 'equals' }))
    act(() => result.current.reuse('2'))
    act(() => result.current.handleAction({ type: 'equals' }))
    await act(async () => first.resolve(1))
    expect(result.current.loading).toBe(true)
    expect(result.current.result).toBeNull()
    expect(onSuccess).not.toHaveBeenCalled()
    act(() => result.current.handleAction({ type: 'equals' }))
    expect(evaluate).toHaveBeenCalledTimes(2)
    await act(async () => second.resolve(2))
    expect(result.current.loading).toBe(false)
    expect(result.current.result).toBe('2')
    expect(onSuccess).toHaveBeenCalledExactlyOnceWith('2', '2')
  })
  it('aborts on unmount and never stores its late result', async () => {
    const request = deferred<number>()
    evaluate.mockReturnValue(request.promise)
    const onSuccess = vi.fn()
    const { result, unmount } = renderHook(() => useCalculator(onSuccess), {
      wrapper: StrictMode,
    })
    act(() => result.current.reuse('1'))
    act(() => result.current.handleAction({ type: 'equals' }))
    const signal = evaluate.mock.calls[0]?.[1]
    unmount()
    expect(signal?.aborted).toBe(true)
    await act(async () => request.resolve(1))
    expect(onSuccess).not.toHaveBeenCalled()
  })
  it('silences AbortError and contains unexpected exceptions', async () => {
    evaluate
      .mockRejectedValueOnce(new DOMException('Canceled', 'AbortError'))
      .mockRejectedValueOnce(new Error('private implementation'))
    const { result } = renderHook(() => useCalculator(vi.fn()))
    act(() => result.current.reuse('1'))
    await act(async () => result.current.handleAction({ type: 'equals' }))
    expect(result.current.error).toBeNull()
    expect(result.current.loading).toBe(false)
    await act(async () => result.current.handleAction({ type: 'equals' }))
    expect(result.current.error).toBe('Calculation failed. Please try again.')
  })
})
