import { cleanup, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import App from './App'

beforeEach(() => {
  localStorage.clear()
  vi.stubGlobal(
    'fetch',
    vi.fn().mockImplementation(() =>
      Promise.resolve(
        new Response('{"result":3.5}', {
          headers: { 'Content-Type': 'application/json' },
        }),
      ),
    ),
  )
})
afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
})

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
  it('edits text including decimal and parentheses and displays only the server result', async () => {
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
    await waitFor(() =>
      expect(screen.getByLabelText('Result')).toHaveTextContent('3.5'),
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

  it('preserves input and theme and does not add a failed calculation to history', async () => {
    const user = userEvent.setup()
    vi.mocked(fetch).mockResolvedValueOnce(
      new Response(
        JSON.stringify({
          type: 'about:blank',
          title: 'Unprocessable Content',
          status: 422,
          detail: 'Expression syntax is invalid.',
          code: 'SYNTAX_ERROR',
        }),
        {
          status: 422,
          headers: { 'Content-Type': 'application/problem+json' },
        },
      ),
    )
    render(<App />)
    await user.click(screen.getByRole('radio', { name: 'Dark' }))
    await user.click(screen.getByRole('button', { name: 'Digit 7' }))
    await user.click(screen.getByRole('button', { name: 'Addition' }))
    await user.click(
      screen.getByRole('button', { name: 'Calculation history' }),
    )
    await user.click(screen.getByRole('button', { name: 'Equals' }))
    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Expression syntax is invalid.',
    )
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
