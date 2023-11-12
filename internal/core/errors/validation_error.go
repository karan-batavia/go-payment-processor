package errors

import "errors"

type ValidationError struct {
	err error
}

func NewValidationError(errs ...error) *ValidationError {
	return &ValidationError{
		err: errors.Join(errs...),
	}
}

func (e *ValidationError) Error() string {
	return e.Error()
}

func (e *ValidationError) Unwrap() []error {
	return e.err.(interface{ Unwrap() []error }).Unwrap()
}
