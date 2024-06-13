package requests

type CreateOrderRequest struct {
	UserId    int `validate:"required" json:"user_id"`
	PackageId int `validate:"required" json:"package_id"`
	Amount    int `json:"amount"`
}
