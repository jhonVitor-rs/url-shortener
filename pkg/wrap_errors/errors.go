package wraperrors

import (
	"errors"
	"strings"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrValidation      = errors.New("validation failed")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrInternal        = errors.New("internal error")
	ErrUniqueViolation = errors.New("unique constraint violation")
)

type AppError struct {
	Code    int
	Message string
	Type    error
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Is(target error) bool {
	targetErr, ok := target.(*AppError)
	if !ok {
		return errors.Is(e.Type, target)
	}

	return errors.Is(e.Type, targetErr.Type)
}

func (e *AppError) As(target any) bool {
	return errors.As(e, target)
}

func New(message string, errorType error, code int, err error) *AppError {
	return &AppError{
		Code:    code,
		Type:    errorType,
		Message: message,
		Err:     err,
	}
}

func NotFoundErr(msg string) *AppError {
	return New(msg, ErrNotFound, 404, nil)
}

func AlreadyExistsErr(msg string) *AppError {
	return New(msg, ErrAlreadyExists, 409, nil)
}

func ValidationErr(msg string) *AppError {
	return New(msg, ErrValidation, 400, nil)
}

func UnauthorizedErr(msg string) *AppError {
	return New(msg, ErrUnauthorized, 401, nil)
}

func InternalErr(msg string, err error) *AppError {
	return New(msg, ErrInternal, 500, err)
}

func UniqueViolationErr(message string) error {
	return &AppError{
		Type:    ErrUniqueViolation,
		Message: message,
		Err:     errors.New(message),
	}
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func IsAlreadyExistsError(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}

func IsValidationError(err error) bool {
	return errors.Is(err, ErrValidation)
}

func IsUnauthorizedError(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

func IsInternalError(err error) bool {
	return errors.Is(err, ErrInternal)
}

func IsUniqueViolation(err error) bool {
	if errors.Is(err, ErrUniqueViolation) {
		return true
	}

	return err != nil && strings.Contains(err.Error(), "unique constraint")
}
