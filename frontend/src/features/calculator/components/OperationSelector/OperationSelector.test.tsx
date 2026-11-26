import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { OperationSelector } from './OperationSelector'

afterEach(cleanup)

describe('OperationSelector', () => {
  it.each([
    { label: 'Addition', operation: 'add' },
    { label: 'Subtraction', operation: 'subtract' },
    { label: 'Multiplication', operation: 'multiply' },
    { label: 'Division', operation: 'divide' },
    { label: 'Exponentiation', operation: 'power' },
    { label: 'Square root', operation: 'sqrt' },
    { label: 'Percentage', operation: 'percentage' },
  ])(
    'offers $label and reports its operation',
    async ({ label, operation }) => {
      const user = userEvent.setup()
      const onSelect = vi.fn()
      render(<OperationSelector selectedOperation={null} onSelect={onSelect} />)
      expect(screen.getAllByRole('button')).toHaveLength(7)
      await user.click(screen.getByRole('button', { name: label }))
      expect(onSelect).toHaveBeenCalledExactlyOnceWith(operation)
      expect(screen.getByRole('button', { name: label })).toHaveAttribute(
        'aria-pressed',
        'false',
      )
    },
  )
  it('only changes selection when controlled props change', () => {
    const { rerender } = render(
      <OperationSelector selectedOperation="add" onSelect={vi.fn()} />,
    )
    expect(screen.getByRole('button', { name: 'Addition' })).toHaveAttribute(
      'aria-pressed',
      'true',
    )
    rerender(<OperationSelector selectedOperation="sqrt" onSelect={vi.fn()} />)
    expect(screen.getByRole('button', { name: 'Addition' })).toHaveAttribute(
      'aria-pressed',
      'false',
    )
    expect(screen.getByRole('button', { name: 'Square root' })).toHaveAttribute(
      'aria-pressed',
      'true',
    )
  })
  it('supports keyboard activation', async () => {
    const user = userEvent.setup()
    const onSelect = vi.fn()
    render(<OperationSelector selectedOperation={null} onSelect={onSelect} />)
    await user.tab()
    expect(screen.getByRole('button', { name: 'Addition' })).toHaveFocus()
    await user.keyboard('{Enter} ')
    expect(onSelect).toHaveBeenCalledTimes(2)
    expect(onSelect).toHaveBeenLastCalledWith('add')
  })
  it('disables every operation and suppresses callbacks', async () => {
    const user = userEvent.setup()
    const onSelect = vi.fn()
    render(
      <OperationSelector
        selectedOperation="divide"
        onSelect={onSelect}
        disabled
      />,
    )
    for (const button of screen.getAllByRole('button')) {
      expect(button).toBeDisabled()
      await user.click(button)
    }
    expect(onSelect).not.toHaveBeenCalled()
  })
})
