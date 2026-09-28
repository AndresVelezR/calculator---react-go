import { useRef, useState } from 'react'
import { calculate } from './api'
import type { BinaryOperation, CalculateRequest, Operation } from './types'

export function useCalculator() {
  const [display, setDisplay] = useState('0')
  const [previousNumber, setPreviousNumber] = useState<number | null>(null)
  const [pendingOperation, setPendingOperation] = useState<BinaryOperation | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [isNewEntry, setIsNewEntry] = useState(true)
  // A ref also blocks a second click before React renders the loading state.
  const requestInFlight = useRef(false)

  function enterDigit(digit: string) {
    if (requestInFlight.current || !/^\d$/.test(digit)) return
    setError(null)
    setDisplay(currentDisplay => isNewEntry || currentDisplay === '0' ? digit : currentDisplay + digit)
    setIsNewEntry(false)
  }

  function enterDecimal() {
    if (requestInFlight.current) return
    setError(null)
    if (isNewEntry) {
      setDisplay('0.')
      setIsNewEntry(false)
    } else if (!display.includes('.')) {
      setDisplay(display + '.')
    }
  }

  function clear() {
    if (requestInFlight.current) return
    setDisplay('0')
    setPreviousNumber(null)
    setPendingOperation(null)
    setError(null)
    setIsNewEntry(true)
  }

  function deleteDigit() {
    if (requestInFlight.current) return
    setError(null)
    if (isNewEntry && pendingOperation !== null) return
    const shortenedDisplay = isNewEntry ? '' : display.slice(0, -1)
    setDisplay(shortenedDisplay === '' || shortenedDisplay === '-' ? '0' : shortenedDisplay)
    setIsNewEntry(false)
  }

  function readNumber(): number | null {
    const number = Number(display)
    if (display.trim() === '' || !Number.isFinite(number)) {
      setError('Enter a finite number.')
      return null
    }
    return number
  }

  async function submit(request: CalculateRequest, nextOperation: BinaryOperation | null = null) {
    requestInFlight.current = true
    setLoading(true)
    setError(null)
    try {
      const result = await calculate(request)
      setDisplay(String(result))
      setPreviousNumber(nextOperation === null ? null : result)
      setPendingOperation(nextOperation)
      setIsNewEntry(true)
    } catch (failure) {
      setError(failure instanceof Error ? failure.message : 'Unable to calculate. Please try again.')
    } finally {
      requestInFlight.current = false
      setLoading(false)
    }
  }

  async function selectOperation(operation: Operation) {
    if (requestInFlight.current) return
    const number = readNumber()
    if (number === null) return
    if (operation === 'square_root') {
      if (pendingOperation !== null) {
        setError('Finish the current calculation before using square root.')
        return
      }
      await submit({ operation, a: number })
      return
    }
    if (pendingOperation !== null && previousNumber !== null && !isNewEntry) {
      await submit({ operation: pendingOperation, a: previousNumber, b: number }, operation)
      return
    }
    setPreviousNumber(number)
    setPendingOperation(operation)
    setIsNewEntry(true)
    setError(null)
  }

  async function equals() {
    if (requestInFlight.current) return
    if (pendingOperation === null || previousNumber === null) {
      setError('Choose an operation first.')
      return
    }
    if (isNewEntry) {
      setError('Enter the second number.')
      return
    }
    const number = readNumber()
    if (number === null) return
    await submit({ operation: pendingOperation, a: previousNumber, b: number })
  }

  return {
    display, previousNumber, pendingOperation, error, loading, isNewEntry,
    enterDigit, enterDecimal, selectOperation, equals, clear, deleteDigit,
  }
}
