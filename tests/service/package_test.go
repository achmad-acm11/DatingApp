package service

import (
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

func setupInitialPackage(db *gorm.DB) (services.PackageService, *gin.Context) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	repository := repositories.NewPackageRepository()
	service := services.NewPackageService(db, validator.New(), repository)

	return service, c
}

// Get All Package Test (Service)
func TestAllPackage(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	service, ctx := setupInitialPackage(db)

	p1 := tests.PackageListDummy[0]
	pRow1 := tests.MappingPackageStore(p1, 1)

	t.Run("Success Test", func(t *testing.T) {
		packages := sqlmock.NewRows(tests.PackageCols).
			AddRow(pRow1...)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectAllPackageSQL).WillReturnRows(packages)
		mock.ExpectCommit()

		data := service.GetAllPackage(ctx)
		t.Logf("%+v", data)
		assert.NotNil(t, data)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Create Package Test (Service)
func TestCreatePackage(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	service, ctx := setupInitialPackage(db)

	p1 := tests.PackageListDummy[0]
	argsCtInst := tests.AnyThingArgs(2)

	t.Run("Success Test", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(tests.InsertPackageSQL).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		data := service.CreatePackage(ctx, requests.CreatePackageRequest{
			NamePackage: p1.NamePackage,
			Amount:      p1.Amount,
		})

		t.Logf("%+v", data)
		assert.NotNil(t, data)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Update Package Test (Service)
func TestUpdatePackage(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	service, ctx := setupInitialPackage(db)

	p1 := tests.PackageListDummy[0]
	p1.Id = 1
	pRow1 := tests.MappingPackageStore(p1, 1)
	argsCtInst := tests.AnyThingArgs(4)

	t.Run("Success Test", func(t *testing.T) {
		product := sqlmock.NewRows(tests.PackageCols).
			AddRow(pRow1...)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOnePackageSQL).WithArgs(1).WillReturnRows(product)
		mock.ExpectExec(tests.UpdatePackageSQL).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		data := service.UpdatePackage(ctx, requests.UpdatePackageRequest{
			NamePackage: p1.NamePackage,
			Amount:      1,
		}, p1.Id)

		t.Logf("%+v", data)
		assert.NotNil(t, data)
		assert.Equal(t, 1, data.Amount)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Get One Package Test (Service)
func TestGetOnePackage(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	service, ctx := setupInitialPackage(db)

	p1 := tests.PackageListDummy[0]
	pRow1 := tests.MappingPackageStore(p1, 1)

	t.Run("Success Test", func(t *testing.T) {
		product := sqlmock.NewRows(tests.PackageCols).
			AddRow(pRow1...)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOnePackageSQL).WithArgs(1).WillReturnRows(product)
		mock.ExpectCommit()

		data := service.GetOnePackage(ctx, 1)

		t.Logf("%+v", data)
		assert.NotNil(t, data)
		assert.Equal(t, 1, data.Id)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Delete One Package Test (Service)
func TestDeleteOnePackage(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	service, ctx := setupInitialPackage(db)

	p1 := tests.PackageListDummy[0]
	pRow1 := tests.MappingPackageStore(p1, 1)
	argsCtInst := tests.AnyThingArgs(2)

	t.Run("Success Test", func(t *testing.T) {
		product := sqlmock.NewRows(tests.PackageCols).
			AddRow(pRow1...)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOnePackageSQL).WithArgs(1).WillReturnRows(product)
		mock.ExpectExec(regexp.QuoteMeta(tests.DeleteOnePackageSQL)).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		service.DeleteOnePackage(ctx, 1)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
