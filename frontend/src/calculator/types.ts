export type BinaryOperation =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'power'
  | 'percentage'

export type Operation = BinaryOperation | 'square_root'

export type CalculateRequest =
  | { operation: BinaryOperation; a: number; b: number }
  | { operation: 'square_root'; a: number; b?: number }

export interface CalculateResponse {
  result: number
}

export interface ErrorResponse {
  error: string
}
