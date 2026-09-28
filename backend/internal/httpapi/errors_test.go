package httpapi

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/AndresVelezR/calculator---react-go/backend/internal/domain"
)

func TestDomainErrorStatus(t *testing.T) {
	testCases := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{"unknown operation", domain.ErrUnknownOperation, http.StatusBadRequest},
		{"division by zero", domain.ErrDivisionByZero, http.StatusUnprocessableEntity},
		{"negative square root", domain.ErrNegativeSquareRoot, http.StatusUnprocessableEntity},
		{"non-finite result", domain.ErrNonFiniteResult, http.StatusUnprocessableEntity},
		{"wrapped domain error", fmt.Errorf("calculate: %w", domain.ErrDivisionByZero), http.StatusUnprocessableEntity},
		{"unexpected error", errors.New("unexpected failure"), http.StatusInternalServerError},
		{"same text is not a sentinel", errors.New("division by zero"), http.StatusInternalServerError},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actualStatus := domainErrorStatus(testCase.err); actualStatus != testCase.expectedStatus {
				t.Errorf("status = %d, expected %d", actualStatus, testCase.expectedStatus)
			}
		})
	}
}
