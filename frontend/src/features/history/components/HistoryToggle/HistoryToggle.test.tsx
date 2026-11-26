import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { HistoryToggle } from './HistoryToggle'

afterEach(cleanup)

describe('HistoryToggle', () => {
  it.each([false, true])(
    'reflects controlled open=%s and requests the opposite state',
    async (open) => {
      const user = userEvent.setup()
      const onOpenChange = vi.fn()
      const { rerender } = render(
        <HistoryToggle
          open={open}
          panelId="history-content"
          onOpenChange={onOpenChange}
        />,
      )
      const button = screen.getByRole('button', { name: 'Calculation history' })
      expect(button).toHaveAttribute('aria-expanded', String(open))
      expect(button).toHaveAttribute('aria-controls', 'history-content')
      expect(button).toHaveAttribute('data-testid', 'history-toggle')
      await user.click(button)
      expect(onOpenChange).toHaveBeenCalledExactlyOnceWith(!open)
      expect(button).toHaveAttribute('aria-expanded', String(open))
      rerender(
        <HistoryToggle
          open={!open}
          panelId="history-content"
          onOpenChange={onOpenChange}
        />,
      )
      expect(button).toHaveAttribute('aria-expanded', String(!open))
    },
  )
  it.each(['{Enter}', ' '])(
    'supports native keyboard activation with %s',
    async (key) => {
      const user = userEvent.setup()
      const onOpenChange = vi.fn()
      render(
        <HistoryToggle
          open={false}
          panelId="history-content"
          onOpenChange={onOpenChange}
        />,
      )
      await user.tab()
      expect(screen.getByRole('button')).toHaveFocus()
      await user.keyboard(key)
      expect(onOpenChange).toHaveBeenCalledExactlyOnceWith(true)
    },
  )
})
