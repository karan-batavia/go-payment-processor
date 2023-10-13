package errors

type Acquirer struct {
	Code int
	Err  error
}

func NewAcquirer(code int, err error) *Acquirer {
	return &Acquirer{
		Code: code,
		Err:  err,
	}
}

func (e *Acquirer) Error() string {
	return e.Err.Error()
}

func (e *Acquirer) Unwrap() error {
	return e.Err
}
