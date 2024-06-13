package exceptions

type BadRequestError struct {
	Message string `json:"message"`
}

func NewBadRequestError(message string) BadRequestError {
	return BadRequestError{
		Message: message,
	}
}
