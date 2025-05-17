package errors

type DatabaseNotFoundError struct {
	Message string
}

func NewDatabaseNotFoundError(message string) *DatabaseNotFoundError {
	return &DatabaseNotFoundError{Message: message}
}

func (e *DatabaseNotFoundError) Error() string {
	return e.Message
}
