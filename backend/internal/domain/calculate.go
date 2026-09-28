package domain

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
	switch operation {
	case OperationAdd:
		return Add(firstNumber, secondNumber), nil
	case OperationSubtract:
		return Subtract(firstNumber, secondNumber), nil
	case OperationMultiply:
		return Multiply(firstNumber, secondNumber), nil
	case OperationDivide:
		return Divide(firstNumber, secondNumber)
	case OperationPower:
		return Power(firstNumber, secondNumber), nil
	case OperationSquareRoot:
		return SquareRoot(firstNumber)
	case OperationPercentage:
		return Percentage(firstNumber, secondNumber), nil
	default:
		// nobody told us how to do this one
		return 0, ErrUnknownOperation
	}
}