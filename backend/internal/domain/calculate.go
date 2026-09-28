package domain

import "math"

// Operation is the name of a math operation, like "add" or "divide".
type Operation string

// The operations the calculator knows about.
const (
	OperationAdd        Operation = "add"
	OperationSubtract   Operation = "subtract"
	OperationMultiply   Operation = "multiply"
	OperationDivide     Operation = "divide"
	OperationPower      Operation = "power"
	OperationSquareRoot Operation = "square_root"
	OperationPercentage Operation = "percentage"
)

// Calculate picks the right operation by name and runs it.
// This is the single entry point the HTTP layer will call.
// Note: square root only needs one number, so secondNumber is ignored there.
func Calculate(operation Operation, firstNumber, secondNumber float64) (float64, error) {
	var result float64
	var err error
	switch operation {
	case OperationAdd:
		result = Add(firstNumber, secondNumber)
	case OperationSubtract:
		result = Subtract(firstNumber, secondNumber)
	case OperationMultiply:
		result = Multiply(firstNumber, secondNumber)
	case OperationDivide:
		result, err = Divide(firstNumber, secondNumber)
	case OperationPower:
		result = Power(firstNumber, secondNumber)
	case OperationSquareRoot:
		result, err = SquareRoot(firstNumber)
	case OperationPercentage:
		result = Percentage(firstNumber, secondNumber)
	default:
		// nobody told us how to do this one
		return 0, ErrUnknownOperation
	}
	if err != nil {
		return 0, err
	}
	if math.IsInf(result, 0) || math.IsNaN(result) {
		return 0, ErrNonFiniteResult
	}
	return result, nil
}
