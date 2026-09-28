package domain

// returns firstNumber + secondNumber.
func Add(firstNumber, secondNumber float64) float64 {
	return firstNumber + secondNumber
}

// returns firstNumber - secondNumber.
func Subtract(firstNumber, secondNumber float64) float64 {
	return firstNumber - secondNumber
}

// returns firstNumber * secondNumber.
func Multiply(firstNumber, secondNumber float64) float64 {
	return firstNumber * secondNumber
}

// Divide returns numerator / denominator.
// This one can fail, so it gives back two things: the result and an error.
// If denominator is zero we return ErrDivisionByZero
func Divide(numerator, denominator float64) (float64, error) {
	// can't divide by zero, stop right here
	if denominator == 0 {
		return 0, ErrDivisionByZero
	}

	// no error to report so we return nil
	return numerator / denominator, nil
}