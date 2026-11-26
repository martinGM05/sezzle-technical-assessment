import { act, cleanup, renderHook } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { useClipboard } from './useClipboard'

function deferred() {
  let resolve!: () => void
  const promise = new Promise<void>((done) => {
    resolve = done
  })
  return { promise, resolve }
}

const original = Object.getOwnPropertyDescriptor(navigator, 'clipboard')
afterEach(() => {
  cleanup()
  if (original) Object.defineProperty(navigator, 'clipboard', original)
  else Reflect.deleteProperty(navigator, 'clipboard')
})

// Matrix: confirmed copy, unavailable browser API, rejected permission, stale
// completion and unmount. Clipboard writes happen only from an explicit action.
describe('useClipboard', () => {
  it('copies the exact result only on request and confirms success', async () => {
    const writeText = vi.fn().mockResolvedValue(undefined)
    Object.defineProperty(navigator, 'clipboard', {
      configurable: true,
      value: { writeText },
    })
    const { result } = renderHook(useClipboard)
    expect(writeText).not.toHaveBeenCalled()
    await act(async () => result.current.copy('0.30000000000000004'))
    expect(writeText).toHaveBeenCalledExactlyOnceWith('0.30000000000000004')
    expect(result.current.message).toBe('Result copied.')
    expect(result.current.error).toBeNull()
  })
  it.each(['unavailable', 'denied'] as const)(
    'reports %s clipboard without throwing',
    async (mode) => {
      Object.defineProperty(navigator, 'clipboard', {
        configurable: true,
        value:
          mode === 'unavailable'
            ? undefined
            : {
                writeText: vi
                  .fn()
                  .mockRejectedValue(new Error('private permission info')),
              },
      })
      const { result } = renderHook(useClipboard)
      await act(async () => result.current.copy('5'))
      expect(result.current.error).toBe(
        'Could not copy the result. Select and copy it manually.',
      )
      expect(result.current.message).toBeNull()
    },
  )
  it('ignores an older completion after another copy fails', async () => {
    const first = deferred()
    const writeText = vi
      .fn()
      .mockReturnValueOnce(first.promise)
      .mockRejectedValueOnce(new Error('denied'))
    Object.defineProperty(navigator, 'clipboard', {
      configurable: true,
      value: { writeText },
    })
    const { result } = renderHook(useClipboard)
    act(() => {
      void result.current.copy('1')
    })
    await act(async () => result.current.copy('2'))
    await act(async () => first.resolve())
    expect(result.current.message).toBeNull()
    expect(result.current.error).not.toBeNull()
  })
  it('safely settles after unmount', async () => {
    const write = deferred()
    Object.defineProperty(navigator, 'clipboard', {
      configurable: true,
      value: { writeText: vi.fn().mockReturnValue(write.promise) },
    })
    const { result, unmount } = renderHook(useClipboard)
    let pending: Promise<void> | undefined
    act(() => {
      pending = result.current.copy('1')
    })
    unmount()
    write.resolve()
    await expect(pending).resolves.toBeUndefined()
  })
})
