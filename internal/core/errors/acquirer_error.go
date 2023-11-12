package errors

type AcquirerError struct {
	Code int
	Err  error
}

func NewAcquirerError(code int, err error) *AcquirerError {
	return &AcquirerError{
		Code: code,
		Err:  err,
	}
}

func (e *AcquirerError) Error() string {
	return e.Err.Error()
}

func (e *AcquirerError) Unwrap() error {
	return e.Err
}
