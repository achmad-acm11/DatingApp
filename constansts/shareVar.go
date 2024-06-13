package constants

import "github.com/sirupsen/logrus"

const UserNotFound = "user not found"
const UserMatchNotFound = "User Match not found"
const PaymentRequired = "You need to pay to access this content."
const PackageNotFound = "Package not found"
const OrderNotFound = "Order not found"
const OrderInsufficient = "Order Insufficient"
const PackagePurchased = "The package has been purchased"
const UserConflict = "user conflict"
const LoginFailed = "Login failed"
const UserDeleted = "User deleted"
const PackageDeleted = "Package deleted"
const OrderDeleted = "Order deleted"

type LayerName string

const (
	Service LayerName = "Service"
)

type NameEntity string

const (
	User    NameEntity = "User"
	Order   NameEntity = "Order"
	Package NameEntity = "Package"
)

var Logger *logrus.Logger

type LogBasicKeyField string

const (
	Request     LogBasicKeyField = "request"
	Message     LogBasicKeyField = "message"
	Response    LogBasicKeyField = "response"
	TimeElapsed LogBasicKeyField = "time_elapsed"
)
