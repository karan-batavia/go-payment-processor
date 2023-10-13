package errors

import "errors"

type Validation struct {
	err error
}

func NewValidation(errs ...error) *Validation {
	return &Validation{
		err: errors.Join(errs...),
	}
}

func (e *Validation) Error() string {
	return e.Error()
}

func (e *Validation) Unwrap() []error {
	return e.err.(interface{ Unwrap() []error }).Unwrap()
}
