import { cleanup, render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { CalculatorKeypad } from './CalculatorKeypad'

afterEach(cleanup)

describe('CalculatorKeypad', () => {
  it('exposes every digit and sends the exact input', async () => {
    const user = userEvent.setup()
    const onAction = vi.fn()
    render(<CalculatorKeypad onAction={onAction} />)
    expect(
      within(
        screen.getByRole('group', { name: 'Calculator keypad' }),
      ).getAllByRole('button'),
    ).toHaveLength(20)
    for (const digit of '0123456789') {
      await user.click(screen.getByRole('button', { name: `Digit ${digit}` }))
      expect(onAction).toHaveBeenLastCalledWith({ type: 'input', value: digit })
    }
    expect(onAction).toHaveBeenCalledTimes(10)
  })
  it.each([
    { label: 'Decimal point', action: { type: 'input', value: '.' } },
    { label: 'Open parenthesis', action: { type: 'input', value: '(' } },
    { label: 'Close parenthesis', action: { type: 'input', value: ')' } },
    { label: 'Add', action: { type: 'operation', operation: 'add' } },
    { label: 'Subtract', action: { type: 'operation', operation: 'subtract' } },
    { label: 'Multiply', action: { type: 'operation', operation: 'multiply' } },
    { label: 'Divide', action: { type: 'operation', operation: 'divide' } },
    { label: 'Clear', action: { type: 'clear' } },
    { label: 'Backspace', action: { type: 'backspace' } },
    { label: 'Equals', action: { type: 'equals' } },
  ])('dispatches $label without evaluating', async ({ label, action }) => {
    const user = userEvent.setup()
    const onAction = vi.fn()
    render(<CalculatorKeypad onAction={onAction} />)
    await user.click(screen.getByRole('button', { name: label }))
    expect(onAction).toHaveBeenCalledExactlyOnceWith(action)
  })
  it('disables all controls together and suppresses actions', async () => {
    const user = userEvent.setup()
    const onAction = vi.fn()
    render(<CalculatorKeypad disabled onAction={onAction} />)
    for (const button of screen.getAllByRole('button')) {
      expect(button).toBeDisabled()
      await user.click(button)
    }
    expect(onAction).not.toHaveBeenCalled()
  })
  it('preserves row order for keyboard navigation', async () => {
    const user = userEvent.setup()
    render(<CalculatorKeypad onAction={vi.fn()} />)
    for (const label of [
      'Clear',
      'Backspace',
      'Open parenthesis',
      'Close parenthesis',
      'Digit 7',
    ]) {
      await user.tab()
      expect(screen.getByRole('button', { name: label })).toHaveFocus()
    }
  })
})
