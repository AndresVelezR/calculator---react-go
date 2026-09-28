interface DisplayProps {
  value: string
  expression: string
  error: string | null
  loading: boolean
}

export function Display({ value, expression, error, loading }: DisplayProps) {
  return (
    <div className="calculator-display" aria-busy={loading}>
      <p className="calculator-expression" aria-label="Pending calculation">{expression || 'Ready when you are'}</p>
      <output className="calculator-value" aria-label="Calculator display" aria-live="polite">{value}</output>
      <p className="calculator-status" role="status">{loading ? 'Calculating…' : ' '}</p>
      {error && <p className="calculator-error" role="alert">{error}</p>}
    </div>
  )
}
