import type { CalculateRequest, CalculateResponse, ErrorResponse } from './types'

const apiBaseUrl = (import.meta.env.VITE_API_URL ?? 'http://localhost:8080').replace(/\/+$/, '')

function isErrorResponse(value: unknown): value is ErrorResponse {
  return typeof value === 'object' && value !== null && 'error' in value &&
    typeof value.error === 'string' && value.error.trim().length > 0
}

function isCalculateResponse(value: unknown): value is CalculateResponse {
  return typeof value === 'object' && value !== null && 'result' in value &&
    typeof value.result === 'number' && Number.isFinite(value.result)
}

export async function calculate(request: CalculateRequest): Promise<number> {
  let response: Response
  try {
    response = await fetch(`${apiBaseUrl}/api/v1/calculate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request),
    })
  } catch {
    throw new Error('Unable to reach the calculator. Please try again.')
  }

  let body: unknown
  try {
    body = await response.json()
  } catch {
    throw new Error('The calculator returned an invalid response.')
  }

  if (!response.ok) {
    throw new Error(isErrorResponse(body) ? body.error : `Calculation failed (${response.status}).`)
  }
  if (!isCalculateResponse(body)) {
    throw new Error('The calculator returned an invalid result.')
  }
  return body.result
}
