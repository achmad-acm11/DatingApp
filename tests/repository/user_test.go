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

func setupInitialUser() (repositories.UserRepository, *gin.Context) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	repository := repositories.NewUserRepository()

	return repository, c
}

// Get All User Test (Repository)
func TestGetAllUser(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	repo, ctx := setupInitialUser()

	u1 := tests.UserListDummy[0]
	uRow1 := tests.MappingUserStore(u1, 1)

	t.Run("Success Test", func(t *testing.T) {
		users := sqlmock.NewRows(tests.UserCols).
			AddRow(uRow1...)

		mock.ExpectQuery(tests.SelectAllUserSQL).WillReturnRows(users)

		data := repo.GetAll(ctx, db)

		t.Logf("%+v", data)

		assert.Len(t, data, 1)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Create User Test (Repository)
func TestCreateUser(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()
	repo, ctx := setupInitialUser()

	u1 := tests.UserListDummy[0]
	argsCtInst := tests.AnyThingArgs(7)

	t.Run("Success Test", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(tests.InsertUserSQL).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		data := repo.Create(ctx, db, u1)

		t.Logf("%+v", data)

		assert.NotNil(t, data)
		assert.Equal(t, 1, data.Id)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Get All User Ten Match Test (Repository)
func TestGetAllTenMatchUser(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	repo, ctx := setupInitialUser()

	u1 := tests.UserListDummy[0]
	uRow1 := tests.MappingUserStore(u1, 1)

	t.Run("Success Test", func(t *testing.T) {
		users := sqlmock.NewRows(tests.UserCols).
			AddRow(uRow1...)

		mock.ExpectQuery(regexp.QuoteMeta(tests.SelectAllTenMatchSQL)).WillReturnRows(users)

		data := repo.GetAllTenMatch(ctx, db, 1, "male", []string{})

		t.Logf("%+v", data)

		assert.Len(t, data, 1)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Success Test Not in userIds", func(t *testing.T) {
		users := sqlmock.NewRows(tests.UserCols).
			AddRow(uRow1...)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE id != ? AND gender != ? AND id NOT IN (?) AND `users`.`deleted_at` IS NULL LIMIT 10")).WillReturnRows(users)

		data := repo.GetAllTenMatch(ctx, db, 1, "male", []string{"1"})

		t.Logf("%+v", data)

		assert.Len(t, data, 1)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Get One User By Email Test (Repository)
func TestGetOneUserByEmail(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	repo, ctx := setupInitialUser()

	u1 := tests.UserListDummy[0]
	uRow1 := tests.MappingUserStore(u1, 1)

	t.Run("Success Test", func(t *testing.T) {
		user := sqlmock.NewRows(tests.UserCols).
			AddRow(uRow1...)

		mock.ExpectQuery(tests.SelectOneUserByEmailSQL).WithArgs("budi@mail.com").WillReturnRows(user)

		data := repo.GetOneByEmail(ctx, db, "budi@mail.com")
		assert.Equal(t, 1, data.Id)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Get One User Test (Repository)
func TestGetOneUser(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	repo, ctx := setupInitialUser()

	u1 := tests.UserListDummy[0]
	uRow1 := tests.MappingUserStore(u1, 1)

	t.Run("Success Test", func(t *testing.T) {
		user := sqlmock.NewRows(tests.UserCols).
			AddRow(uRow1...)

		mock.ExpectQuery(tests.SelectOneUserSQL).WithArgs(1).WillReturnRows(user)

		data := repo.GetOneById(ctx, db, 1)
		assert.Equal(t, 1, data.Id)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Update User Test (Repository)
func TestUpdateUser(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()
	repo, ctx := setupInitialUser()

	u1 := tests.UserListDummy[0]
	u1.Id = 1
	argsCtInst := tests.AnyThingArgs(6)

	t.Run("Success Test", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(tests.UpdateUserSQL).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		data := repo.Update(ctx, db, u1)

		t.Logf("%+v", data)

		assert.NotNil(t, data)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Delete One User Test (Repository)
func TestDeleteOneUser(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	repo, ctx := setupInitialUser()

	u1 := tests.UserListDummy[0]
	u1.Id = 1
	argsCtInst := tests.AnyThingArgs(2)

	t.Run("Success Test", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(tests.DeleteOneUserSQL)).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo.DeleteOne(ctx, db, u1)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
