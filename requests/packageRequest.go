package requests

type CreatePackageRequest struct {
	NamePackage string `validate:"required" json:"name_package"`
	Amount      int    `json:"amount"`
}

type UpdatePackageRequest struct {
	NamePackage string `validate:"required" json:"name_package"`
	Amount      int    `json:"amount"`
}
