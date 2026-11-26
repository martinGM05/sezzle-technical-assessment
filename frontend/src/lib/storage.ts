const READ_ERROR =
  'Saved preferences could not be read. Changes will remain available for this session.'
const WRITE_ERROR =
  'Changes could not be saved. They remain available for this session.'

export function readStorage(key: string): {
  value: unknown
  error: string | null
} {
  try {
    const raw = window.localStorage.getItem(key)
    if (raw === null) return { value: null, error: null }
    const value: unknown = JSON.parse(raw)
    return { value, error: value === null ? READ_ERROR : null }
  } catch {
    return { value: null, error: READ_ERROR }
  }
}

export function writeStorage(key: string, value: unknown): string | null {
  try {
    window.localStorage.setItem(key, JSON.stringify(value))
    return null
  } catch {
    return WRITE_ERROR
  }
}
