package constants

const UserNotFound = "user not found"
const UserCoupleNotFound = "User Couple not found"
const PaymentRequired = "You need to pay to access this content."
const PackageNotFound = "Package not found"
const OrderNotFound = "Order not found"
const OrderInsufficient = "Order Insufficient"
const PackagePurchased = "The package has been purchased"
const ReferralNotFound = "referral not found"
const ReferralConflict = "referral conflict"
const EmailConflict = "email conflict"
const UserConflict = "user conflict"
const LoginFailed = "Login failed"
const AccountIsNotVerify = "Account is not verify"
const CurrentPasswordWrong = "Current password is wrong"
const TokenPasswordWrong = "Token password is wrong"
const BadReqPassword = "New Password is not match with Verify Password"
const UserDeleted = "User deleted"
const PackageDeleted = "Package deleted"
const OrderDeleted = "Order deleted"
const UserForgetPassword = "User Changed Password"
const UserReset = "User reset password"
const ReferralDeleted = "Referral deleted"
const UserCreated = "User created"
const OTPNotMatch = "OTP is not match"
const OTPMessage = "Hi %s\nSelamat Datang di Qonstanta\n\nSaat ini Kamu mendaftar sebagai Siswa.\nqonstanta.id\n\nEmail: %s\n\nPassword: %s\n\nKode Verifikasi: %d\nDemi keamanan Anda, jangan bagikan kode ini.\n\nTerima Kasih\nminQ 😊"
const OTPPasswordMessage = "Hi %s, Kamu telah melakukan request pergantion password akun di qonstanta.id\n\nKode Verifikasi: %d\n\n Silahkan masukkan kode diatas sebagai konfirmasi\nDemi keamanan Anda, jangan bagikan kode ini.\n\nTerima Kasih\nminQ 😊"

type LayerName string

const (
	Controller LayerName = "Controller"
	Service    LayerName = "Service"
	Repository LayerName = "Repository"
)

type NameEntity string

const (
	User NameEntity = "User"
)

//var Logger *logrus.Logger

type LogBasicKeyField string

const (
	Request     LogBasicKeyField = "request"
	Message     LogBasicKeyField = "message"
	Response    LogBasicKeyField = "response"
	TimeElapsed LogBasicKeyField = "time_elapsed"
)
