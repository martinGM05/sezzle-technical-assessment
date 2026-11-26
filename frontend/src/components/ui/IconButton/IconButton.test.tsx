import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { IconButton } from './IconButton'

afterEach(cleanup)
const icon = (
  <svg data-testid="sample-icon">
    <title>Decorative star</title>
    <path d="M1 1h10" />
  </svg>
)

describe('IconButton', () => {
  it('names the button without exposing its decorative icon', () => {
    render(<IconButton aria-label="Favorite">{icon}</IconButton>)
    expect(screen.getByRole('button', { name: 'Favorite' })).toBeEnabled()
    expect(screen.getByTestId('sample-icon')).toBeInTheDocument()
    expect(screen.queryByRole('img')).not.toBeInTheDocument()
    expect(
      screen.getByTestId('sample-icon').closest('[aria-hidden="true"]'),
    ).not.toBeNull()
  })
  it('supports click, Tab, Enter and Space', async () => {
    const user = userEvent.setup()
    const onClick = vi.fn()
    render(
      <IconButton aria-label="Favorite" onClick={onClick}>
        {icon}
      </IconButton>,
    )
    await user.tab()
    expect(screen.getByRole('button')).toHaveFocus()
    await user.keyboard('{Enter} ')
    await user.click(screen.getByRole('button'))
    expect(onClick).toHaveBeenCalledTimes(3)
  })
  it('prevents disabled activation and tab focus', async () => {
    const user = userEvent.setup()
    const onClick = vi.fn()
    render(
      <IconButton aria-label="Favorite" disabled onClick={onClick}>
        {icon}
      </IconButton>,
    )
    await user.click(screen.getByRole('button'))
    await user.tab()
    await user.keyboard('{Enter} ')
    expect(screen.getByRole('button')).not.toHaveFocus()
    expect(onClick).not.toHaveBeenCalled()
  })
})
