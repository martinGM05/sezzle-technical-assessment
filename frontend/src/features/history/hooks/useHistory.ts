import { useCallback, useRef, useState } from 'react'
import { readStorage, writeStorage } from '../../../lib/storage'
import type { HistoryEntry } from '../types'

const KEY = 'sezzle.history.v1'
const LIMIT = 100
const INVALID = 'Saved history is invalid and could not be loaded.'
const INVALID_ADDITION = 'This calculation could not be added to history.'
type HistoryState = {
  entries: readonly HistoryEntry[]
  storageError: string | null
}

function isEntry(value: unknown): value is HistoryEntry {
  if (typeof value !== 'object' || value === null) return false
  const item = value as Record<string, unknown>
  if (
    typeof item.id !== 'string' ||
    !item.id.trim() ||
    item.id.length > 128 ||
    typeof item.expression !== 'string' ||
    !item.expression.trim() ||
    item.expression.length > 4096 ||
    typeof item.result !== 'string' ||
    item.result.length > 128 ||
    !/^-?(?:\d+(?:\.\d+)?|\.\d+)(?:e[+-]?\d+)?$/i.test(item.result) ||
    !Number.isFinite(Number(item.result)) ||
    typeof item.timestamp !== 'string' ||
    item.timestamp.length !== 24
  )
    return false
  const time = new Date(item.timestamp)
  return (
    Number.isFinite(time.getTime()) && time.toISOString() === item.timestamp
  )
}

function loadHistory(): HistoryState {
  const { value, error } = readStorage(KEY)
  if (error || value === null) return { entries: [], storageError: error }
  if (typeof value !== 'object') return { entries: [], storageError: INVALID }
  const data = value as Record<string, unknown>
  if (
    data.version !== 1 ||
    !Array.isArray(data.entries) ||
    data.entries.length > LIMIT ||
    !data.entries.every(isEntry) ||
    new Set(data.entries.map((entry) => entry.id)).size !== data.entries.length
  )
    return { entries: [], storageError: INVALID }
  return {
    entries: data.entries.toSorted((a, b) =>
      b.timestamp.localeCompare(a.timestamp),
    ),
    storageError: null,
  }
}

export function useHistory() {
  const [state, setState] = useState(loadHistory)
  const entriesRef = useRef(state.entries)

  const update = useCallback((entries: readonly HistoryEntry[]) => {
    entriesRef.current = entries
    const storageError = writeStorage(KEY, { version: 1, entries })
    setState({ entries, storageError })
  }, [])

  const add = useCallback(
    (expression: string, result: string) => {
      const entry = {
        id: crypto.randomUUID(),
        expression,
        result,
        timestamp: new Date().toISOString(),
      }
      if (!isEntry(entry)) {
        setState((previous) => ({
          ...previous,
          storageError: INVALID_ADDITION,
        }))
        return
      }
      update([entry, ...entriesRef.current].slice(0, LIMIT))
    },
    [update],
  )

  const remove = useCallback(
    (id: string) => {
      update(entriesRef.current.filter((entry) => entry.id !== id))
    },
    [update],
  )

  const clear = useCallback(() => update([]), [update])
  return { ...state, add, remove, clear }
}
