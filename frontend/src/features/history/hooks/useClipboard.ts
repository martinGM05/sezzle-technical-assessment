import { useEffect, useRef, useState } from 'react'

export function useClipboard() {
  const [message, setMessage] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const generation = useRef(0)

  useEffect(
    () => () => {
      generation.current++
    },
    [],
  )

  async function copy(value: string) {
    const current = ++generation.current
    setMessage(null)
    setError(null)
    try {
      if (!navigator.clipboard?.writeText)
        throw new Error('Clipboard unavailable')
      await navigator.clipboard.writeText(value)
      if (generation.current === current) setMessage('Result copied.')
    } catch {
      if (generation.current === current)
        setError('Could not copy the result. Select and copy it manually.')
    }
  }

  return { copy, message, error }
}
