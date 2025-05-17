package errors

type UnprocessableError struct {
	Message string
}

func NewUnprocessableError(message string) *UnprocessableError {
	return &UnprocessableError{Message: message}
}

func (e *UnprocessableError) Error() string {
	return e.Message
}
