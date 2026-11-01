package errors

import "errors"

var (
	ErrNotFound        = errors.New("resource not found")
	ErrConflict        = errors.New("resource conflict")
	ErrInvalidInput    = errors.New("invalid input")
	ErrDuplicate       = errors.New("duplicate resource")
	ErrAlreadyAnswered = errors.New("authorization already answered")
	ErrInvalidStatus   = errors.New("invalid authorization status")
)
