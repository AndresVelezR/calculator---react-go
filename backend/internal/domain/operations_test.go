package domain

import (
	"errors"
	"testing"
)

// Each test below uses a "table": a list of cases, and we run
// the same check on every one of them.

func TestAdd(t *testing.T) {
	testCases := []struct {
		name           string
		firstNumber    float64
		secondNumber   float64
		expectedResult float64
	}{
		{"two positive numbers", 2, 3, 5},
		{"positive and negative", 5, -8, -3},
		{"adding zero", 7, 0, 7},
		{"decimals", 1.5, 2.25, 3.75},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualResult := Add(testCase.firstNumber, testCase.secondNumber)
			if actualResult != testCase.expectedResult {
				t.Errorf("Add(%v, %v) = %v, expected %v",
					testCase.firstNumber, testCase.secondNumber, actualResult, testCase.expectedResult)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	testCases := []struct {
		name           string
		firstNumber    float64
		secondNumber   float64
		expectedResult float64
	}{
		{"simple subtraction", 10, 4, 6},
		{"result is negative", 4, 10, -6},
		{"subtracting zero", 5, 0, 5},
		{"decimals", 2.5, 0.5, 2},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualResult := Subtract(testCase.firstNumber, testCase.secondNumber)
			if actualResult != testCase.expectedResult {
				t.Errorf("Subtract(%v, %v) = %v, expected %v",
					testCase.firstNumber, testCase.secondNumber, actualResult, testCase.expectedResult)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	testCases := []struct {
		name           string
		firstNumber    float64
		secondNumber   float64
		expectedResult float64
	}{
		{"two positive numbers", 3, 4, 12},
		{"one negative", -3, 4, -12},
		{"two negatives", -3, -4, 12},
		{"times zero", 5, 0, 0},
		{"decimals", 0.5, 4, 2},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualResult := Multiply(testCase.firstNumber, testCase.secondNumber)
			if actualResult != testCase.expectedResult {
				t.Errorf("Multiply(%v, %v) = %v, expected %v",
					testCase.firstNumber, testCase.secondNumber, actualResult, testCase.expectedResult)
			}
		})
	}
}

// Divide is the only one that can fail, so this table also has an
// expectedError column. nil there means "we expect no error".
func TestDivide(t *testing.T) {
	testCases := []struct {
		name           string
		numerator      float64
		denominator    float64
		expectedResult float64
		expectedError  error
	}{
		{"exact division", 10, 2, 5, nil},
		{"negative result", -9, 3, -3, nil},
		{"decimal result", 1, 4, 0.25, nil},
		{"zero divided by something", 0, 5, 0, nil},
		{"divide by zero", 5, 0, 0, ErrDivisionByZero},
		{"zero divided by zero", 0, 0, 0, ErrDivisionByZero},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualResult, actualError := Divide(testCase.numerator, testCase.denominator)

			// errors.Is also works when both are nil, so one check covers both cases
			if !errors.Is(actualError, testCase.expectedError) {
				t.Fatalf("Divide(%v, %v) error = %v, expected %v",
					testCase.numerator, testCase.denominator, actualError, testCase.expectedError)
			}
			if actualResult != testCase.expectedResult {
				t.Errorf("Divide(%v, %v) = %v, expected %v",
					testCase.numerator, testCase.denominator, actualResult, testCase.expectedResult)
			}
		})
	}
}