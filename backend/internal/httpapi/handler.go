package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/AndresVelezR/calculator---react-go/backend/internal/domain"
)

func CalculateHandler(writer http.ResponseWriter, request *http.Request) {
	request.Body = http.MaxBytesReader(writer, request.Body, 4096)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	var input CalculateRequest
	if err := decoder.Decode(&input); err != nil {
		writeJSON(writer, http.StatusBadRequest, ErrorResponse{Error: "request must be a valid JSON object with operation, a, and b"})
		return
	}
	// Reject a second JSON value or trailing garbage.
	if err := decoder.Decode(new(any)); err != io.EOF {
		writeJSON(writer, http.StatusBadRequest, ErrorResponse{Error: "request must contain a single JSON object"})
		return
	}
	if input.Operation == "" {
		writeJSON(writer, http.StatusBadRequest, ErrorResponse{Error: "operation is required"})
		return
	}
	if input.A == nil {
		writeJSON(writer, http.StatusBadRequest, ErrorResponse{Error: "a is required"})
		return
	}
	if input.B == nil && input.Operation != domain.OperationSquareRoot {
		writeJSON(writer, http.StatusBadRequest, ErrorResponse{Error: "b is required for this operation"})
		return
	}

	var secondNumber float64
	if input.B != nil {
		secondNumber = *input.B
	}
	result, err := domain.Calculate(input.Operation, *input.A, secondNumber)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, domain.ErrUnknownOperation) {
			statusCode = http.StatusBadRequest
		} else if errors.Is(err, domain.ErrDivisionByZero) || errors.Is(err, domain.ErrNegativeSquareRoot) || errors.Is(err, domain.ErrNonFiniteResult) {
			statusCode = http.StatusUnprocessableEntity
		}
		writeJSON(writer, statusCode, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(writer, http.StatusOK, CalculateResponse{Result: result})
}

func writeJSON(writer http.ResponseWriter, statusCode int, response any) {
	// Encode before sending headers so encoding failures can return a proper status.
	body, err := json.Marshal(response)
	if err != nil {
		log.Printf("encode response: %v", err)
		statusCode = http.StatusInternalServerError
		body = []byte(`{"error":"could not encode response"}`)
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	if _, err := writer.Write(body); err != nil {
		log.Printf("write response: %v", err)
	}
}
