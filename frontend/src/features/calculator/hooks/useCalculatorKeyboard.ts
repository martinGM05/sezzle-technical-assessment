import { useEffect } from 'react'
import type { CalculatorAction, InputCharacter, Operation } from '../types'

const inputCharacters: readonly InputCharacter[] = [
  '0',
  '1',
  '2',
  '3',
  '4',
  '5',
  '6',
  '7',
  '8',
  '9',
  '.',
  '(',
  ')',
]
const keyboardOperations: Record<string, Operation> = {
  '+': 'add',
  '-': 'subtract',
  '*': 'multiply',
  '/': 'divide',
  '^': 'power',
  '%': 'percentage',
  r: 'sqrt',
}

export function useCalculatorKeyboard(
  onAction: (action: CalculatorAction) => void,
) {
  useEffect(() => {
    function handleKey(event: KeyboardEvent) {
      if (
        event.defaultPrevented ||
        event.isComposing ||
        event.ctrlKey ||
        event.metaKey ||
        event.altKey
      )
        return
      const target = event.target
      if (target instanceof Element) {
        if (target.closest('input, textarea, select, [contenteditable]')) return
        if (
          target !== document.body &&
          target !== document.documentElement &&
          !target.closest('[data-testid="calculator"]')
        )
          return
        // Space retains native button activation. Enter is the calculator's
        // evaluate shortcut regardless of which calculator key retains focus.
      }
      let action: CalculatorAction | undefined
      const input = inputCharacters.find((value) => value === event.key)
      const operation = Object.hasOwn(keyboardOperations, event.key)
        ? keyboardOperations[event.key]
        : undefined
      if (input) action = { type: 'input', value: input }
      else if (operation) action = { type: 'operation', operation }
      else if (event.key === 'Enter' || event.key === '=') {
        // Cancel the native button click even when ignoring held-key repeats.
        event.preventDefault()
        if (event.repeat) return
        action = { type: 'equals' }
      } else if (event.key === 'Backspace') action = { type: 'backspace' }
      else if (event.key === 'Escape' || event.key === 'Delete')
        action = { type: 'clear' }
      if (!action) return
      event.preventDefault()
      onAction(action)
    }
    window.addEventListener('keydown', handleKey)
    return () => window.removeEventListener('keydown', handleKey)
  }, [onAction])
}
