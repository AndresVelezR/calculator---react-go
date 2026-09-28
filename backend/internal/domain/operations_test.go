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

// Divide can fail, so this table also has an
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

func TestPower(t *testing.T) {
	testCases := []struct {
		name           string
		base           float64
		exponent       float64
		expectedResult float64
	}{
		{"positive exponent", 2, 3, 8},
		{"zero exponent", 5, 0, 1},
		{"negative exponent", 2, -2, 0.25},
		{"fractional exponent", 9, 0.5, 3},
		{"negative base", -2, 3, -8},
		{"zero base", 0, 3, 0},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualResult := Power(testCase.base, testCase.exponent)
			if actualResult != testCase.expectedResult {
				t.Errorf("Power(%v, %v) = %v, expected %v",
					testCase.base, testCase.exponent, actualResult, testCase.expectedResult)
			}
		})
	}
}

func TestSquareRoot(t *testing.T) {
	testCases := []struct {
		name           string
		number         float64
		expectedResult float64
		expectedError  error
	}{
		{"positive number", 9, 3, nil},
		{"zero", 0, 0, nil},
		{"decimal", 0.25, 0.5, nil},
		{"negative number", -4, 0, ErrNegativeSquareRoot},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualResult, actualError := SquareRoot(testCase.number)
			if !errors.Is(actualError, testCase.expectedError) {
				t.Fatalf("SquareRoot(%v) error = %v, expected %v",
					testCase.number, actualError, testCase.expectedError)
			}
			if actualResult != testCase.expectedResult {
				t.Errorf("SquareRoot(%v) = %v, expected %v",
					testCase.number, actualResult, testCase.expectedResult)
			}
		})
	}
}

func TestPercentage(t *testing.T) {
	testCases := []struct {
		name           string
		value          float64
		percent        float64
		expectedResult float64
	}{
		{"twenty percent", 50, 20, 10},
		{"zero percent", 50, 0, 0},
		{"zero value", 0, 20, 0},
		{"over one hundred percent", 50, 150, 75},
		{"negative value", -50, 20, -10},
		{"negative percent", 50, -20, -10},
		{"decimal percent", 200, 0.5, 1},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualResult := Percentage(testCase.value, testCase.percent)
			if actualResult != testCase.expectedResult {
				t.Errorf("Percentage(%v, %v) = %v, expected %v",
					testCase.value, testCase.percent, actualResult, testCase.expectedResult)
			}
		})
	}
}
