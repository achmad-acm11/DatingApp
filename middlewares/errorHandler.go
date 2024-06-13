package middlewares

import (
	"DatingApp/exceptions"
	"DatingApp/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//if constant.Logger != nil {
				//	constant.Logger.Error(fmt.Sprintf("%+v", err))
				//}

				var errType interface{}
				var message string
				var status string
				var code int

				switch err.(type) {
				case exceptions.InternalServerError:
					errType = err.(exceptions.InternalServerError)
					message = "Internal Server Error"
					status = "error"
					code = http.StatusInternalServerError
				case exceptions.NotFoundError:
					errType = err.(exceptions.NotFoundError)
					message = "Not Found"
					status = "error"
					code = http.StatusNotFound
				case exceptions.ValidationError:
					errType = err.(exceptions.ValidationError)
					message = "Bad Request"
					status = "error"
					code = http.StatusBadRequest
				case exceptions.BadRequestError:
					errType = err.(exceptions.BadRequestError)
					message = "Bad Request"
					status = "error"
					code = http.StatusBadRequest
				case exceptions.ConflictError:
					errType = err.(exceptions.ConflictError)
					message = "Conflict"
					status = "error"
					code = http.StatusConflict
				case exceptions.UnauthorizedError:
					errType = err.(exceptions.UnauthorizedError)
					message = "Unauthorized"
					status = "error"
					code = http.StatusUnauthorized
				case exceptions.PaymentRequiredError:
					errType = err.(exceptions.PaymentRequiredError)
					message = "Payment Required"
					status = "error"
					code = http.StatusPaymentRequired
				}

				context.JSON(code, helpers.APIResponse(message, status, code, errType))
			}
		}()
		context.Next()
	}
}
