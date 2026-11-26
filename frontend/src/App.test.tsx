import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it } from 'vitest'
import App from './App'

afterEach(cleanup)

describe('App', () => {
  it('renders the calculator rather than the temporary component showcase', () => {
    render(<App />)
    expect(
      screen.getByRole('heading', { level: 1, name: 'Calculator' }),
    ).toBeInTheDocument()
    expect(screen.getByTestId('calculator')).toBeInTheDocument()
    expect(
      screen.queryByRole('button', { name: 'Try a button' }),
    ).not.toBeInTheDocument()
    expect(screen.getByLabelText('Result')).toHaveTextContent('—')
  })
  it('starts with the system theme and applies all selections to the document', async () => {
    const user = userEvent.setup()
    render(<App />)
    expect(screen.getByRole('radio', { name: 'System' })).toBeChecked()
    expect(document.documentElement).toHaveAttribute('data-theme', 'system')
    for (const theme of ['Light', 'Dark', 'System']) {
      await user.click(screen.getByRole('radio', { name: theme }))
      expect(document.documentElement).toHaveAttribute(
        'data-theme',
        theme.toLowerCase(),
      )
    }
  })
  it('edits text including decimal and parentheses without evaluating or parsing', async () => {
    const user = userEvent.setup()
    render(<App />)
    for (const name of [
      'Open parenthesis',
      'Digit 1',
      'Decimal point',
      'Digit 5',
      'Close parenthesis',
      'Add',
      'Digit 2',
    ]) {
      await user.click(screen.getByRole('button', { name }))
    }
    expect(screen.getByLabelText('Expression')).toHaveTextContent('(1.5)+2')
    expect(screen.getByRole('button', { name: 'Addition' })).toHaveAttribute(
      'aria-pressed',
      'true',
    )
    await user.click(screen.getByRole('button', { name: 'Equals' }))
    expect(screen.getByLabelText('Result')).toHaveTextContent('—')
    expect(
      screen.getByRole('status', { name: 'Integration notice' }),
    ).toHaveTextContent(
      'Calculation will be available when backend integration is added.',
    )
    expect(screen.getByTestId('calculator-display')).toHaveAttribute(
      'aria-busy',
      'false',
    )
    await user.click(screen.getByRole('button', { name: 'Decimal point' }))
    await user.click(screen.getByRole('button', { name: 'Decimal point' }))
    expect(screen.getByLabelText('Expression')).toHaveTextContent('(1.5)+2..')
    expect(
      screen.queryByRole('status', { name: 'Integration notice' }),
    ).not.toBeInTheDocument()
  })
  it.each([
    { label: 'Addition', symbol: '+' },
    { label: 'Subtraction', symbol: '−' },
    { label: 'Multiplication', symbol: '×' },
    { label: 'Division', symbol: '÷' },
    { label: 'Exponentiation', symbol: '^' },
    { label: 'Square root', symbol: '√' },
    { label: 'Percentage', symbol: '%' },
  ])('selects $label and appends its text token', async ({ label, symbol }) => {
    const user = userEvent.setup()
    render(<App />)
    await user.click(screen.getByRole('button', { name: label }))
    expect(screen.getByRole('button', { name: label })).toHaveAttribute(
      'aria-pressed',
      'true',
    )
    expect(screen.getByLabelText('Expression')).toHaveTextContent(symbol)
  })
  it('clears input, selection and notices and safely backspaces empty text', async () => {
    const user = userEvent.setup()
    render(<App />)
    await user.click(screen.getByRole('button', { name: 'Backspace' }))
    expect(screen.getByLabelText('Expression')).toHaveTextContent('0')
    await user.click(screen.getByRole('button', { name: 'Digit 9' }))
    await user.click(screen.getByRole('button', { name: 'Digit 8' }))
    await user.click(screen.getByRole('button', { name: 'Backspace' }))
    expect(screen.getByLabelText('Expression')).toHaveTextContent('9')
    await user.click(screen.getByRole('button', { name: 'Percentage' }))
    await user.click(screen.getByRole('button', { name: 'Equals' }))
    await user.click(screen.getByRole('button', { name: 'Clear' }))
    expect(screen.getByLabelText('Expression')).toHaveTextContent('0')
    expect(screen.getByRole('button', { name: 'Percentage' })).toHaveAttribute(
      'aria-pressed',
      'false',
    )
    expect(
      screen.queryByRole('status', { name: 'Integration notice' }),
    ).not.toBeInTheDocument()
  })
})

describe('App history', () => {
  it('starts collapsed with no invented calculations and reveals an empty panel', async () => {
    const user = userEvent.setup()
    render(<App />)
    const toggle = screen.getByRole('button', { name: 'Calculation history' })
    expect(toggle).toHaveAttribute('aria-expanded', 'false')
    expect(
      screen.queryByRole('list', { name: 'Calculation entries' }),
    ).not.toBeInTheDocument()
    await user.click(toggle)
    expect(toggle).toHaveAttribute('aria-expanded', 'true')
    expect(
      screen.getByRole('heading', { name: 'No calculations yet' }),
    ).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Clear history' })).toBeDisabled()
    await user.click(toggle)
    expect(
      screen.queryByRole('heading', { name: 'No calculations yet' }),
    ).not.toBeInTheDocument()
  })

  it('preserves calculator input and theme across toggles and never fabricates an equals result', async () => {
    const user = userEvent.setup()
    render(<App />)
    await user.click(screen.getByRole('radio', { name: 'Dark' }))
    await user.click(screen.getByRole('button', { name: 'Digit 7' }))
    await user.click(screen.getByRole('button', { name: 'Addition' }))
    await user.click(
      screen.getByRole('button', { name: 'Calculation history' }),
    )
    await user.click(screen.getByRole('button', { name: 'Equals' }))
    expect(screen.getByLabelText('Expression')).toHaveTextContent('7+')
    expect(screen.getByRole('button', { name: 'Addition' })).toHaveAttribute(
      'aria-pressed',
      'true',
    )
    expect(screen.getByRole('radio', { name: 'Dark' })).toBeChecked()
    expect(screen.getByLabelText('Result')).toHaveTextContent('—')
    expect(
      screen.getByRole('heading', { name: 'No calculations yet' }),
    ).toBeInTheDocument()
    expect(screen.queryByRole('listitem')).not.toBeInTheDocument()
  })
})
