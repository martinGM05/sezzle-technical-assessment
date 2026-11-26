import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { CalculatorDisplay } from './CalculatorDisplay'

afterEach(cleanup)

describe('CalculatorDisplay', () => {
  it('labels the expression and shows empty placeholders without inventing a result', () => {
    render(<CalculatorDisplay expression="" result={null} />)
    expect(screen.getByLabelText('Expression')).toHaveTextContent('0')
    expect(screen.getByLabelText('Result')).toHaveTextContent('—')
  })
  it('preserves zero and the complete expression', () => {
    render(<CalculatorDisplay expression="1 − 1" result="0" />)
    expect(screen.getByLabelText('Expression')).toHaveTextContent('1 − 1')
    expect(screen.getByLabelText('Result')).toHaveTextContent('0')
  })
  it('does not truncate long values', () => {
    const expression = '1234567890'.repeat(30)
    const result = '9876543210'.repeat(30)
    render(<CalculatorDisplay expression={expression} result={result} />)
    expect(screen.getByLabelText('Expression').textContent).toBe(expression)
    expect(screen.getByLabelText('Result').textContent).toBe(result)
  })
  it('announces loading and marks the display busy', () => {
    render(<CalculatorDisplay expression="2 + 2" result={null} loading />)
    expect(screen.getByTestId('calculator-display')).toHaveAttribute(
      'aria-busy',
      'true',
    )
    expect(
      screen.getByRole('status', { name: 'Calculation status' }),
    ).toHaveTextContent('Calculating…')
  })
  it('renders the supplied error as an alert', () => {
    render(
      <CalculatorDisplay
        expression="1 ÷ 0"
        result={null}
        error="Cannot divide by zero."
      />,
    )
    expect(screen.getByRole('alert')).toHaveTextContent(
      'Cannot divide by zero.',
    )
  })
  it('updates from controlled props and clears stale error markup', () => {
    const { rerender } = render(
      <CalculatorDisplay expression="1" result={null} error="Invalid input." />,
    )
    rerender(<CalculatorDisplay expression="2" result="2" />)
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
    expect(screen.getByLabelText('Result')).toHaveTextContent('2')
  })
})
