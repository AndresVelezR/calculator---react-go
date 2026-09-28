import { beforeEach, describe, expect, it, vi } from 'vitest'

const fetchMock = vi.fn<typeof fetch>()

beforeEach(() => {
  vi.resetModules()
  fetchMock.mockReset()
  vi.stubGlobal('fetch', fetchMock)
  vi.stubEnv('VITE_API_URL', undefined)
})

describe('calculate API client', () => {
  it('sends JSON to the default endpoint and returns the result', async () => {
    fetchMock.mockResolvedValue(new Response(JSON.stringify({ result: 2 })))
    const { calculate } = await import('./api')
    await expect(calculate({ operation: 'add', a: 0, b: 2 })).resolves.toBe(2)
    expect(fetchMock).toHaveBeenCalledWith('http://localhost:8080/api/v1/calculate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: '{"operation":"add","a":0,"b":2}',
    })
  })

  it('omits b for a square root', async () => {
    fetchMock.mockResolvedValue(new Response(JSON.stringify({ result: 3 })))
    const { calculate } = await import('./api')
    await calculate({ operation: 'square_root', a: 9 })
    expect(fetchMock.mock.calls[0][1]?.body).toBe('{"operation":"square_root","a":9}')
  })

  it.each([
    ['https://calculator.example/', 'https://calculator.example/api/v1/calculate'],
    ['', '/api/v1/calculate'],
  ])('uses the configured base URL %s', async (baseUrl, expectedUrl) => {
    vi.stubEnv('VITE_API_URL', baseUrl)
    fetchMock.mockResolvedValue(new Response(JSON.stringify({ result: 3 })))
    const { calculate } = await import('./api')
    await calculate({ operation: 'add', a: 1, b: 2 })
    expect(fetchMock.mock.calls[0][0]).toBe(expectedUrl)
  })

  it.each([400, 422, 500])('preserves the backend error for status %s', async status => {
    fetchMock.mockResolvedValue(new Response(JSON.stringify({ error: 'division by zero' }), { status }))
    const { calculate } = await import('./api')
    await expect(calculate({ operation: 'divide', a: 1, b: 0 })).rejects.toThrow('division by zero')
  })

  it.each([null, {}, { error: '' }, { error: ' ' }, { error: 123 }])('provides a fallback for an error body %j', async body => {
    fetchMock.mockResolvedValue(new Response(JSON.stringify(body), { status: 503 }))
    const { calculate } = await import('./api')
    await expect(calculate({ operation: 'add', a: 1, b: 2 })).rejects.toThrow('Calculation failed (503).')
  })

  it('reports a network failure clearly', async () => {
    fetchMock.mockRejectedValue(new TypeError('Failed to fetch'))
    const { calculate } = await import('./api')
    await expect(calculate({ operation: 'add', a: 1, b: 2 })).rejects.toThrow('Unable to reach the calculator.')
  })

  it('rejects malformed response JSON', async () => {
    fetchMock.mockResolvedValue(new Response('not JSON'))
    const { calculate } = await import('./api')
    await expect(calculate({ operation: 'add', a: 1, b: 2 })).rejects.toThrow('The calculator returned an invalid response.')
  })

  it.each([null, [], {}, { result: '3' }, { result: null }, { result: true }])('rejects an invalid result %j', async body => {
    fetchMock.mockResolvedValue(new Response(JSON.stringify(body)))
    const { calculate } = await import('./api')
    await expect(calculate({ operation: 'add', a: 1, b: 2 })).rejects.toThrow('The calculator returned an invalid result.')
  })
})
