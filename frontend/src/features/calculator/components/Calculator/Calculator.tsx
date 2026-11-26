import {
  CalculatorDisplay,
  type CalculatorDisplayProps,
} from '../CalculatorDisplay/CalculatorDisplay'
import { CalculatorKeypad } from '../CalculatorKeypad/CalculatorKeypad'
import { OperationSelector } from '../OperationSelector/OperationSelector'
import type { CalculatorAction, Operation } from '../../types'
import './Calculator.css'

export type CalculatorProps = CalculatorDisplayProps & {
  selectedOperation: Operation | null
  onAction: (action: CalculatorAction) => void
  disabled?: boolean
}

export function Calculator({
  expression,
  result,
  selectedOperation,
  onAction,
  loading = false,
  error = null,
  disabled = false,
}: CalculatorProps) {
  const controlsDisabled = disabled || loading
  return (
    <section
      className="calculator"
      aria-label="Calculator"
      data-testid="calculator"
    >
      <header className="calculator__header">
        <span className="calculator__badge" aria-hidden="true">
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="1.8"
            focusable="false"
          >
            <rect x="5" y="2" width="14" height="20" rx="3" />
            <path
              d="M8 6h8M8 11h1m3 0h1m3 0h0M8 15h1m3 0h1m3 0h0M8 19h1m3 0h1m3 0h0"
              strokeLinecap="round"
            />
          </svg>
        </span>
        <div className="calculator__intro">
          <h1>Calculator</h1>
          <p>Seven operations, one place</p>
        </div>
      </header>
      <div className="calculator__layout">
        <div className="calculator__display">
          <CalculatorDisplay
            expression={expression}
            result={result}
            loading={loading}
            error={error}
          />
        </div>
        <div className="calculator__operations">
          <OperationSelector
            selectedOperation={selectedOperation}
            disabled={controlsDisabled}
            onSelect={(operation) => onAction({ type: 'operation', operation })}
          />
        </div>
        <div className="calculator__keypad">
          <CalculatorKeypad onAction={onAction} disabled={controlsDisabled} />
        </div>
      </div>
    </section>
  )
}
