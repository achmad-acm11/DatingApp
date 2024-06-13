package tests

import (
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type AnyThing struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyThing) Match(v driver.Value) bool {
	return true
}

func AnyThingArgs(size int) []driver.Value {
	var values []driver.Value
	for i := 0; i < size; i++ {
		values = append(values, AnyThing{})
	}

	return values
}

func DbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// a SELECT VERSION() query will be run when gorm opens the database
	// so we need to expect that here
	columns := []string{"version"}
	mock.ExpectQuery("SELECT VERSION()").WithArgs().WillReturnRows(
		mock.NewRows(columns).FromCSVString("1"),
	)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:       sqlDB,
		DriverName: "mysql",
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return sqlDB, gormDB, mock
}

// SQL Package
const InsertPackageSQL = "INSERT INTO `product_packages`"
const UpdatePackageSQL = "UPDATE `product_packages`"
const SelectOnePackageSQL = "SELECT (.+) FROM `product_packages` WHERE id = \\?"
const SelectAllPackageSQL = "SELECT (.+) FROM `product_packages`"
const DeleteOnePackageSQL = "UPDATE `product_packages` SET `deleted_at`=? WHERE `product_packages`.`id` = ? AND `product_packages`.`deleted_at` IS NULL"

// SQL Order
const InsertOrderSQL = "INSERT INTO `orders`"
const UpdateOrderSQL = "UPDATE `orders`"
const SelectOneOrderSQL = "SELECT (.+) FROM `orders` WHERE id = \\?"
const SelectOneOrderUserSQL = "SELECT (.+) FROM `orders` WHERE package_id = \\? AND user_id = \\?"
const SelectAllOrderSQL = "SELECT (.+) FROM `orders`"
const DeleteOneOrderSQL = "UPDATE `orders` SET `deleted_at`=? WHERE `orders`.`id` = ? AND `orders`.`deleted_at` IS NULL"

// SQL Match Queue
const InsertMatchQueueSQL = "INSERT INTO `match_queues`"
const UpdateMatchQueueSQL = "UPDATE `match_queues`"
const SelectOneMatchQueueSQL = "SELECT (.+) FROM `match_queues` WHERE user_id = \\?"
const SelectOneMatchQueueCurrentSQL = "SELECT (.+) FROM `match_queues` WHERE user_id = \\? AND date = \\?"
const SelectAllMatchQueueSQL = "SELECT (.+) FROM `match_queues`"
const DeleteOneMatchQueueSQL = "UPDATE `match_queues` SET `deleted_at`=? WHERE `match_queues`.`id` = ? AND `match_queues`.`deleted_at` IS NULL"

// SQL User
const SelectAllUserSQL = "SELECT (.+) FROM `users`"
const InsertUserSQL = "INSERT INTO `users`"
const SelectOneUserByEmailSQL = "SELECT (.+) FROM `users` WHERE email = \\?"
const UpdateUserSQL = "UPDATE `users`"
const SelectOneUserSQL = "SELECT (.+) FROM `users` WHERE id = \\?"
const SelectAllTenMatchSQL = "SELECT * FROM `users` WHERE id != ? AND gender != ? AND `users`.`deleted_at` IS NULL LIMIT 10"
const DeleteOneUserSQL = "UPDATE `users` SET `deleted_at`=? WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL"
