package requests

type CreateUserRequest struct {
	Name        string `validate:"required" json:"name"`
	Gender      string `validate:"required,oneof=male female'" json:"gender"`
	PhoneNumber string `validate:"required,e164" json:"phone_number"`
	Email       string `validate:"required,email" json:"email"`
	Password    string `validate:"required" json:"password"`
}

type SignUpUserRequest struct {
	Name        string `validate:"required" json:"name"`
	PhoneNumber string `validate:"required,e164" json:"phone_number"`
	Email       string `validate:"required,email" json:"email"`
	Gender      string `validate:"required,oneof=male female'" json:"gender"`
	Password    string `validate:"required" json:"password"`
}

type LoginRequest struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required" json:"password"`
}

//type CheckOTPRequest struct {
//	PhoneNumber string `validate:"required,e164" json:"phone_number"`
//	Otp         int    `validate:"required,OTP=4" json:"otp"`
//}

//type CheckUserRequest struct {
//	Email string `validate:"required" json:"email"`
//}

//type SendOTPUserRequest struct {
//	Name        string `validate:"required" json:"name"`
//	Email       string `validate:"required" json:"email"`
//	Password    string `validate:"required" json:"password"`
//	PhoneNumber string `validate:"required,e164" json:"phone_number"`
//}
