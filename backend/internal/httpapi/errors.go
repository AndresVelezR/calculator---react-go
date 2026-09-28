package httpapi

import (
	"errors"
	"net/http"

	"github.com/AndresVelezR/calculator---react-go/backend/internal/domain"
)

func domainErrorStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrUnknownOperation):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrDivisionByZero),
		errors.Is(err, domain.ErrNegativeSquareRoot),
		errors.Is(err, domain.ErrNonFiniteResult):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
