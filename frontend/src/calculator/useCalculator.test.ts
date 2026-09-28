import { act, renderHook } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { calculate } from './api'
import { useCalculator } from './useCalculator'
import type { BinaryOperation } from './types'

vi.mock('./api', () => ({ calculate: vi.fn() }))
const calculateMock = vi.mocked(calculate)

beforeEach(() => {
  calculateMock.mockReset()
})

function enterNumber(result: { current: ReturnType<typeof useCalculator> }, number: string) {
  for (const character of number) {
    act(() => character === '.' ? result.current.enterDecimal() : result.current.enterDigit(character))
  }
}

describe('useCalculator', () => {
  it('starts at zero, ignores invalid digits, and accepts only one decimal point', () => {
    const { result } = renderHook(useCalculator)
    expect(result.current.display).toBe('0')
    act(() => result.current.enterDigit('x'))
    enterNumber(result, '00.1.2')
    expect(result.current.display).toBe('0.12')
    expect(calculateMock).not.toHaveBeenCalled()
  })

  it('starts decimal input with a leading zero', () => {
    const { result } = renderHook(useCalculator)
    act(() => result.current.enterDecimal())
    expect(result.current.display).toBe('0.')
  })

  it('deletes digits without leaving an empty display and clears pending state', async () => {
    const { result } = renderHook(useCalculator)
    enterNumber(result, '12')
    act(() => result.current.deleteDigit())
    expect(result.current.display).toBe('1')
    act(() => result.current.deleteDigit())
    expect(result.current.display).toBe('0')
    await act(() => result.current.selectOperation('add'))
    await act(() => result.current.equals())
    expect(result.current.error).toBe('Enter the second number.')
    act(() => result.current.clear())
    expect(result.current).toMatchObject({ display: '0', previousNumber: null, pendingOperation: null, error: null, loading: false })
  })

  const operations: BinaryOperation[] = ['add', 'subtract', 'multiply', 'divide', 'power', 'percentage']
  it.each(operations)('sends %s with both operands to the API', async operation => {
    calculateMock.mockResolvedValue(42)
    const { result } = renderHook(useCalculator)
    enterNumber(result, '12')
    await act(() => result.current.selectOperation(operation))
    enterNumber(result, '3')
    await act(() => result.current.equals())
    expect(calculateMock).toHaveBeenCalledWith({ operation, a: 12, b: 3 })
    expect(result.current).toMatchObject({ display: '42', pendingOperation: null, previousNumber: null, loading: false, error: null })
    enterNumber(result, '7')
    expect(result.current.display).toBe('7')
  })

  it('sends square root without a second operand', async () => {
    calculateMock.mockResolvedValue(3)
    const { result } = renderHook(useCalculator)
    enterNumber(result, '9')
    await act(() => result.current.selectOperation('square_root'))
    expect(calculateMock).toHaveBeenCalledWith({ operation: 'square_root', a: 9 })
    expect(result.current.display).toBe('3')
  })

  it('requires an operation and a second operand before equals', async () => {
    const { result } = renderHook(useCalculator)
    await act(() => result.current.equals())
    expect(result.current.error).toBe('Choose an operation first.')
    await act(() => result.current.selectOperation('add'))
    await act(() => result.current.equals())
    expect(result.current.error).toBe('Enter the second number.')
    expect(calculateMock).not.toHaveBeenCalled()
  })

  it('replaces a pending operation without calculating and preserves it when deleting', async () => {
    const { result } = renderHook(useCalculator)
    enterNumber(result, '5')
    await act(() => result.current.selectOperation('add'))
    act(() => result.current.deleteDigit())
    await act(() => result.current.selectOperation('multiply'))
    expect(result.current).toMatchObject({ previousNumber: 5, pendingOperation: 'multiply', display: '5' })
    expect(calculateMock).not.toHaveBeenCalled()
  })

  it('evaluates chained operations left to right', async () => {
    calculateMock.mockResolvedValueOnce(5).mockResolvedValueOnce(20)
    const { result } = renderHook(useCalculator)
    enterNumber(result, '2')
    await act(() => result.current.selectOperation('add'))
    enterNumber(result, '3')
    await act(() => result.current.selectOperation('multiply'))
    expect(result.current).toMatchObject({ previousNumber: 5, pendingOperation: 'multiply', display: '5' })
    enterNumber(result, '4')
    await act(() => result.current.equals())
    expect(calculateMock).toHaveBeenNthCalledWith(1, { operation: 'add', a: 2, b: 3 })
    expect(calculateMock).toHaveBeenNthCalledWith(2, { operation: 'multiply', a: 5, b: 4 })
    expect(result.current.display).toBe('20')
  })

  it('asks to finish a binary calculation before square root', async () => {
    const { result } = renderHook(useCalculator)
    await act(() => result.current.selectOperation('add'))
    await act(() => result.current.selectOperation('square_root'))
    expect(result.current.error).toBe('Finish the current calculation before using square root.')
    expect(calculateMock).not.toHaveBeenCalled()
  })

  it('preserves full precision for the next calculation', async () => {
    calculateMock.mockResolvedValue(0.1 + 0.2)
    const { result } = renderHook(useCalculator)
    enterNumber(result, '0.1')
    await act(() => result.current.selectOperation('add'))
    enterNumber(result, '0.2')
    await act(() => result.current.equals())
    await act(() => result.current.selectOperation('multiply'))
    enterNumber(result, '10')
    await act(() => result.current.equals())
    expect(calculateMock).toHaveBeenLastCalledWith({ operation: 'multiply', a: 0.1 + 0.2, b: 10 })
  })

  it('shows backend errors and allows correcting an operand', async () => {
    calculateMock.mockRejectedValueOnce(new Error('division by zero')).mockResolvedValueOnce(2)
    const { result } = renderHook(useCalculator)
    enterNumber(result, '4')
    await act(() => result.current.selectOperation('divide'))
    enterNumber(result, '0')
    await act(() => result.current.equals())
    expect(result.current).toMatchObject({ error: 'division by zero', loading: false })
    enterNumber(result, '2')
    expect(result.current.error).toBeNull()
    await act(() => result.current.equals())
    expect(result.current.display).toBe('2')
  })

  it('handles an unexpected rejection without losing the loading state', async () => {
    calculateMock.mockRejectedValue('unexpected')
    const { result } = renderHook(useCalculator)
    await act(() => result.current.selectOperation('square_root'))
    expect(result.current).toMatchObject({ error: 'Unable to calculate. Please try again.', loading: false })
  })

  it('blocks duplicate requests and input while a calculation is pending', async () => {
    let resolveCalculation!: (value: number) => void
    calculateMock.mockReturnValue(new Promise<number>(resolve => { resolveCalculation = resolve }))
    const { result } = renderHook(useCalculator)
    enterNumber(result, '1')
    await act(() => result.current.selectOperation('add'))
    enterNumber(result, '2')
    let pendingCalculation!: Promise<void>
    act(() => { pendingCalculation = result.current.equals() })
    expect(result.current.loading).toBe(true)
    act(() => {
      result.current.enterDigit('9')
      result.current.enterDecimal()
      result.current.deleteDigit()
      result.current.clear()
    })
    await act(() => result.current.equals())
    await act(() => result.current.selectOperation('square_root'))
    expect(calculateMock).toHaveBeenCalledTimes(1)
    expect(result.current.display).toBe('2')
    await act(async () => { resolveCalculation(3); await pendingCalculation })
    expect(result.current).toMatchObject({ display: '3', loading: false })
  })

  it('rejects an overflowing input before requesting a calculation', async () => {
    const { result } = renderHook(useCalculator)
    enterNumber(result, '9'.repeat(310))
    await act(() => result.current.selectOperation('add'))
    expect(result.current.error).toBe('Enter a finite number.')
    expect(calculateMock).not.toHaveBeenCalled()
    act(() => result.current.clear())
    enterNumber(result, '1')
    await act(() => result.current.selectOperation('add'))
    enterNumber(result, '9'.repeat(310))
    await act(() => result.current.equals())
    expect(result.current.error).toBe('Enter a finite number.')
    expect(calculateMock).not.toHaveBeenCalled()
  })

  it('starts fresh when deleting or entering a decimal after a result', async () => {
    calculateMock.mockResolvedValue(3)
    const { result } = renderHook(useCalculator)
    await act(() => result.current.selectOperation('square_root'))
    act(() => result.current.deleteDigit())
    expect(result.current.display).toBe('0')
    await act(() => result.current.selectOperation('square_root'))
    act(() => result.current.enterDecimal())
    expect(result.current.display).toBe('0.')
  })
})
