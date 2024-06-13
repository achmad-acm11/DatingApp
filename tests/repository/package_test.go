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

func setupInitialPackage() (repositories.PackageRepository, *gin.Context) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	repository := repositories.NewPackageRepository()

	return repository, c
}

// Get All Package Test (Repository)
func TestGetAllPackage(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	repo, ctx := setupInitialPackage()

	p1 := tests.PackageListDummy[0]
	pRow1 := tests.MappingPackageStore(p1, 1)

	t.Run("Success Test", func(t *testing.T) {
		packages := sqlmock.NewRows(tests.PackageCols).
			AddRow(pRow1...)

		mock.ExpectQuery(tests.SelectAllPackageSQL).WillReturnRows(packages)

		data := repo.GetAll(ctx, db)

		t.Logf("%+v", data)

		assert.Len(t, data, 1)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Create Package Test (Repository)
func TestCreatePackage(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()
	repo, ctx := setupInitialPackage()

	p1 := tests.PackageListDummy[0]
	argsCtInst := tests.AnyThingArgs(2)

	t.Run("Success Test", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(tests.InsertPackageSQL).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		data := repo.Create(ctx, db, p1)

		t.Logf("%+v", data)

		assert.NotNil(t, data)
		assert.Equal(t, 1, data.Id)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Update Package Test (Repository)
func TestUpdatePackage(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()
	repo, ctx := setupInitialPackage()

	p1 := tests.PackageListDummy[0]
	p1.Id = 1
	argsCtInst := tests.AnyThingArgs(3)

	t.Run("Success Test", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(tests.UpdatePackageSQL).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		data := repo.Update(ctx, db, p1)

		t.Logf("%+v", data)

		assert.NotNil(t, data)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Get One Package Test (Repository)
func TestGetOnePackage(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	repo, ctx := setupInitialPackage()

	p1 := tests.PackageListDummy[0]
	pRow1 := tests.MappingPackageStore(p1, 1)

	t.Run("Success Test", func(t *testing.T) {
		product := sqlmock.NewRows(tests.PackageCols).
			AddRow(pRow1...)

		mock.ExpectQuery(tests.SelectOnePackageSQL).WithArgs(1).WillReturnRows(product)

		data := repo.GetOneById(ctx, db, 1)
		assert.Equal(t, 1, data.Id)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Delete One Package Test (Repository)
func TestDeleteOnePackage(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	repo, ctx := setupInitialPackage()

	p1 := tests.PackageListDummy[0]
	p1.Id = 1
	argsCtInst := tests.AnyThingArgs(2)

	t.Run("Success Test", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(tests.DeleteOnePackageSQL)).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo.DeleteOne(ctx, db, p1)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
