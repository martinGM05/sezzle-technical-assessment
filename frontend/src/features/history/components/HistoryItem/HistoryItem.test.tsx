import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { HistoryItem } from './HistoryItem'

const entry = Object.freeze({
  id: 'calculation-1',
  expression: '1 ÷ 3',
  result: '0.3333333333333333',
  timestamp: '2026-06-01T12:30:00.000Z',
})

function setup(overrides = {}) {
  const callbacks = { onReuse: vi.fn(), onCopy: vi.fn(), onDelete: vi.fn() }
  render(
    <ul>
      <HistoryItem entry={entry} {...callbacks} {...overrides} />
    </ul>,
  )
  return callbacks
}

afterEach(cleanup)

describe('HistoryItem', () => {
  it('preserves the expression, full result and machine-readable timestamp', () => {
    setup()
    expect(screen.getByRole('listitem')).toHaveAttribute(
      'data-testid',
      'history-item-calculation-1',
    )
    expect(screen.getByText(entry.expression)).toBeInTheDocument()
    expect(screen.getByText(entry.result)).toBeInTheDocument()
    const time = screen.getByTestId('history-timestamp')
    expect(time.tagName).toBe('TIME')
    expect(time).toHaveAttribute('datetime', entry.timestamp)
    expect(time.textContent).not.toBe('')
  })
  it.each(['invalid', ''])(
    'shows an honest fallback for timestamp %j',
    (timestamp) => {
      setup({ entry: { ...entry, timestamp } })
      expect(screen.getByText('Time unavailable')).toBeInTheDocument()
      expect(screen.getByText('Time unavailable')).not.toHaveAttribute(
        'datetime',
      )
    },
  )
  it('emits exact controlled callbacks without mutating the entry or announcing a copy', async () => {
    const user = userEvent.setup()
    const callbacks = setup()
    await user.click(screen.getByRole('button', { name: 'Reuse expression' }))
    await user.click(screen.getByRole('button', { name: 'Copy result' }))
    await user.click(screen.getByRole('button', { name: 'Delete entry' }))
    expect(callbacks.onReuse).toHaveBeenCalledExactlyOnceWith(entry)
    expect(callbacks.onCopy).toHaveBeenCalledExactlyOnceWith(entry)
    expect(callbacks.onDelete).toHaveBeenCalledExactlyOnceWith(entry.id)
    expect(screen.getByText(entry.result)).toBeInTheDocument()
    expect(screen.queryByRole('status')).not.toBeInTheDocument()
  })
  it('supports keyboard activation for every action', async () => {
    const user = userEvent.setup()
    const callbacks = setup()
    await user.tab()
    expect(
      screen.getByRole('button', { name: 'Reuse expression' }),
    ).toHaveFocus()
    await user.keyboard('{Enter}')
    await user.tab()
    await user.keyboard(' ')
    await user.tab()
    await user.keyboard('{Enter}')
    expect(callbacks.onReuse).toHaveBeenCalledOnce()
    expect(callbacks.onCopy).toHaveBeenCalledOnce()
    expect(callbacks.onDelete).toHaveBeenCalledOnce()
  })
  it('disables copying when no integration callback is available', async () => {
    const user = userEvent.setup()
    setup({ onCopy: undefined })
    expect(screen.getByRole('button', { name: 'Copy result' })).toBeDisabled()
    await user.tab()
    await user.tab()
    expect(screen.getByRole('button', { name: 'Delete entry' })).toHaveFocus()
  })
  it('renders long expression and result text without truncating their values', () => {
    const expression = '1234567890+'.repeat(100)
    const result = '1234567890'.repeat(100)
    setup({ entry: { ...entry, expression, result } })
    expect(screen.getByText(expression)).toHaveTextContent(expression)
    expect(screen.getByText(result)).toHaveTextContent(result)
  })
})
