import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { Alert } from './Alert'

afterEach(cleanup)

describe('Alert', () => {
  it.each([
    { variant: 'info', role: 'status' },
    { variant: 'success', role: 'status' },
    { variant: 'warning', role: 'status' },
    { variant: 'error', role: 'alert' },
  ] as const)(
    'renders $variant content with $role semantics',
    ({ variant, role }) => {
      render(
        <Alert variant={variant}>
          <strong>Important message</strong>
        </Alert>,
      )
      expect(screen.getByRole(role)).toHaveTextContent('Important message')
    },
  )
  it('defaults to a polite informational status', () => {
    render(<Alert>Ready</Alert>)
    expect(screen.getByRole('status')).toHaveTextContent('Ready')
  })
  it('allows explicit urgency and updates its message', () => {
    const { rerender } = render(
      <Alert variant="warning" role="alert">
        Connection lost
      </Alert>,
    )
    expect(screen.getByRole('alert')).toHaveTextContent('Connection lost')
    rerender(
      <Alert variant="success" role="status">
        Connected
      </Alert>,
    )
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
    expect(screen.getByRole('status')).toHaveTextContent('Connected')
  })
})
