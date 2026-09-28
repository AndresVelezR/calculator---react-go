package domain

import (
	"testing"
)

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
				t.Errorf(
					"Add(%v, %v) = %v, expected %v",
					testCase.firstNumber,
					testCase.secondNumber,
					actualResult,
					testCase.expectedResult,
				)
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
			actualResult := Subtract(
				testCase.firstNumber,
				testCase.secondNumber,
			)

			if actualResult != testCase.expectedResult {
				t.Errorf(
					"Subtract(%v, %v) = %v, expected %v",
					testCase.firstNumber,
					testCase.secondNumber,
					actualResult,
					testCase.expectedResult,
				)
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
			actualResult := Multiply(
				testCase.firstNumber,
				testCase.secondNumber,
			)

			if actualResult != testCase.expectedResult {
				t.Errorf(
					"Multiply(%v, %v) = %v, expected %v",
					testCase.firstNumber,
					testCase.secondNumber,
					actualResult,
					testCase.expectedResult,
				)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	testCases := []struct {
		name           string
		numerator      float64
		denominator    float64
		expectedResult float64
	}{
		{"exact division", 10, 2, 5},
		{"negative result", -9, 3, -3},
		{"decimal result", 1, 4, 0.25},
		{"zero divided by something", 0, 5, 0},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualResult, err := Divide(
				testCase.numerator,
				testCase.denominator,
			)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actualResult != testCase.expectedResult {
				t.Errorf(
					"Divide(%v, %v) = %v, expected %v",
					testCase.numerator,
					testCase.denominator,
					actualResult,
					testCase.expectedResult,
				)
			}
		})
	}
}