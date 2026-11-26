import { useEffect, useRef, useState } from 'react'
import { CalculationApiError, evaluateExpression } from '../api/client'
import { operations } from '../operations'
import type { CalculatorAction, Operation } from '../types'

type PendingRequest = { controller: AbortController }

export function useCalculator(
  onSuccess: (expression: string, result: string) => void,
) {
  const [expression, setExpression] = useState('')
  const [result, setResult] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [selectedOperation, setSelectedOperation] = useState<Operation | null>(
    null,
  )
  // Refs close the same-event gap before React commits an updated render.
  const expressionRef = useRef('')
  const pending = useRef<PendingRequest | null>(null)

  useEffect(
    () => () => {
      pending.current?.controller.abort()
      pending.current = null
    },
    [],
  )

  function edit(value: string) {
    pending.current?.controller.abort()
    pending.current = null
    expressionRef.current = value
    setExpression(value)
    setResult(null)
    setError(null)
    setLoading(false)
  }

  async function submit() {
    if (pending.current) return
    const submitted = expressionRef.current
    if (!submitted.trim()) {
      setError('Enter an expression.')
      return
    }
    const request: PendingRequest = { controller: new AbortController() }
    pending.current = request
    setLoading(true)
    setError(null)
    setResult(null)
    try {
      const value = await evaluateExpression(
        submitted,
        request.controller.signal,
      )
      if (pending.current !== request || request.controller.signal.aborted)
        return
      const displayed = Object.is(value, -0) ? '-0' : String(value)
      setResult(displayed)
      onSuccess(submitted, displayed)
    } catch (failure) {
      if (pending.current !== request || request.controller.signal.aborted)
        return
      if (
        typeof failure === 'object' &&
        failure !== null &&
        'name' in failure &&
        failure.name === 'AbortError'
      )
        return
      setError(
        failure instanceof CalculationApiError
          ? failure.message
          : 'Calculation failed. Please try again.',
      )
    } finally {
      if (pending.current === request) {
        pending.current = null
        setLoading(false)
      }
    }
  }

  function handleAction(action: CalculatorAction) {
    switch (action.type) {
      case 'input':
        edit(expressionRef.current + action.value)
        break
      case 'operation':
        edit(
          expressionRef.current +
            operations[action.operation].input +
            (action.operation === 'sqrt' ? '(' : ''),
        )
        setSelectedOperation(action.operation)
        break
      case 'clear':
        edit('')
        setSelectedOperation(null)
        break
      case 'backspace':
        edit(expressionRef.current.slice(0, -1))
        break
      case 'equals':
        void submit()
        break
    }
  }

  function reuse(value: string) {
    edit(value)
    setSelectedOperation(null)
  }

  return {
    expression,
    result,
    error,
    loading,
    selectedOperation,
    handleAction,
    reuse,
  }
}
