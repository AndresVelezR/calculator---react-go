import { describe, expect, it } from 'vitest'
import { formatResult } from './formatResult'

describe('formatResult', () => {
  it.each([
    [0.1 + 0.2, '0.3'],
    [1 / 3, '0.3333333333'],
    [10, '10'],
    [1.25, '1.25'],
    [-0.1 - 0.2, '-0.3'],
    [-0, '0'],
    [1e-12, '1e-12'],
    [-1e-12, '-1e-12'],
    [1e21, '1e+21'],
  ])('formats %s as %s', (number, expected) => {
    expect(formatResult(number)).toBe(expected)
  })
})
