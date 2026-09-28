package httpapi

import "github.com/AndresVelezR/calculator---react-go/backend/internal/domain"

type CalculateRequest struct {
	Operation domain.Operation `json:"operation"`
	// Pointers distinguish a missing number from an explicit zero.
	A *float64 `json:"a"`
	B *float64 `json:"b"`
}

type CalculateResponse struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
