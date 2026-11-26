import { cleanup, render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { HistoryPanel, type HistoryPanelProps } from './HistoryPanel'
import type { HistoryEntry } from '../../types'

const entries: readonly HistoryEntry[] = Object.freeze([
  Object.freeze({
    id: 'first',
    expression: '1 + 2',
    result: '3',
    timestamp: '2026-01-01T12:00:00Z',
  }),
  Object.freeze({
    id: 'second',
    expression: '9 ÷ 2',
    result: '4.5',
    timestamp: '2026-01-01T13:00:00Z',
  }),
])

function props(overrides: Partial<HistoryPanelProps> = {}): HistoryPanelProps {
  return {
    entries,
    open: true,
    onOpenChange: vi.fn(),
    onReuse: vi.fn(),
    onCopy: vi.fn(),
    onDelete: vi.fn(),
    onClear: vi.fn(),
    ...overrides,
  }
}

afterEach(cleanup)

describe('HistoryPanel', () => {
  it('renders a named region, semantic list and supplied entries in order', () => {
    render(<HistoryPanel {...props()} />)
    expect(
      screen.getByRole('region', { name: 'Calculation history' }),
    ).toHaveAttribute('data-testid', 'history-panel')
    expect(
      screen
        .getAllByRole('listitem')
        .map((item) => item.getAttribute('data-testid')),
    ).toEqual(['history-item-first', 'history-item-second'])
    expect(
      screen.getByRole('button', { name: 'Calculation history' }),
    ).toHaveAttribute('aria-controls', screen.getByTestId('history-content').id)
    expect(
      screen.queryByRole('group', { name: 'Clear all history?' }),
    ).not.toBeInTheDocument()
  })
  it('renders empty history without entry actions and disables clear', () => {
    render(<HistoryPanel {...props({ entries: [] })} />)
    expect(
      screen.getByRole('heading', { name: 'No calculations yet' }),
    ).toBeInTheDocument()
    expect(screen.queryByRole('list')).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Clear history' })).toBeDisabled()
  })
  it('keeps open state controlled and removes collapsed content from keyboard navigation', async () => {
    const user = userEvent.setup()
    const initial = props({ open: false })
    const { rerender } = render(<HistoryPanel {...initial} />)
    expect(screen.getByTestId('history-content')).not.toBeVisible()
    expect(screen.queryByRole('list')).not.toBeInTheDocument()
    expect(screen.getAllByRole('button')).toHaveLength(1)
    await user.tab()
    await user.keyboard('{Enter}')
    expect(initial.onOpenChange).toHaveBeenCalledExactlyOnceWith(true)
    expect(screen.queryByRole('list')).not.toBeInTheDocument()
    rerender(<HistoryPanel {...initial} open />)
    expect(screen.getByRole('list')).toBeVisible()
  })
  it('forwards exact reuse, copy and delete payloads without changing controlled entries', async () => {
    const user = userEvent.setup()
    const callbacks = props()
    render(<HistoryPanel {...callbacks} />)
    const item = within(screen.getByTestId('history-item-second'))
    await user.click(item.getByRole('button', { name: 'Reuse expression' }))
    await user.click(item.getByRole('button', { name: 'Copy result' }))
    await user.click(item.getByRole('button', { name: 'Delete entry' }))
    expect(callbacks.onReuse).toHaveBeenCalledExactlyOnceWith(entries[1])
    expect(callbacks.onCopy).toHaveBeenCalledExactlyOnceWith(entries[1])
    expect(callbacks.onDelete).toHaveBeenCalledExactlyOnceWith('second')
    expect(screen.getAllByRole('listitem')).toHaveLength(2)
  })
  it('requires explicit confirmation, initially focuses cancel, and does not trap focus', async () => {
    const user = userEvent.setup()
    const callbacks = props()
    render(<HistoryPanel {...callbacks} />)
    await user.click(screen.getByRole('button', { name: 'Clear history' }))
    expect(callbacks.onClear).not.toHaveBeenCalled()
    const confirmation = screen.getByRole('group', {
      name: 'Clear all history?',
    })
    expect(
      within(confirmation).getByRole('button', { name: 'Cancel' }),
    ).toHaveFocus()
    await user.tab()
    expect(
      within(confirmation).getByRole('button', { name: 'Clear all entries' }),
    ).toHaveFocus()
    await user.tab()
    expect(
      within(screen.getByTestId('history-item-first')).getByRole('button', {
        name: 'Reuse expression',
      }),
    ).toHaveFocus()
    await user.click(
      within(confirmation).getByRole('button', { name: 'Clear all entries' }),
    )
    expect(callbacks.onClear).toHaveBeenCalledExactlyOnceWith()
    expect(
      screen.queryByRole('group', { name: 'Clear all history?' }),
    ).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Clear history' })).toHaveFocus()
    expect(screen.getAllByRole('listitem')).toHaveLength(2)
  })
  it.each(['cancel button', 'Escape'])(
    'cancels through %s and restores focus without clearing',
    async (method) => {
      const user = userEvent.setup()
      const callbacks = props()
      render(<HistoryPanel {...callbacks} />)
      await user.click(screen.getByRole('button', { name: 'Clear history' }))
      if (method === 'Escape') await user.keyboard('{Escape}')
      else await user.click(screen.getByRole('button', { name: 'Cancel' }))
      expect(callbacks.onClear).not.toHaveBeenCalled()
      expect(
        screen.queryByRole('group', { name: 'Clear all history?' }),
      ).not.toBeInTheDocument()
      expect(
        screen.getByRole('button', { name: 'Clear history' }),
      ).toHaveFocus()
    },
  )
  it('ignores unrelated confirmation keys', async () => {
    const user = userEvent.setup()
    render(<HistoryPanel {...props()} />)
    await user.click(screen.getByRole('button', { name: 'Clear history' }))
    await user.keyboard('x')
    expect(
      screen.getByRole('group', { name: 'Clear all history?' }),
    ).toBeVisible()
  })
  it('restores toggle focus on programmatic collapse and resets pending confirmation', async () => {
    const user = userEvent.setup()
    const initial = props()
    const { rerender } = render(<HistoryPanel {...initial} />)
    await user.click(screen.getByRole('button', { name: 'Clear history' }))
    rerender(<HistoryPanel {...initial} open={false} />)
    expect(
      screen.getByRole('button', { name: 'Calculation history' }),
    ).toHaveFocus()
    expect(
      screen.queryByRole('button', { name: 'Cancel' }),
    ).not.toBeInTheDocument()
    rerender(<HistoryPanel {...initial} />)
    expect(
      screen.queryByRole('group', { name: 'Clear all history?' }),
    ).not.toBeInTheDocument()
    expect(initial.onClear).not.toHaveBeenCalled()
  })
  it('does not steal outside focus when the parent collapses the panel', async () => {
    const user = userEvent.setup()
    const initial = props()
    const { rerender } = render(
      <>
        <button>Outside</button>
        <HistoryPanel {...initial} />
      </>,
    )
    await user.click(screen.getByRole('button', { name: 'Clear history' }))
    await user.click(screen.getByRole('button', { name: 'Outside' }))
    rerender(
      <>
        <button>Outside</button>
        <HistoryPanel {...initial} open={false} />
      </>,
    )
    expect(screen.getByRole('button', { name: 'Outside' })).toHaveFocus()
  })
  it('safely removes a pending confirmation when controlled entries become empty', async () => {
    const user = userEvent.setup()
    const initial = props()
    const { rerender } = render(<HistoryPanel {...initial} />)
    await user.click(screen.getByRole('button', { name: 'Clear history' }))
    rerender(<HistoryPanel {...initial} entries={[]} />)
    expect(
      screen.queryByRole('button', { name: 'Clear all entries' }),
    ).not.toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: 'Calculation history' }),
    ).toHaveFocus()
    expect(initial.onClear).not.toHaveBeenCalled()
    rerender(<HistoryPanel {...initial} />)
    expect(
      screen.queryByRole('group', { name: 'Clear all history?' }),
    ).not.toBeInTheDocument()
  })
  it('shows large lists in twenty-entry increments without losing items or order', async () => {
    const user = userEvent.setup()
    const largeEntries = Array.from({ length: 45 }, (_, index) => ({
      id: String(index),
      expression: `${index} + 1`,
      result: String(index + 1),
      timestamp: '2026-01-01T12:00:00Z',
    }))
    render(<HistoryPanel {...props({ entries: largeEntries })} />)
    expect(screen.getAllByRole('listitem')).toHaveLength(20)
    expect(screen.getByText('Showing 20 of 45 entries')).toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: 'Show more history' }))
    expect(screen.getAllByRole('listitem')).toHaveLength(40)
    await user.click(screen.getByRole('button', { name: 'Show more history' }))
    expect(screen.getAllByRole('listitem')).toHaveLength(45)
    expect(screen.getAllByRole('listitem').at(-1)).toHaveAttribute(
      'data-testid',
      'history-item-44',
    )
    expect(
      screen.queryByRole('button', { name: 'Show more history' }),
    ).not.toBeInTheDocument()
  })
  it('updates displayed entries from props and keeps independent panel identifiers', () => {
    const initial = props()
    const { rerender } = render(<HistoryPanel {...initial} />)
    rerender(<HistoryPanel {...initial} entries={entries.slice(1)} />)
    expect(screen.queryByTestId('history-item-first')).not.toBeInTheDocument()
    expect(screen.getAllByRole('listitem')).toHaveLength(1)
    render(<HistoryPanel {...props({ entries: [], open: false })} />)
    const controls = screen
      .getAllByRole('button', { name: 'Calculation history' })
      .map((button) => button.getAttribute('aria-controls'))
    expect(new Set(controls).size).toBe(2)
  })
})
