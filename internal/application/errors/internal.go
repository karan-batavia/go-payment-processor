package errors

type Internal struct {
	err error
}

func NewInternal(err error) *Internal {
	return &Internal{
		err: err,
	}
}

func (e *Internal) Error() string {
	return e.err.Error()
}

func (e *Internal) Unwrap() error {
	return e.err
}
