import {
  act,
  cleanup,
  render,
  screen,
  waitFor,
  within,
} from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import App from './App'

const fetchMock = vi.fn<typeof fetch>()
function success(result: number) {
  return new Response(JSON.stringify({ result }), {
    headers: { 'Content-Type': 'application/json' },
  })
}
beforeEach(() => {
  localStorage.clear()
  fetchMock.mockReset()
  vi.stubGlobal('fetch', fetchMock)
})
afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

// Matrix: public UI -> native fetch boundary, every operation and keyboard,
// loading/duplicates/clear cancellation, errors/retry, history actions/reload,
// theme reload, invalid/unavailable persistence and clipboard feedback.
describe('functional calculator integration', () => {
  it.each([
    { keys: '2+3*4', ascii: '2+3*4', result: 14 },
    { keys: '(2+3)*4', ascii: '(2+3)*4', result: 20 },
    { keys: '-5+3', ascii: '-5+3', result: -2 },
    { keys: '10/4', ascii: '10/4', result: 2.5 },
    { keys: '2^3^2', ascii: '2^3^2', result: 512 },
    { keys: 'r9+7)', ascii: 'sqrt(9+7)', result: 4 },
    { keys: '200*15%', ascii: '200*15%', result: 30 },
    { keys: '.5+.25', ascii: '.5+.25', result: 0.75 },
  ])(
    'submits keyboard $keys through the expression API',
    async ({ keys, ascii, result }) => {
      const user = userEvent.setup()
      fetchMock.mockResolvedValueOnce(success(result))
      render(<App />)
      await user.keyboard(keys + '{Enter}')
      await waitFor(() =>
        expect(screen.getByTestId('calculator-result').textContent).toBe(
          String(result),
        ),
      )
      expect(fetchMock).toHaveBeenCalledExactlyOnceWith(
        '/api/v1/expressions',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ expression: ascii }),
          signal: expect.any(AbortSignal),
        }),
      )
      await user.click(
        screen.getByRole('button', { name: 'Calculation history' }),
      )
      expect(screen.getAllByRole('listitem')).toHaveLength(1)
      expect(screen.getByTestId('history-timestamp')).toHaveAttribute(
        'datetime',
      )
    },
  )
  it.each([
    'Clear',
    'Backspace',
    'Digit 4',
    'Multiply',
    'Multiplication',
    'Open parenthesis',
    'Decimal point',
    'Equals',
  ])('evaluates once with Enter after interacting with %s', async (name) => {
    const user = userEvent.setup()
    fetchMock.mockResolvedValueOnce(success(20))
    render(<App />)
    const button = screen.getByRole('button', { name })
    await user.click(button)
    // Escape edits the expression without moving focus off the clicked button.
    await user.keyboard('{Escape}(2+3)*4')
    expect(button).toHaveFocus()
    expect(screen.getByTestId('calculator-expression')).toHaveTextContent(
      '(2+3)×4',
    )
    await user.keyboard('{Enter}')
    await waitFor(() =>
      expect(screen.getByTestId('calculator-result').textContent).toBe('20'),
    )
    expect(screen.getByTestId('calculator-expression')).toHaveTextContent(
      '(2+3)×4',
    )
    expect(fetchMock).toHaveBeenCalledExactlyOnceWith(
      '/api/v1/expressions',
      expect.objectContaining({
        body: JSON.stringify({ expression: '(2+3)*4' }),
      }),
    )
    expect(
      JSON.parse(localStorage.getItem('sezzle.history.v1') ?? '{}').entries,
    ).toHaveLength(1)
  })

  it('keeps Clear usable while loading, blocks duplicate equals and ignores an aborted response', async () => {
    const user = userEvent.setup()
    let resolve!: (value: Response) => void
    fetchMock.mockReturnValueOnce(
      new Promise<Response>((done) => {
        resolve = done
      }),
    )
    render(<App />)
    await user.keyboard('1===')
    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect(
      screen.getByRole('status', { name: 'Calculation status' }),
    ).toHaveTextContent('Calculating')
    expect(screen.getByRole('button', { name: 'Equals' })).toBeDisabled()
    await user.click(screen.getByRole('button', { name: 'Clear' }))
    expect(fetchMock.mock.calls[0]?.[1]?.signal?.aborted).toBe(true)
    await act(async () => resolve(success(1)))
    expect(screen.getByTestId('calculator-result')).toHaveTextContent('—')
    expect(localStorage.getItem('sezzle.history.v1')).toBeNull()
  })
  it('recovers from network failure and stores only the successful retry', async () => {
    const user = userEvent.setup()
    fetchMock
      .mockRejectedValueOnce(new TypeError('offline'))
      .mockResolvedValueOnce(success(2))
    render(<App />)
    await user.keyboard('1+1{Enter}')
    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Could not reach',
    )
    expect(localStorage.getItem('sezzle.history.v1')).toBeNull()
    await user.keyboard('{Enter}')
    await waitFor(() =>
      expect(screen.getByTestId('calculator-result')).toHaveTextContent('2'),
    )
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
    expect(
      JSON.parse(localStorage.getItem('sezzle.history.v1') ?? '{}').entries,
    ).toHaveLength(1)
  })
  it('persists successes and theme, reuses, copies, deletes and confirms clearing', async () => {
    const user = userEvent.setup()
    const copy = vi
      .spyOn(navigator.clipboard, 'writeText')
      .mockResolvedValue(undefined)
    fetchMock
      .mockResolvedValueOnce(success(5))
      .mockResolvedValueOnce(success(4))
    const first = render(<App />)
    await user.keyboard('2+3{Enter}')
    await waitFor(() =>
      expect(screen.getByTestId('calculator-result')).toHaveTextContent('5'),
    )
    await user.keyboard('{Escape}2+2{Enter}')
    await waitFor(() =>
      expect(screen.getByTestId('calculator-result')).toHaveTextContent('4'),
    )
    await user.click(screen.getByRole('radio', { name: 'Dark' }))
    first.unmount()
    render(<App />)
    expect(screen.getByRole('radio', { name: 'Dark' })).toBeChecked()
    await user.click(
      screen.getByRole('button', { name: 'Calculation history' }),
    )
    const items = screen.getAllByRole('listitem')
    expect(items).toHaveLength(2)
    await user.click(
      within(items[0]).getByRole('button', { name: 'Reuse expression' }),
    )
    expect(screen.getByTestId('calculator-expression')).toHaveTextContent('2+2')
    expect(fetchMock).toHaveBeenCalledTimes(2)
    await user.click(
      within(items[0]).getByRole('button', { name: 'Copy result' }),
    )
    expect(copy).toHaveBeenCalledExactlyOnceWith('4')
    expect(
      await screen.findByRole('status', { name: 'Clipboard status' }),
    ).toHaveTextContent('Result copied.')
    await user.click(
      within(items[0]).getByRole('button', { name: 'Delete entry' }),
    )
    expect(screen.getAllByRole('listitem')).toHaveLength(1)
    await user.click(screen.getByRole('button', { name: 'Clear history' }))
    await user.click(screen.getByRole('button', { name: 'Cancel' }))
    expect(screen.getAllByRole('listitem')).toHaveLength(1)
    await user.click(screen.getByRole('button', { name: 'Clear history' }))
    await user.click(screen.getByRole('button', { name: 'Clear all entries' }))
    expect(
      screen.getByRole('heading', { name: 'No calculations yet' }),
    ).toBeInTheDocument()
    expect(
      JSON.parse(localStorage.getItem('sezzle.history.v1') ?? '{}').entries,
    ).toEqual([])
    expect(localStorage.getItem('sezzle.theme.v1')).toContain('dark')
  })
  it('keeps calculations usable with invalid saved data', async () => {
    localStorage.setItem('sezzle.history.v1', 'not-json')
    localStorage.setItem('sezzle.theme.v1', '{"version":99}')
    const user = userEvent.setup()
    fetchMock.mockResolvedValueOnce(success(3))
    render(<App />)
    expect(screen.getByRole('radio', { name: 'System' })).toBeChecked()
    expect(
      screen.getByRole('status', { name: 'Storage status' }),
    ).toBeInTheDocument()
    await user.keyboard('1+2{Enter}')
    await waitFor(() =>
      expect(screen.getByTestId('calculator-result')).toHaveTextContent('3'),
    )
  })
  it('keeps the result and history in memory when localStorage rejects writes', async () => {
    vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {
      throw new DOMException('Quota', 'QuotaExceededError')
    })
    const user = userEvent.setup()
    fetchMock.mockResolvedValueOnce(success(2))
    render(<App />)
    await user.keyboard('1+1{Enter}')
    await waitFor(() =>
      expect(screen.getByTestId('calculator-result')).toHaveTextContent('2'),
    )
    expect(
      screen.getByRole('status', { name: 'Storage status' }),
    ).toHaveTextContent('could not be saved')
    await user.click(
      screen.getByRole('button', { name: 'Calculation history' }),
    )
    expect(screen.getAllByRole('listitem')).toHaveLength(1)
  })
})
