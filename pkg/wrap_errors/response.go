package wraperrors

import "errors"

const (
	NotFound         = "NOT_FOUND"
	AlreadyExists    = "ALREADY_EXISTS"
	ValidationFailed = "VALIDATION_FAILED"
	Unauthorized     = "UNAUTHORIZED"
	Forbidden        = "FORBIDDEN"
	Internal         = "INTERNAL"
)

type AppError struct {
	Code    int
	Message string
	Type    string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Is(target error) bool {
	return errors.Is(e, target)
}

func New(message, errorType string, code int, err error) *AppError {
	return &AppError{
		Code:    code,
		Type:    errorType,
		Message: message,
		Err:     err,
	}
}

func NotFoundErr(msg string) *AppError {
	return New(msg, NotFound, 404, nil)
}

func AlreadyExistsErr(msg string) *AppError {
	return New(msg, AlreadyExists, 409, nil)
}

func ValidationErr(msg string) *AppError {
	return New(msg, ValidationFailed, 400, nil)
}

func UnauthorizedErr(msg string) *AppError {
	return New(msg, Unauthorized, 401, nil)
}

func ForbiddenErr(msg string) *AppError {
	return New(msg, Forbidden, 403, nil)
}

func InternalErr(msg string, err error) *AppError {
	return New(msg, Internal, 500, err)
}
