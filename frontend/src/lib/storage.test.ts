import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { readStorage, writeStorage } from './storage'

// Matrix: missing/valid/malformed JSON, access/read/write errors, unrelated keys.
describe('safe JSON storage', () => {
  beforeEach(() => localStorage.clear())
  afterEach(() => vi.restoreAllMocks())

  it('reads missing and valid JSON without changing other keys', () => {
    localStorage.setItem('unrelated', 'keep')
    expect(readStorage('missing')).toEqual({ value: null, error: null })
    localStorage.setItem('example', '{"version":1}')
    expect(readStorage('example')).toEqual({
      value: { version: 1 },
      error: null,
    })
    expect(localStorage.getItem('unrelated')).toBe('keep')
  })

  it('reports malformed JSON without removing it', () => {
    localStorage.setItem('example', '{')
    expect(readStorage('example').error).toEqual(expect.any(String))
    expect(localStorage.getItem('example')).toBe('{')
  })

  it('contains storage access and read failures', () => {
    vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => {
      throw new Error('private')
    })
    expect(readStorage('example')).toEqual({
      value: null,
      error: expect.any(String),
    })
    vi.restoreAllMocks()
    vi.spyOn(window, 'localStorage', 'get').mockImplementation(() => {
      throw new Error('private')
    })
    expect(readStorage('example')).toEqual({
      value: null,
      error: expect.any(String),
    })
    expect(writeStorage('example', {})).toEqual(expect.any(String))
  })

  it('writes JSON and reports quota errors without exposing their details', () => {
    expect(writeStorage('example', { version: 1 })).toBeNull()
    expect(JSON.parse(localStorage.getItem('example')!)).toEqual({ version: 1 })
    vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {
      throw new DOMException('private', 'QuotaExceededError')
    })
    const error = writeStorage('example', {})
    expect(error).toEqual(expect.any(String))
    expect(error).not.toContain('private')
  })
})
