package errors

import "errors"

var (
	ErrNotFound        = errors.New("resource not found")
	ErrConflict        = errors.New("resource conflict")
	ErrInvalidInput    = errors.New("invalid input")
	ErrAlreadyAnswered = errors.New("authorization already answered")
	ErrInvalidStatus   = errors.New("invalid authorization status")
)
