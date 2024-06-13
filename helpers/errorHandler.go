package helpers

import (
	"DatingApp/exceptions"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func ErrorHandler(err error) {
	if err != nil {
		fmt.Printf("%+v", err)
		panic(exceptions.NewInternalServerError(err.Error()))
	}
}

func customMessage(field string, typeMessage string, param string) string {
	var msg string
	switch typeMessage {
	case "required":
		msg = fmt.Sprintf("%s %s", field, "field is required")
	case "oneof":
		msg = fmt.Sprintf("%s %s", field, "field is oneof "+param)
	case "e164":
		msg = fmt.Sprintf("%s %s", field, "field invalid phone number format, format example: +62***********")
	case "OTP":
		msg = fmt.Sprintf("%s %s", field, "field invalid OTP format, number format and lenght OTP is "+param)
	case "email":
		msg = fmt.Sprintf("%s %s", field, "field invalid Email format")
	case "pin":
		msg = fmt.Sprintf("%s %s", field, "field invalid Pin format")
	}
	return msg
}

func ErrorHandlerValidator(err error) {
	if err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			msg := customMessage(strings.ToLower(e.Field()), e.Tag(), e.Param())
			errors[strings.ToLower(e.Field())] = msg
		}
		panic(exceptions.NewValidationError(errors))
	}
}
