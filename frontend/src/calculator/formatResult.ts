// Keep the original number for later calculations; this is only for display.
export function formatResult(number: number): string {
  if (number !== 0 && Math.abs(number) < 1e-10) {
    return Number(number.toPrecision(10)).toString()
  }
  return Number(number.toFixed(10)).toString()
}
