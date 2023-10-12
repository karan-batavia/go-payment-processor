package errors

type Validation string

func (e Validation) Error() string {
	return string(e)
}
