import { Display } from './Display'
import { Keypad } from './Keypad'
import { useCalculator } from './useCalculator'
import type { BinaryOperation } from './types'

const operationSymbols: Record<BinaryOperation, string> = {
  add: '+', subtract: '−', multiply: '×', divide: '÷', power: '^', percentage: '%',
}

export function Calculator() {
  const calculator = useCalculator()
  const expression = calculator.pendingOperation === null ? '' :
    `${calculator.previousNumber} ${operationSymbols[calculator.pendingOperation]}`

  return (
    <section className="calculator" aria-label="Calculator">
      <Display value={calculator.display} expression={expression} error={calculator.error} loading={calculator.loading} />
      <Keypad
        disabled={calculator.loading}
        onDigit={calculator.enterDigit}
        onDecimal={calculator.enterDecimal}
        onOperation={calculator.selectOperation}
        onEquals={calculator.equals}
        onClear={calculator.clear}
        onDelete={calculator.deleteDigit}
      />
      <p className="calculator-hint">Percentage: enter a value, press %, then enter the percent and =.</p>
      <p className="calculator-hint">Square root: enter a number and press √.</p>
    </section>
  )
}
