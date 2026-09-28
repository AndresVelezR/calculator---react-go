package domain

import "math"

// Add returns firstNumber + secondNumber.
func Add(firstNumber, secondNumber float64) float64 {
	return firstNumber + secondNumber
}

// Subtract returns firstNumber - secondNumber.
func Subtract(firstNumber, secondNumber float64) float64 {
	return firstNumber - secondNumber
}

// Multiply returns firstNumber * secondNumber.
func Multiply(firstNumber, secondNumber float64) float64 {
	return firstNumber * secondNumber
}

// Divide returns numerator / denominator.
// This one can fail, so it gives back two things: the result and an error.
// If denominator is zero we return ErrDivisionByZero (the 0 is just a
// placeholder, nobody should use it when there's an error).
func Divide(numerator, denominator float64) (float64, error) {
	// can't divide by zero, stop right here
	if denominator == 0 {
		return 0, ErrDivisionByZero
	}

	// all good, no error to report so we return nil
	return numerator / denominator, nil
}

// Power returns base raised to the exponent, e.g. Power(2, 3) is 8.
func Power(base, exponent float64) float64 {
	return math.Pow(base, exponent)
}

// SquareRoot returns the square root of number.
// Negative numbers have no real square root, so that's an error.
func SquareRoot(number float64) (float64, error) {
	if number < 0 {
		return 0, ErrNegativeSquareRoot
	}
	return math.Sqrt(number), nil
}

// Percentage answers "what is percent% of value?"
// Example: Percentage(50, 20) is 10, because 20% of 50 is 10.
func Percentage(value, percent float64) float64 {
	return value * percent / 100
}
