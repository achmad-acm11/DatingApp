package requests

type CreatePackageRequest struct {
	NamePackage string `validate:"required" json:"name_package"`
	Amount      int    `validate:"required" json:"amount"`
}

type UpdatePackageRequest struct {
	NamePackage string `validate:"required" json:"name_package"`
	Amount      int    `validate:"required" json:"amount"`
}
