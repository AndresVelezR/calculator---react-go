import type { Operation } from './types'

interface KeypadProps {
  disabled: boolean
  onDigit: (digit: string) => void
  onDecimal: () => void
  onOperation: (operation: Operation) => void
  onEquals: () => void
  onClear: () => void
  onDelete: () => void
}

export function Keypad({ disabled, onDigit, onDecimal, onOperation, onEquals, onClear, onDelete }: KeypadProps) {
  return (
    <fieldset className="calculator-keypad" disabled={disabled}>
      <legend className="visually-hidden">Calculator keys</legend>
      <button type="button" className="key-utility" onClick={onClear} aria-label="Clear">C</button>
      <button type="button" className="key-utility" onClick={onDelete} aria-label="Delete last digit">⌫</button>
      <button type="button" className="key-operation" onClick={() => onOperation('square_root')} aria-label="Square root">√</button>
      <button type="button" className="key-operation" onClick={() => onOperation('divide')} aria-label="Divide">÷</button>
      {['7', '8', '9'].map(digit => <button type="button" key={digit} onClick={() => onDigit(digit)}>{digit}</button>)}
      <button type="button" className="key-operation" onClick={() => onOperation('multiply')} aria-label="Multiply">×</button>
      {['4', '5', '6'].map(digit => <button type="button" key={digit} onClick={() => onDigit(digit)}>{digit}</button>)}
      <button type="button" className="key-operation" onClick={() => onOperation('subtract')} aria-label="Subtract">−</button>
      {['1', '2', '3'].map(digit => <button type="button" key={digit} onClick={() => onDigit(digit)}>{digit}</button>)}
      <button type="button" className="key-operation" onClick={() => onOperation('add')} aria-label="Add">+</button>
      <button type="button" onClick={() => onDigit('0')}>0</button>
      <button type="button" onClick={onDecimal} aria-label="Decimal point">.</button>
      <button type="button" className="key-operation" onClick={() => onOperation('power')} aria-label="Power">xʸ</button>
      <button type="button" className="key-operation" onClick={() => onOperation('percentage')} aria-label="Percentage">%</button>
      <button type="button" className="key-equals" onClick={onEquals} aria-label="Equals">=</button>
    </fieldset>
  )
}
