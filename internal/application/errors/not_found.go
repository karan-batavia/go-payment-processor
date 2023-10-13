package errors

type NotFound struct {
	err error
}

func NewNotFound(err error) *NotFound {
	return &NotFound{
		err: err,
	}
}

func (e *NotFound) Error() string {
	return e.err.Error()
}

func (e *NotFound) Unwrap() error {
	return e.err
}
