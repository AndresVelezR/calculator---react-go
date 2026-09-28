import { act, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import App from '../App'

const fetchMock = vi.fn<typeof fetch>()

beforeEach(() => {
  fetchMock.mockReset()
  vi.stubGlobal('fetch', fetchMock)
})

async function pressKeys(keys: string[]) {
  const user = userEvent.setup()
  for (const name of keys) {
    await user.click(screen.getByRole('button', { name }))
  }
}

describe('calculator interactions', () => {
  it.each([
    { keys: ['1', 'Add', '2', 'Equals'], operation: 'add', a: 1, b: 2, result: 3 },
    { keys: ['1', 'Subtract', '2', 'Equals'], operation: 'subtract', a: 1, b: 2, result: -1 },
    { keys: ['3', 'Multiply', '4', 'Equals'], operation: 'multiply', a: 3, b: 4, result: 12 },
    { keys: ['1', 'Divide', '4', 'Equals'], operation: 'divide', a: 1, b: 4, result: 0.25 },
    { keys: ['2', 'Power', '3', 'Equals'], operation: 'power', a: 2, b: 3, result: 8 },
    { keys: ['5', '0', 'Percentage', '2', '0', 'Equals'], operation: 'percentage', a: 50, b: 20, result: 10 },
  ])('calculates $operation through the HTTP client', async ({ keys, operation, a, b, result }) => {
    fetchMock.mockResolvedValue(new Response(JSON.stringify({ result })))
    render(<App />)
    await pressKeys(keys)
    await waitFor(() => expect(screen.getByLabelText('Calculator display')).toHaveTextContent(String(result)))
    expect(JSON.parse(String(fetchMock.mock.calls[0][1]?.body))).toEqual({ operation, a, b })
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
  })

  it('calculates a square root with one number', async () => {
    fetchMock.mockResolvedValue(new Response(JSON.stringify({ result: 3 })))
    render(<App />)
    await pressKeys(['9', 'Square root'])
    await waitFor(() => expect(screen.getByLabelText('Calculator display')).toHaveTextContent('3'))
    expect(JSON.parse(String(fetchMock.mock.calls[0][1]?.body))).toEqual({ operation: 'square_root', a: 9 })
  })

  it('rounds a result but preserves decimal input and its original precision', async () => {
    fetchMock.mockResolvedValue(new Response(JSON.stringify({ result: 0.1 + 0.2 })))
    render(<App />)
    await pressKeys(['0', 'Decimal point', '1', 'Add', '0', 'Decimal point', '2', 'Equals'])
    await waitFor(() => expect(screen.getByLabelText('Calculator display')).toHaveTextContent(/^0\.3$/))
    await pressKeys(['Multiply'])
    expect(screen.getByLabelText('Pending calculation')).toHaveTextContent('0.3 ×')
    await pressKeys(['1', '0', 'Equals'])
    expect(JSON.parse(String(fetchMock.mock.calls[1][1]?.body))).toEqual({ operation: 'multiply', a: 0.1 + 0.2, b: 10 })
    await pressKeys(['Clear', 'Decimal point', '2', '0'])
    expect(screen.getByLabelText('Calculator display')).toHaveTextContent(/^0\.20$/)
  })

  it('shows a backend error and recovers after clearing', async () => {
    fetchMock.mockResolvedValueOnce(new Response(JSON.stringify({ error: 'division by zero' }), { status: 422 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ result: 3 })))
    render(<App />)
    await pressKeys(['1', 'Divide', '0', 'Equals'])
    expect(await screen.findByRole('alert')).toHaveTextContent('division by zero')
    await pressKeys(['Clear'])
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
    expect(screen.getByLabelText('Calculator display')).toHaveTextContent(/^0$/)
    await pressKeys(['1', 'Add', '2', 'Equals'])
    await waitFor(() => expect(screen.getByLabelText('Calculator display')).toHaveTextContent(/^3$/))
  })

  it('shows a clear network error', async () => {
    fetchMock.mockRejectedValue(new TypeError('Failed to fetch'))
    render(<App />)
    await pressKeys(['9', 'Square root'])
    expect(await screen.findByRole('alert')).toHaveTextContent('Unable to reach the calculator.')
  })

  it('validates missing input and supports decimal and delete keys', async () => {
    render(<App />)
    await pressKeys(['Equals'])
    expect(screen.getByRole('alert')).toHaveTextContent('Choose an operation first.')
    await pressKeys(['Clear', 'Decimal point', 'Decimal point', '5'])
    expect(screen.getByLabelText('Calculator display')).toHaveTextContent(/^0\.5$/)
    await pressKeys(['Delete last digit'])
    expect(screen.getByLabelText('Calculator display')).toHaveTextContent(/^0\.$/)
    await pressKeys(['Add', 'Equals'])
    expect(screen.getByRole('alert')).toHaveTextContent('Enter the second number.')
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('disables all keys while awaiting the response', async () => {
    let resolveResponse!: (response: Response) => void
    fetchMock.mockReturnValue(new Promise<Response>(resolve => { resolveResponse = resolve }))
    render(<App />)
    await pressKeys(['1', 'Add', '2', 'Equals'])
    expect(screen.getByText('Calculating…')).toBeInTheDocument()
    for (const button of screen.getAllByRole('button')) expect(button).toBeDisabled()
    await pressKeys(['Equals'])
    expect(fetchMock).toHaveBeenCalledTimes(1)
    await act(async () => { resolveResponse(new Response(JSON.stringify({ result: 3 }))) })
    await waitFor(() => expect(screen.getByRole('button', { name: 'Equals' })).toBeEnabled())
    expect(screen.getByLabelText('Calculator display')).toHaveTextContent(/^3$/)
  })
})
