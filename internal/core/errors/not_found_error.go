package errors

type NotFoundError struct {
	err error
}

func NewNotFoundError(err error) *NotFoundError {
	return &NotFoundError{
		err: err,
	}
}

func (e *NotFoundError) Error() string {
	return e.err.Error()
}

func (e *NotFoundError) Unwrap() error {
	return e.err
}
