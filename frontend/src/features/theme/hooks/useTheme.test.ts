import { StrictMode } from 'react'
import { act, cleanup, renderHook } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { useTheme } from './useTheme'

const key = 'sezzle.theme.v1'

// Matrix: default/system/light/dark reload; malformed/version/value rejection;
// storage getter/read/write failures; DOM synchronization and key isolation.
describe('useTheme', () => {
  beforeEach(() => {
    localStorage.clear()
    delete document.documentElement.dataset.theme
  })
  afterEach(() => {
    cleanup()
    vi.restoreAllMocks()
  })

  it('defaults to system without writing storage, including StrictMode remounts', () => {
    const write = vi.spyOn(Storage.prototype, 'setItem')
    const { result } = renderHook(useTheme, { wrapper: StrictMode })
    expect(result.current.theme).toBe('system')
    expect(result.current.storageError).toBeNull()
    expect(document.documentElement.dataset.theme).toBe('system')
    expect(write).not.toHaveBeenCalled()
  })

  it.each(['light', 'dark', 'system'] as const)(
    'applies, persists and reloads %s',
    (theme) => {
      localStorage.setItem('unrelated', 'keep')
      const { result, unmount } = renderHook(useTheme)
      act(() => result.current.setTheme(theme))
      expect(result.current.theme).toBe(theme)
      expect(document.documentElement.dataset.theme).toBe(theme)
      expect(JSON.parse(localStorage.getItem(key)!)).toEqual({
        version: 1,
        theme,
      })
      expect(localStorage.getItem('unrelated')).toBe('keep')
      unmount()
      expect(renderHook(useTheme).result.current.theme).toBe(theme)
    },
  )

  it.each([
    '{',
    'null',
    '[]',
    '{}',
    '"dark"',
    '{"version":2,"theme":"dark"}',
    '{"version":1,"theme":"invalid"}',
    '{"version":1,"theme":null}',
  ])('uses system and exposes a warning for invalid storage (%#)', (raw) => {
    localStorage.setItem(key, raw)
    const { result } = renderHook(useTheme)
    expect(result.current.theme).toBe('system')
    expect(result.current.storageError).toEqual(expect.any(String))
    expect(document.documentElement.dataset.theme).toBe('system')
    expect(localStorage.getItem(key)).toBe(raw)
  })

  it.each(['getter', 'read', 'write'])(
    'keeps changes in memory when storage %s fails',
    (failure) => {
      if (failure === 'getter')
        vi.spyOn(window, 'localStorage', 'get').mockImplementation(() => {
          throw new Error('private')
        })
      if (failure === 'read')
        vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => {
          throw new Error('private')
        })
      if (failure === 'write')
        vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {
          throw new DOMException('private', 'QuotaExceededError')
        })
      const { result } = renderHook(useTheme)
      if (failure !== 'write')
        expect(result.current.storageError).toEqual(expect.any(String))
      act(() => result.current.setTheme('dark'))
      expect(result.current.theme).toBe('dark')
      expect(document.documentElement.dataset.theme).toBe('dark')
      if (failure !== 'read')
        expect(result.current.storageError).toEqual(expect.any(String))
    },
  )
})
