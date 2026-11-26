import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { useCalculatorKeyboard } from './useCalculatorKeyboard'
import type { CalculatorAction } from '../types'

afterEach(cleanup)

function Harness({
  onAction,
}: {
  onAction: (action: CalculatorAction) => void
}) {
  useCalculatorKeyboard(onAction)
  return (
    <>
      <section data-testid="calculator">
        <button>Digit</button>
      </section>
      <input aria-label="Other input" />
      <textarea aria-label="Notes" />
      <select aria-label="Choice">
        <option>One</option>
      </select>
      <div contentEditable suppressContentEditableWarning aria-label="Editor">
        Edit
      </div>
      <button>History action</button>
      <input type="radio" aria-label="Theme" />
    </>
  )
}

// Matrix: digits, syntax keys, all operations, editing, native activation,
// editable/other controls, modifiers/IME/repeat, listener replacement and cleanup.
describe('useCalculatorKeyboard', () => {
  it.each([
    { key: '7', action: { type: 'input', value: '7' } },
    { key: '.', action: { type: 'input', value: '.' } },
    { key: '(', action: { type: 'input', value: '(' } },
    { key: ')', action: { type: 'input', value: ')' } },
    { key: '+', action: { type: 'operation', operation: 'add' } },
    { key: '-', action: { type: 'operation', operation: 'subtract' } },
    { key: '*', action: { type: 'operation', operation: 'multiply' } },
    { key: '/', action: { type: 'operation', operation: 'divide' } },
    { key: '^', action: { type: 'operation', operation: 'power' } },
    { key: '%', action: { type: 'operation', operation: 'percentage' } },
    { key: 'r', action: { type: 'operation', operation: 'sqrt' } },
    { key: 'Enter', action: { type: 'equals' } },
    { key: '=', action: { type: 'equals' } },
    { key: 'Backspace', action: { type: 'backspace' } },
    { key: 'Escape', action: { type: 'clear' } },
    { key: 'Delete', action: { type: 'clear' } },
  ])('dispatches $key without evaluating', ({ key, action }) => {
    const onAction = vi.fn()
    render(<Harness onAction={onAction} />)
    expect(fireEvent.keyDown(document.body, { key })).toBe(false)
    expect(onAction).toHaveBeenCalledExactlyOnceWith(action)
  })

  it('evaluates on button Enter without a native click, preserving Space activation', async () => {
    const user = userEvent.setup()
    const onAction = vi.fn()
    render(<Harness onAction={onAction} />)
    const button = screen.getByRole('button', { name: 'Digit' })
    const clicked = vi.fn()
    button.addEventListener('click', clicked)
    button.focus()
    await user.keyboard('{Enter}')
    expect(clicked).not.toHaveBeenCalled()
    expect(onAction).toHaveBeenCalledExactlyOnceWith({ type: 'equals' })
    await user.keyboard(' ')
    expect(clicked).toHaveBeenCalledTimes(1)
    expect(onAction).toHaveBeenCalledTimes(1)
    // Held Enter must suppress both repeat evaluation and native activation.
    expect(fireEvent.keyDown(button, { key: 'Enter', repeat: true })).toBe(
      false,
    )
    expect(onAction).toHaveBeenCalledTimes(1)
  })

  it('preserves native Enter and Space on controls outside the calculator', async () => {
    const user = userEvent.setup()
    const onAction = vi.fn()
    render(<Harness onAction={onAction} />)
    const clicked = vi.fn()
    const button = screen.getByRole('button', { name: 'History action' })
    button.addEventListener('click', clicked)
    button.focus()
    await user.keyboard('{Enter} ')
    expect(clicked).toHaveBeenCalledTimes(2)
    const theme = screen.getByRole('radio', { name: 'Theme' })
    theme.focus()
    await user.keyboard(' ')
    expect(theme).toBeChecked()
    expect(onAction).not.toHaveBeenCalled()
  })

  it.each(['Other input', 'Notes', 'Choice', 'Editor', 'Theme'])(
    'does not intercept %s',
    (name) => {
      const onAction = vi.fn()
      render(<Harness onAction={onAction} />)
      expect(fireEvent.keyDown(screen.getByLabelText(name), { key: '7' })).toBe(
        true,
      )
      expect(onAction).not.toHaveBeenCalled()
    },
  )

  it('does not hijack history controls or composing/browser shortcuts', () => {
    const onAction = vi.fn()
    render(<Harness onAction={onAction} />)
    fireEvent.keyDown(screen.getByRole('button', { name: 'History action' }), {
      key: 'Escape',
    })
    for (const options of [
      { ctrlKey: true },
      { metaKey: true },
      { altKey: true },
      { isComposing: true },
    ]) {
      expect(
        fireEvent.keyDown(document.body, { key: 'Enter', ...options }),
      ).toBe(true)
    }
    expect(fireEvent.keyDown(document.body, { key: 'F5' })).toBe(true)
    expect(onAction).not.toHaveBeenCalled()
  })

  it('uses the latest callback and removes its listener on unmount', () => {
    const old = vi.fn()
    const latest = vi.fn()
    const { rerender, unmount } = render(<Harness onAction={old} />)
    rerender(<Harness onAction={latest} />)
    fireEvent.keyDown(document.body, { key: '2' })
    expect(old).not.toHaveBeenCalled()
    expect(latest).toHaveBeenCalledTimes(1)
    unmount()
    fireEvent.keyDown(document.body, { key: '3' })
    expect(latest).toHaveBeenCalledTimes(1)
  })
})
