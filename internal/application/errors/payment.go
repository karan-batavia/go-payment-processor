package errors

type Payment string

func (e Payment) Error() string {
	return string(e)
}
