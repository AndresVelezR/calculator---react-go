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