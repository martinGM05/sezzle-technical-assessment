import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Calculator } from './Calculator'

afterEach(cleanup)

const initial = { expression: '', result: null, selectedOperation: null }

describe('Calculator', () => {
  it('composes display, keypad and operation selector without calculating', async () => {
    const user = userEvent.setup()
    const onAction = vi.fn()
    render(<Calculator {...initial} onAction={onAction} />)
    expect(
      screen.getByRole('heading', { name: 'Calculator', level: 1 }),
    ).toBeInTheDocument()
    expect(screen.getByTestId('calculator-keypad')).toBeInTheDocument()
    expect(screen.getByTestId('operation-selector')).toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: 'Digit 2' }))
    await user.click(screen.getByRole('button', { name: 'Square root' }))
    await user.click(screen.getByRole('button', { name: 'Equals' }))
    expect(onAction.mock.calls).toEqual([
      [{ type: 'input', value: '2' }],
      [{ type: 'operation', operation: 'sqrt' }],
      [{ type: 'equals' }],
    ])
    expect(screen.getByLabelText('Expression')).toHaveTextContent('0')
    expect(screen.getByLabelText('Result')).toHaveTextContent('—')
  })
  it.each([
    { state: 'loading', loading: true, disabled: false },
    { state: 'disabled', loading: false, disabled: true },
  ])('blocks all controls when $state', async ({ loading, disabled }) => {
    const user = userEvent.setup()
    const onAction = vi.fn()
    render(
      <Calculator
        {...initial}
        onAction={onAction}
        loading={loading}
        disabled={disabled}
      />,
    )
    for (const button of screen.getAllByRole('button')) {
      expect(button).toBeDisabled()
      await user.click(button)
    }
    expect(onAction).not.toHaveBeenCalled()
    expect(screen.getByTestId('calculator-display')).toHaveAttribute(
      'aria-busy',
      String(loading),
    )
  })
  it('updates controlled values, selected operations and recovery states', () => {
    const { rerender } = render(
      <Calculator
        expression="1÷0"
        result={null}
        selectedOperation="divide"
        error="Cannot divide by zero."
        onAction={vi.fn()}
      />,
    )
    expect(screen.getByRole('alert')).toHaveTextContent(
      'Cannot divide by zero.',
    )
    expect(screen.getByRole('button', { name: 'Division' })).toHaveAttribute(
      'aria-pressed',
      'true',
    )
    rerender(
      <Calculator
        expression="2×0"
        result="0"
        selectedOperation="multiply"
        onAction={vi.fn()}
      />,
    )
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
    expect(screen.getByLabelText('Expression')).toHaveTextContent('2×0')
    expect(screen.getByLabelText('Result')).toHaveTextContent('0')
    expect(
      screen.getByRole('button', { name: 'Multiplication' }),
    ).toHaveAttribute('aria-pressed', 'true')
  })
  it('preserves long expressions and results through composition', () => {
    const value = '1234567890'.repeat(50)
    render(
      <Calculator
        {...initial}
        expression={value}
        result={value}
        onAction={vi.fn()}
      />,
    )
    expect(screen.getByLabelText('Expression').textContent).toBe(value)
    expect(screen.getByLabelText('Result').textContent).toBe(value)
  })
})
