package tests

import (
	"DatingApp/entities"
	"database/sql/driver"
)

var OrderCols = []string{"id", "user_id", "package_id", "name_package", "amount"}

var OrderListDummy = []entities.Order{
	entities.Order{
		UserId:      1,
		PackageId:   1,
		NamePackage: "Premium",
		Amount:      0,
	}}

func MappingOrderStore(item entities.Order, numId int) []driver.Value {
	var values []driver.Value

	values = append(values, numId)
	values = append(values, item.UserId)
	values = append(values, item.PackageId)
	values = append(values, item.NamePackage)
	values = append(values, item.Amount)

	return values
}
