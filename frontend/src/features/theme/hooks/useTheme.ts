import { useCallback, useEffect, useState } from 'react'
import type { Theme } from '../../../components/ui/ThemeToggle/ThemeToggle'
import { readStorage, writeStorage } from '../../../lib/storage'

const KEY = 'sezzle.theme.v1'
type ThemeState = { theme: Theme; storageError: string | null }

function loadTheme(): ThemeState {
  const { value, error } = readStorage(KEY)
  if (error || value === null) return { theme: 'system', storageError: error }
  if (typeof value === 'object') {
    const data = value as Record<string, unknown>
    if (
      data.version === 1 &&
      (data.theme === 'light' ||
        data.theme === 'dark' ||
        data.theme === 'system')
    ) {
      return { theme: data.theme, storageError: null }
    }
  }
  return {
    theme: 'system',
    storageError: 'Saved theme is invalid. Using the system theme.',
  }
}

export function useTheme() {
  const [state, setState] = useState(loadTheme)

  useEffect(() => {
    document.documentElement.dataset.theme = state.theme
  }, [state.theme])

  const setTheme = useCallback((theme: Theme) => {
    const storageError = writeStorage(KEY, { version: 1, theme })
    setState({ theme, storageError })
  }, [])

  return { ...state, setTheme }
}
