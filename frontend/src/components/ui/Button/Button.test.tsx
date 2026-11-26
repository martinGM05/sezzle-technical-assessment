import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Button } from './Button'

afterEach(cleanup)

describe('Button', () => {
  it('renders its accessible name and defaults to a non-submitting button', () => {
    render(<Button>Continue</Button>)
    expect(screen.getByRole('button', { name: 'Continue' })).toHaveAttribute(
      'type',
      'button',
    )
  })

  it.each(['primary', 'secondary', 'danger', 'ghost'] as const)(
    'renders the %s variant with composed content',
    (variant) => {
      render(
        <Button variant={variant}>
          <span>Continue</span>
        </Button>,
      )
      expect(screen.getByRole('button', { name: 'Continue' })).toBeEnabled()
    },
  )

  it('responds to pointer and keyboard activation', async () => {
    const user = userEvent.setup()
    const onClick = vi.fn()
    render(<Button onClick={onClick}>Continue</Button>)
    await user.tab()
    expect(screen.getByRole('button')).toHaveFocus()
    await user.keyboard('{Enter} ')
    await user.click(screen.getByRole('button'))
    expect(onClick).toHaveBeenCalledTimes(3)
  })

  it('does not activate or receive tab focus when disabled', async () => {
    const user = userEvent.setup()
    const onClick = vi.fn()
    render(
      <Button disabled onClick={onClick}>
        Continue
      </Button>,
    )
    await user.click(screen.getByRole('button'))
    await user.tab()
    await user.keyboard('{Enter} ')
    expect(screen.getByRole('button')).not.toHaveFocus()
    expect(onClick).not.toHaveBeenCalled()
  })

  it('forwards native attributes and supports explicit form submission', async () => {
    const user = userEvent.setup()
    const onSubmit = vi.fn((event: React.FormEvent) => event.preventDefault())
    render(
      <form onSubmit={onSubmit}>
        <Button type="submit" name="action" value="save" title="Save changes">
          Save
        </Button>
      </form>,
    )
    const button = screen.getByRole('button', { name: 'Save' })
    expect(button).toHaveAttribute('name', 'action')
    expect(button).toHaveAttribute('value', 'save')
    expect(button).toHaveAttribute('title', 'Save changes')
    await user.click(button)
    expect(onSubmit).toHaveBeenCalledOnce()
  })
})
