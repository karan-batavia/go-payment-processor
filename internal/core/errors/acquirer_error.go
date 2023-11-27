package errors

type AcquirerError struct {
	Code    int
	Message string
}

func NewAcquirerError(code int, message string) *AcquirerError {
	return &AcquirerError{
		Code:    code,
		Message: message,
	}
}

func (e *AcquirerError) Error() string {
	return e.Message
}
