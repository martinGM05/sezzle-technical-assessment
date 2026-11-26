import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { CalculatorKey } from './CalculatorKey'

afterEach(cleanup)

describe('CalculatorKey', () => {
  it('renders a native named button and delivers its exact action', async () => {
    const user = userEvent.setup()
    const onAction = vi.fn()
    render(
      <CalculatorKey
        label="Seven"
        action={{ type: 'input', value: '7' }}
        onAction={onAction}
        testId="key-7"
      >
        7
      </CalculatorKey>,
    )
    const key = screen.getByRole('button', { name: 'Seven' })
    expect(key).toHaveAttribute('type', 'button')
    expect(key).toHaveTextContent('7')
    expect(key).toBe(screen.getByTestId('key-7'))
    await user.click(key)
    expect(onAction).toHaveBeenCalledExactlyOnceWith({
      type: 'input',
      value: '7',
    })
  })
  it('supports Tab, Enter, and Space without custom keyboard handlers', async () => {
    const user = userEvent.setup()
    const onAction = vi.fn()
    render(
      <CalculatorKey
        label="Equals"
        action={{ type: 'equals' }}
        onAction={onAction}
      >
        =
      </CalculatorKey>,
    )
    await user.tab()
    expect(screen.getByRole('button', { name: 'Equals' })).toHaveFocus()
    await user.keyboard('{Enter} ')
    expect(onAction).toHaveBeenCalledTimes(2)
    expect(onAction).toHaveBeenLastCalledWith({ type: 'equals' })
  })
  it('cannot activate or receive sequential focus when disabled', async () => {
    const user = userEvent.setup()
    const onAction = vi.fn()
    render(
      <CalculatorKey
        label="Clear"
        action={{ type: 'clear' }}
        onAction={onAction}
        disabled
      >
        C
      </CalculatorKey>,
    )
    const key = screen.getByRole('button', { name: 'Clear' })
    expect(key).toBeDisabled()
    await user.click(key)
    await user.tab()
    expect(key).not.toHaveFocus()
    expect(onAction).not.toHaveBeenCalled()
  })
})
