import { useState } from 'react'
import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { ThemeToggle, type Theme } from './ThemeToggle'

afterEach(cleanup)

function ControlledToggle() {
  const [theme, setTheme] = useState<Theme>('system')
  return <ThemeToggle value={theme} onChange={setTheme} />
}

describe('ThemeToggle', () => {
  it.each(['light', 'dark', 'system'] as const)(
    'represents controlled %s selection with accessible labels',
    (value) => {
      render(<ThemeToggle value={value} onChange={vi.fn()} />)
      expect(
        screen.getByRole('group', { name: 'Color theme' }),
      ).toHaveAttribute('data-testid', 'theme-toggle')
      for (const theme of ['light', 'dark', 'system']) {
        const radio = screen.getByRole('radio', {
          name: new RegExp(theme, 'i'),
        })
        expect(radio).toHaveAttribute('data-testid', `theme-option-${theme}`)
        if (theme === value) expect(radio).toBeChecked()
        else expect(radio).not.toBeChecked()
      }
    },
  )
  it('requests a selection without owning the parent state', async () => {
    const user = userEvent.setup()
    const onChange = vi.fn()
    const { rerender } = render(
      <ThemeToggle value="system" onChange={onChange} />,
    )
    await user.click(screen.getByRole('radio', { name: 'Dark' }))
    expect(onChange).toHaveBeenCalledExactlyOnceWith('dark')
    expect(screen.getByRole('radio', { name: 'System' })).toBeChecked()
    rerender(<ThemeToggle value="dark" onChange={onChange} />)
    expect(screen.getByRole('radio', { name: 'Dark' })).toBeChecked()
  })
  it('supports tab, arrows, space and all three choices', async () => {
    const user = userEvent.setup()
    render(<ControlledToggle />)
    await user.tab()
    expect(screen.getByRole('radio', { name: 'System' })).toHaveFocus()
    await user.keyboard('{ArrowRight}')
    expect(screen.getByRole('radio', { name: 'Light' })).toBeChecked()
    await user.keyboard('{ArrowRight}')
    expect(screen.getByRole('radio', { name: 'Dark' })).toBeChecked()
    await user.keyboard('{ArrowRight} ')
    expect(screen.getByRole('radio', { name: 'System' })).toBeChecked()
  })
  it('keeps separate instances in independent native radio groups', async () => {
    const user = userEvent.setup()
    render(
      <>
        <ControlledToggle />
        <ControlledToggle />
      </>,
    )
    await user.click(screen.getAllByRole('radio', { name: 'Light' })[0])
    expect(screen.getAllByRole('radio', { name: 'Light' })[0]).toBeChecked()
    expect(screen.getAllByRole('radio', { name: 'System' })[1]).toBeChecked()
  })
})
