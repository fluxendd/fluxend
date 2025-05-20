package errors

type ForbiddenError struct {
	Message string
}

func NewForbiddenError(message string) *ForbiddenError {
	return &ForbiddenError{Message: message}
}

func (e *ForbiddenError) Error() string {
	return e.Message
}
