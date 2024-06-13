package exceptions

type PaymentRequiredError struct {
	Message string `json:"message"`
}

func NewPaymentRequiredError(message string) PaymentRequiredError {
	return PaymentRequiredError{
		Message: message,
	}
}
