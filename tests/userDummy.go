package tests

import (
	"DatingApp/entities"
	"database/sql/driver"
)

var UserCols = []string{"id", "name", "email", "gender", "phone_number", "is_premium", "verified", "password"}

var UserListDummy = []entities.User{entities.User{
	Name:        "Budi Hernandex",
	Email:       "budi@mail.com",
	Gender:      "male",
	PhoneNumber: "+62877XXXXXXXX",
	IsPremium:   false,
	Verified:    false,
}, entities.User{
	Name:        "Muthia",
	Email:       "muthia@mail.com",
	Gender:      "female",
	PhoneNumber: "+62877XXXXXXXX",
	IsPremium:   false,
	Verified:    false,
}}

func MappingUserStore(item entities.User, numId int) []driver.Value {
	var values []driver.Value

	values = append(values, numId)
	values = append(values, item.Name)
	values = append(values, item.Email)
	values = append(values, item.Gender)
	values = append(values, item.PhoneNumber)
	values = append(values, item.IsPremium)
	values = append(values, item.Verified)
	values = append(values, item.Password)

	return values
}
