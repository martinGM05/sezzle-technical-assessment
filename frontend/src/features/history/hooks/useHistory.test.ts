import { StrictMode } from 'react'
import { act, cleanup, renderHook } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { useHistory } from './useHistory'

const key = 'sezzle.history.v1'
const entry = {
  id: 'entry-1',
  expression: '1 − 1',
  result: '0',
  timestamp: '2026-11-26T12:00:00.000Z',
}

// Matrix: empty/reload; invalid schemas and values; ordered/bounded unique additions;
// delete/clear; access/read/quota failures; StrictMode and unrelated-key isolation.
describe('useHistory', () => {
  beforeEach(() => localStorage.clear())
  afterEach(() => {
    cleanup()
    vi.restoreAllMocks()
  })

  it('starts empty without writing or clearing unrelated storage', () => {
    localStorage.setItem('unrelated', 'keep')
    const write = vi.spyOn(Storage.prototype, 'setItem')
    const { result } = renderHook(useHistory, { wrapper: StrictMode })
    expect(result.current.entries).toEqual([])
    expect(result.current.storageError).toBeNull()
    expect(write).not.toHaveBeenCalled()
    expect(localStorage.getItem('unrelated')).toBe('keep')
  })

  it('reloads validated persisted results including zero, newest first', () => {
    const newer = {
      ...entry,
      id: 'entry-2',
      timestamp: '2026-11-26T12:01:00.000Z',
      result: '1e+21',
    }
    localStorage.setItem(
      key,
      JSON.stringify({ version: 1, entries: [entry, newer] }),
    )
    const { result } = renderHook(useHistory)
    expect(result.current.entries).toEqual([newer, entry])
    expect(result.current.storageError).toBeNull()
  })

  it.each([
    '{',
    'null',
    '[]',
    '{}',
    JSON.stringify({ version: 2, entries: [] }),
    JSON.stringify({ version: 1, entries: 'wrong' }),
    ...[
      null,
      { ...entry, id: '' },
      { ...entry, id: 1 },
      { ...entry, expression: '' },
      { ...entry, expression: 'x'.repeat(4097) },
      { ...entry, result: 0 },
      { ...entry, result: '' },
      { ...entry, result: 'NaN' },
      { ...entry, result: 'Infinity' },
      { ...entry, result: '1e999' },
      { ...entry, result: '0x10' },
      { ...entry, result: ' '.repeat(2) },
      { ...entry, result: '1'.repeat(129) },
      { ...entry, timestamp: 'yesterday' },
      { ...entry, timestamp: '2026-02-30T00:00:00.000Z' },
    ].map((bad) => JSON.stringify({ version: 1, entries: [bad] })),
    JSON.stringify({ version: 1, entries: [entry, entry] }),
    JSON.stringify({
      version: 1,
      entries: Array.from({ length: 101 }, (_, i) => ({
        ...entry,
        id: String(i),
      })),
    }),
  ])('rejects invalid persisted data without mutating it (%#)', (raw) => {
    localStorage.setItem(key, raw)
    const { result } = renderHook(useHistory)
    expect(result.current.entries).toEqual([])
    expect(result.current.storageError).toEqual(expect.any(String))
    expect(localStorage.getItem(key)).toBe(raw)
  })

  it('adds once per call under StrictMode with unique IDs and ISO dates, then reloads', () => {
    const { result, unmount } = renderHook(useHistory, { wrapper: StrictMode })
    act(() => {
      result.current.add('1 − 1', '0')
      result.current.add('2 × 3', '6')
    })
    expect(result.current.entries.map(({ result }) => result)).toEqual([
      '6',
      '0',
    ])
    const entries = result.current.entries
    expect(new Set(entries.map(({ id }) => id)).size).toBe(2)
    for (const item of entries) {
      expect(item.id).toMatch(/^[a-f0-9-]{36}$/)
      expect(new Date(item.timestamp).toISOString()).toBe(item.timestamp)
    }
    expect(JSON.parse(localStorage.getItem(key)!)).toEqual({
      version: 1,
      entries,
    })
    unmount()
    expect(renderHook(useHistory).result.current.entries).toEqual(entries)
  })

  it('bounds new history to the 100 most recent entries', () => {
    const { result } = renderHook(useHistory)
    act(() => {
      for (let i = 0; i < 102; i++) result.current.add(String(i), String(i))
    })
    expect(result.current.entries).toHaveLength(100)
    expect(result.current.entries[0].result).toBe('101')
    expect(result.current.entries[99].result).toBe('2')
    expect(JSON.parse(localStorage.getItem(key)!).entries).toHaveLength(100)
  })

  it('persists delete and clear without changing other keys', () => {
    localStorage.setItem('unrelated', 'keep')
    const { result } = renderHook(useHistory)
    act(() => {
      result.current.add('1', '1')
      result.current.add('2', '2')
    })
    act(() => result.current.remove(result.current.entries[0].id))
    expect(result.current.entries.map(({ result }) => result)).toEqual(['1'])
    expect(JSON.parse(localStorage.getItem(key)!).entries).toEqual(
      result.current.entries,
    )
    act(() => result.current.clear())
    expect(result.current.entries).toEqual([])
    expect(JSON.parse(localStorage.getItem(key)!)).toEqual({
      version: 1,
      entries: [],
    })
    expect(localStorage.getItem('unrelated')).toBe('keep')
  })

  it.each(['getter', 'read', 'write'])(
    'retains in-memory interaction when storage %s fails',
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
      const { result } = renderHook(useHistory)
      if (failure !== 'write')
        expect(result.current.storageError).toEqual(expect.any(String))
      act(() => result.current.add('2 + 2', '4'))
      expect(result.current.entries).toHaveLength(1)
      if (failure !== 'read')
        expect(result.current.storageError).toEqual(expect.any(String))
      act(() => result.current.clear())
      expect(result.current.entries).toEqual([])
    },
  )

  it('ignores invalid additions while preserving valid in-memory history', () => {
    const { result } = renderHook(useHistory)
    act(() => result.current.add('1', '1'))
    act(() => {
      result.current.add('', '1')
      result.current.add('1', 'NaN')
    })
    expect(result.current.entries).toHaveLength(1)
    expect(result.current.storageError).toEqual(expect.any(String))
  })
})
