import { Alert } from '../../../../components/ui/Alert/Alert'
import './CalculatorDisplay.css'

export type CalculatorDisplayProps = {
  expression: string
  result: string | null
  loading?: boolean
  error?: string | null
}

export function CalculatorDisplay({
  expression,
  result,
  loading = false,
  error = null,
}: CalculatorDisplayProps) {
  return (
    <div
      className="calculator-display"
      data-testid="calculator-display"
      aria-busy={loading}
    >
      <div
        className="calculator-display__expression"
        aria-label="Expression"
        data-testid="calculator-expression"
      >
        {expression || '0'}
      </div>
      <output
        className="calculator-display__result"
        aria-label="Result"
        data-testid="calculator-result"
      >
        {result ?? '—'}
      </output>
      {loading ? (
        <p
          className="calculator-display__status"
          role="status"
          aria-label="Calculation status"
        >
          Calculating…
        </p>
      ) : null}
      {error ? <Alert variant="error">{error}</Alert> : null}
    </div>
  )
}
