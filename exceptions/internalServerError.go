package exceptions

type InternalServerError struct {
	Message string `json:"message"`
}

func NewInternalServerError(message string) InternalServerError {
	return InternalServerError{
		Message: message,
	}
}
