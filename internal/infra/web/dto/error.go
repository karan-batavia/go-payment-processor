package dto

import "bytes"

type Error struct {
	Messages []string
}

func NewError(messages ...string) *Error {
	return &Error{
		Messages: messages,
	}
}

func (e *Error) Error() string {
	buff := bytes.NewBufferString("")
	for _, msg := range e.Messages {
		buff.WriteString(msg)
		buff.WriteString("\n")
	}
	return buff.String()
}
