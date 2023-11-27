package errors

import "bytes"

type ValidationError struct {
	Messages []string
}

func NewValidationError(messages ...string) *ValidationError {
	return &ValidationError{
		Messages: messages,
	}
}

func (e *ValidationError) Error() string {
	buff := bytes.NewBufferString("")
	for _, msg := range e.Messages {
		buff.WriteString(msg)
		buff.WriteString("\n")
	}
	return buff.String()
}
