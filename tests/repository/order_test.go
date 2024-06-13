package repository

import (
	"DatingApp/repositories"
	"DatingApp/tests"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"regexp"
	"testing"
)

func setupInitialOrder() (repositories.OrderRepository, *gin.Context) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	repository := repositories.NewOrderRepository()

	return repository, c
}

// Get All Order Test (Repository)
func TestGetAllOrder(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	repo, ctx := setupInitialOrder()

	o1 := tests.OrderListDummy[0]
	oRow1 := tests.MappingOrderStore(o1, 1)

	t.Run("Success Test", func(t *testing.T) {
		orders := sqlmock.NewRows(tests.OrderCols).
			AddRow(oRow1...)

		mock.ExpectQuery(tests.SelectAllOrderSQL).WillReturnRows(orders)

		data := repo.GetAll(ctx, db)

		t.Logf("%+v", data)

		assert.Len(t, data, 1)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Create Order Test (Repository)
func TestCreateOrder(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()
	repo, ctx := setupInitialOrder()

	o1 := tests.OrderListDummy[0]
	argsCtInst := tests.AnyThingArgs(4)

	t.Run("Success Test", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(tests.InsertOrderSQL).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		data := repo.Create(ctx, db, o1)

		t.Logf("%+v", data)

		assert.NotNil(t, data)
		assert.Equal(t, 1, data.Id)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Update Order Test (Repository)
func TestUpdateOrder(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()
	repo, ctx := setupInitialOrder()

	o1 := tests.OrderListDummy[0]
	o1.Id = 1
	argsCtInst := tests.AnyThingArgs(5)

	t.Run("Success Test", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(tests.UpdateOrderSQL).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		data := repo.Update(ctx, db, o1)

		t.Logf("%+v", data)

		assert.NotNil(t, data)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Delete One Order Test (Repository)
func TestDeleteOneOrder(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	repo, ctx := setupInitialOrder()

	o1 := tests.OrderListDummy[0]
	o1.Id = 1
	argsCtInst := tests.AnyThingArgs(2)

	t.Run("Success Test", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(tests.DeleteOneOrderSQL)).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo.DeleteOne(ctx, db, o1)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Get One Order Test (Repository)
func TestGetOneOrder(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	repo, ctx := setupInitialOrder()

	o1 := tests.OrderListDummy[0]
	oRow1 := tests.MappingOrderStore(o1, 1)

	t.Run("Success Test", func(t *testing.T) {
		order := sqlmock.NewRows(tests.OrderCols).
			AddRow(oRow1...)

		mock.ExpectQuery(tests.SelectOneOrderSQL).WithArgs(1).WillReturnRows(order)

		data := repo.GetOneById(ctx, db, 1)
		assert.Equal(t, 1, data.Id)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Get One Order User Test (Repository)
func TestGetOneOrderUser(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	repo, ctx := setupInitialOrder()

	o1 := tests.OrderListDummy[0]
	oRow1 := tests.MappingOrderStore(o1, 1)

	t.Run("Success Test", func(t *testing.T) {
		order := sqlmock.NewRows(tests.OrderCols).
			AddRow(oRow1...)

		mock.ExpectQuery(tests.SelectOneOrderUserSQL).
			WithArgs(1, 1).WillReturnRows(order)

		data := repo.GetOrderUser(ctx, db, 1, 1)
		assert.Equal(t, 1, data.Id)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
