package service

import (
	"DatingApp/exceptions"
	"DatingApp/repositories"
	"DatingApp/requests"
	"DatingApp/services"
	"DatingApp/tests"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http/httptest"
	"regexp"
	"testing"
)

func setupInitialOrder(db *gorm.DB) (services.OrderService, *gin.Context) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	repository := repositories.NewOrderRepository()
	repositoryPackage := repositories.NewPackageRepository()
	repositoryUser := repositories.NewUserRepository()
	service := services.NewOrderService(db, validator.New(), repository, repositoryPackage, repositoryUser)

	return service, c
}

// Create Order Test (Service)
func TestCreateOrder(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	service, ctx := setupInitialOrder(db)

	u1 := tests.UserListDummy[0]
	uRow1 := tests.MappingUserStore(u1, 1)

	p1 := tests.PackageListDummy[0]
	pRow1 := tests.MappingPackageStore(p1, 1)

	o1 := tests.OrderListDummy[0]
	oRow1 := tests.MappingOrderStore(o1, 1)
	argsCtInst := tests.AnyThingArgs(4)

	t.Run("Validation Failed", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.ValidationError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()

		service.CreateOrder(ctx, requests.CreateOrderRequest{})

		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("User Not Found", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.NotFoundError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()
		users := sqlmock.NewRows(tests.UserCols)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserSQL).WillReturnRows(users)
		mock.ExpectRollback()

		service.CreateOrder(ctx, requests.CreateOrderRequest{
			UserId:    1,
			PackageId: 1,
			Amount:    0,
		})

		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Package Not Found", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.NotFoundError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()
		users := sqlmock.NewRows(tests.UserCols).AddRow(uRow1...)
		packages := sqlmock.NewRows(tests.PackageCols)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserSQL).WillReturnRows(users)
		mock.ExpectQuery(tests.SelectOnePackageSQL).WillReturnRows(packages)
		mock.ExpectRollback()

		service.CreateOrder(ctx, requests.CreateOrderRequest{
			UserId:    1,
			PackageId: 1,
			Amount:    0,
		})

		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Package Already Purchased", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.ConflictError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()
		users := sqlmock.NewRows(tests.UserCols).AddRow(uRow1...)
		packages := sqlmock.NewRows(tests.PackageCols).AddRow(pRow1...)
		orders := sqlmock.NewRows(tests.OrderCols).AddRow(oRow1...)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserSQL).WillReturnRows(users)
		mock.ExpectQuery(tests.SelectOnePackageSQL).WillReturnRows(packages)
		mock.ExpectQuery(tests.SelectOneOrderUserSQL).WillReturnRows(orders)
		mock.ExpectRollback()

		service.CreateOrder(ctx, requests.CreateOrderRequest{
			UserId:    1,
			PackageId: 1,
			Amount:    0,
		})

		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Success Order", func(t *testing.T) {
		users := sqlmock.NewRows(tests.UserCols).AddRow(uRow1...)
		packages := sqlmock.NewRows(tests.PackageCols).AddRow(pRow1...)
		orders := sqlmock.NewRows(tests.OrderCols)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserSQL).WillReturnRows(users)
		mock.ExpectQuery(tests.SelectOnePackageSQL).WillReturnRows(packages)
		mock.ExpectQuery(tests.SelectOneOrderUserSQL).WillReturnRows(orders)
		mock.ExpectExec(tests.InsertOrderSQL).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(tests.UpdateUserSQL).
			WithArgs(tests.AnyThingArgs(7)...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		data := service.CreateOrder(ctx, requests.CreateOrderRequest{
			UserId:    1,
			PackageId: 1,
			Amount:    0,
		})
		t.Logf("%+v", data)

		assert.NotNil(t, data)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Get One Order Test (Service)
func TestGetOneOrder(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	service, ctx := setupInitialOrder(db)

	o1 := tests.OrderListDummy[0]
	oRow1 := tests.MappingOrderStore(o1, 1)

	t.Run("Order Not Found", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.NotFoundError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()

		order := sqlmock.NewRows(tests.OrderCols)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneOrderSQL).WithArgs(1).WillReturnRows(order)
		mock.ExpectRollback()

		service.GetOneOrder(ctx, 1)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Success Test", func(t *testing.T) {
		order := sqlmock.NewRows(tests.OrderCols).
			AddRow(oRow1...)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneOrderSQL).WithArgs(1).WillReturnRows(order)
		mock.ExpectCommit()

		data := service.GetOneOrder(ctx, 1)

		t.Logf("%+v", data)
		assert.NotNil(t, data)
		assert.Equal(t, 1, data.Id)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Delete One Order Test (Service)
func TestDeleteOneOrder(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	service, ctx := setupInitialOrder(db)

	o1 := tests.OrderListDummy[0]
	oRow1 := tests.MappingOrderStore(o1, 1)
	argsCtInst := tests.AnyThingArgs(2)

	t.Run("Order Not Found", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.NotFoundError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()

		order := sqlmock.NewRows(tests.OrderCols)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneOrderSQL).WithArgs(1).WillReturnRows(order)
		mock.ExpectRollback()

		service.DeleteOneOrder(ctx, 1)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Success Test", func(t *testing.T) {
		order := sqlmock.NewRows(tests.OrderCols).
			AddRow(oRow1...)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneOrderSQL).WithArgs(1).WillReturnRows(order)
		mock.ExpectExec(regexp.QuoteMeta(tests.DeleteOneOrderSQL)).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		service.DeleteOneOrder(ctx, 1)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
