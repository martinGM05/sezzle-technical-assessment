import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { HistoryEmptyState } from './HistoryEmptyState'

afterEach(cleanup)

describe('HistoryEmptyState', () => {
  it('explains empty history without invented entries or actions', () => {
    render(<HistoryEmptyState />)
    expect(
      screen.getByRole('heading', { name: 'No calculations yet' }),
    ).toBeInTheDocument()
    expect(
      screen.getByText('Your completed calculations will appear here.'),
    ).toBeInTheDocument()
    expect(screen.getByTestId('history-empty-state')).toBeInTheDocument()
    expect(screen.queryByRole('list')).not.toBeInTheDocument()
    expect(screen.queryByRole('button')).not.toBeInTheDocument()
  })
})
