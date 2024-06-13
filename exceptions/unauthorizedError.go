package exceptions

type UnauthorizedError struct {
	Message string `json:"message"`
}

func NewUnauthorizedError(message string) UnauthorizedError {
	return UnauthorizedError{
		Message: message,
	}
}
