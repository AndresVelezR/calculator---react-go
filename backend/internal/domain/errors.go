package domain

import "errors"

/*
These are the ways an operation can go wrong.
Other packages check for them with errors.Is(err, domain.ErrSomething)
instead of comparing message strings.
*/
var (
	// Dividing by zero
	ErrDivisionByZero = errors.New("division by zero")

	// Square root of something below zero 
	ErrNegativeSquareRoot = errors.New("square root of a negative number")

	// unknown operation
	ErrUnknownOperation = errors.New("unknown operation")
)