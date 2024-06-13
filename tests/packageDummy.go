package tests

import (
	"DatingApp/entities"
	"database/sql/driver"
)

var PackageCols = []string{"id", "name_package", "amount"}

var PackageListDummy = []entities.ProductPackage{
	entities.ProductPackage{
		NamePackage: "Premium",
		Amount:      0,
	}, entities.ProductPackage{
		NamePackage: "Verified",
		Amount:      0,
	}}

func MappingPackageStore(item entities.ProductPackage, numId int) []driver.Value {
	var values []driver.Value

	values = append(values, numId)
	values = append(values, item.NamePackage)
	values = append(values, item.Amount)

	return values
}
