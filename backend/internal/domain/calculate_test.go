package domain

import (
	"errors"
	"math"
	"testing"
)

func TestCalculate(t *testing.T) {
	testCases := []struct {
		name           string
		operation      Operation
		firstNumber    float64
		secondNumber   float64
		expectedResult float64
		expectedError  error
	}{
		{"add", OperationAdd, 2, 3, 5, nil},
		{"subtract", OperationSubtract, 2, 3, -1, nil},
		{"multiply", OperationMultiply, 2, 3, 6, nil},
		{"divide", OperationDivide, 6, 3, 2, nil},
		{"power", OperationPower, 2, 3, 8, nil},
		{"square root ignores second number", OperationSquareRoot, 9, -1, 3, nil},
		{"percentage", OperationPercentage, 50, 20, 10, nil},
		{"division by zero", OperationDivide, 1, 0, 0, ErrDivisionByZero},
		{"negative square root", OperationSquareRoot, -1, 0, 0, ErrNegativeSquareRoot},
		{"unknown operation", Operation("modulo"), 1, 2, 0, ErrUnknownOperation},
		{"empty operation", Operation(""), 1, 2, 0, ErrUnknownOperation},
		{"power overflow", OperationPower, 10, 1000, 0, ErrNonFiniteResult},
		{"power has no real result", OperationPower, -1, 0.5, 0, ErrNonFiniteResult},
		{"addition overflow", OperationAdd, math.MaxFloat64, math.MaxFloat64, 0, ErrNonFiniteResult},
		{"subtraction overflow", OperationSubtract, -math.MaxFloat64, math.MaxFloat64, 0, ErrNonFiniteResult},
		{"multiplication overflow", OperationMultiply, math.MaxFloat64, 2, 0, ErrNonFiniteResult},
		{"division overflow", OperationDivide, math.MaxFloat64, 0.5, 0, ErrNonFiniteResult},
		{"percentage overflow", OperationPercentage, math.MaxFloat64, 200, 0, ErrNonFiniteResult},
		{"NaN result", OperationSquareRoot, math.NaN(), 0, 0, ErrNonFiniteResult},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualResult, actualError := Calculate(testCase.operation, testCase.firstNumber, testCase.secondNumber)
			if !errors.Is(actualError, testCase.expectedError) {
				t.Fatalf("Calculate(%q, %v, %v) error = %v, expected %v",
					testCase.operation, testCase.firstNumber, testCase.secondNumber, actualError, testCase.expectedError)
			}
			if actualResult != testCase.expectedResult {
				t.Errorf("Calculate(%q, %v, %v) = %v, expected %v",
					testCase.operation, testCase.firstNumber, testCase.secondNumber, actualResult, testCase.expectedResult)
			}
		})
	}
}
