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

func setupInitialUser(db *gorm.DB) (services.UserService, *gin.Context) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("userId", float64(1))

	repository := repositories.NewUserRepository()
	repositoryMatchQueue := repositories.NewMatchQueueRepository(repository)
	service := services.NewUserService(db, validator.New(), repository, repositoryMatchQueue)

	return service, c
}

// Get My User Test (Service)
func TestGetMyUser(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	service, ctx := setupInitialUser(db)

	u1 := tests.UserListDummy[0]
	uRow1 := tests.MappingUserStore(u1, 1)

	t.Run("Success Test", func(t *testing.T) {
		user := sqlmock.NewRows(tests.UserCols).
			AddRow(uRow1...)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserSQL).WithArgs(1).WillReturnRows(user)
		mock.ExpectCommit()

		data := service.GetMyUser(ctx)

		t.Logf("%+v", data)
		assert.NotNil(t, data)
		assert.Equal(t, 1, data.Id)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Get One User Test (Service)
func TestGetOneUser(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	service, ctx := setupInitialUser(db)

	u1 := tests.UserListDummy[0]
	uRow1 := tests.MappingUserStore(u1, 1)

	t.Run("User Not Found", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.NotFoundError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()

		user := sqlmock.NewRows(tests.UserCols)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserSQL).WithArgs(1).WillReturnRows(user)
		mock.ExpectRollback()

		service.GetOneUser(ctx, 1)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Success Test", func(t *testing.T) {
		user := sqlmock.NewRows(tests.UserCols).
			AddRow(uRow1...)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserSQL).WithArgs(1).WillReturnRows(user)
		mock.ExpectCommit()

		data := service.GetOneUser(ctx, 1)

		t.Logf("%+v", data)
		assert.NotNil(t, data)
		assert.Equal(t, 1, data.Id)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Get All User Test (Service)
func TestGetAllUser(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	service, ctx := setupInitialUser(db)

	u1 := tests.UserListDummy[0]
	uRow1 := tests.MappingUserStore(u1, 1)

	t.Run("Success Test", func(t *testing.T) {
		users := sqlmock.NewRows(tests.UserCols).
			AddRow(uRow1...)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectAllUserSQL).WillReturnRows(users)
		mock.ExpectCommit()

		data := service.GetAllUser(ctx)

		t.Logf("%+v", data)
		assert.NotNil(t, data)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// SignUp User Test (Service)
func TestSignUpUser(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	service, ctx := setupInitialUser(db)

	u1 := tests.UserListDummy[0]
	uRow1 := tests.MappingUserStore(u1, 1)
	argsCtInst := tests.AnyThingArgs(7)

	t.Run("Validation Failed", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.ValidationError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()

		service.SignUpUser(ctx, requests.SignUpUserRequest{})
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("User Conflict", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.ConflictError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()
		users := sqlmock.NewRows(tests.UserCols).AddRow(uRow1...)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserByEmailSQL).WillReturnRows(users)
		mock.ExpectRollback()

		service.SignUpUser(ctx, requests.SignUpUserRequest{
			Name:        u1.Name,
			PhoneNumber: "+6289988776655",
			Email:       u1.Email,
			Gender:      u1.Gender,
			Password:    "123",
		})
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Success Test", func(t *testing.T) {
		users := sqlmock.NewRows(tests.UserCols)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserByEmailSQL).WillReturnRows(users)
		mock.ExpectExec(tests.InsertUserSQL).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		data := service.SignUpUser(ctx, requests.SignUpUserRequest{
			Name:        u1.Name,
			PhoneNumber: "+6289988776655",
			Email:       u1.Email,
			Gender:      u1.Gender,
			Password:    "123",
		})

		t.Logf("%+v", data)
		assert.NotNil(t, data)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Login User Test (Service)
func TestLoginUser(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	service, ctx := setupInitialUser(db)

	u1 := tests.UserListDummy[0]
	u1.Password = "$2a$14$iUTzLp.KG7.HyjeQtBSXHepUsxIjflLRgZdeItj/SbSe5LzoUEXzq"
	uRow1 := tests.MappingUserStore(u1, 1)

	t.Run("Validation Failed", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.ValidationError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()

		service.Login(ctx, requests.LoginRequest{})
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("User not found", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.UnauthorizedError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()
		users := sqlmock.NewRows(tests.UserCols)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserByEmailSQL).WillReturnRows(users)
		mock.ExpectRollback()

		service.Login(ctx, requests.LoginRequest{
			Email:    u1.Email,
			Password: "123",
		})
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Wrong password", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.UnauthorizedError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()
		users := sqlmock.NewRows(tests.UserCols).AddRow(uRow1...)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserByEmailSQL).WillReturnRows(users)
		mock.ExpectRollback()

		service.Login(ctx, requests.LoginRequest{
			Email:    u1.Email,
			Password: "1",
		})
		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Success Test", func(t *testing.T) {
		users := sqlmock.NewRows(tests.UserCols).AddRow(uRow1...)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserByEmailSQL).WillReturnRows(users)
		mock.ExpectCommit()

		data := service.Login(ctx, requests.LoginRequest{
			Email:    u1.Email,
			Password: "123",
		})

		t.Logf("%+v", data)
		assert.NotNil(t, data)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Delete One User Test (Service)
func TestDeleteOneUser(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	service, ctx := setupInitialUser(db)

	u1 := tests.UserListDummy[0]
	uRow1 := tests.MappingUserStore(u1, 1)
	argsCtInst := tests.AnyThingArgs(2)

	t.Run("User Not Found", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.NotFoundError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()

		user := sqlmock.NewRows(tests.UserCols)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserSQL).WithArgs(1).WillReturnRows(user)
		mock.ExpectRollback()

		service.DeleteOneUser(ctx, 1)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Success Test", func(t *testing.T) {
		user := sqlmock.NewRows(tests.UserCols).
			AddRow(uRow1...)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserSQL).WithArgs(1).WillReturnRows(user)
		mock.ExpectExec(regexp.QuoteMeta(tests.DeleteOneUserSQL)).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		service.DeleteOneUser(ctx, 1)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Get All Match User Test (Service)
func TestGetAllMatchUser(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	service, ctx := setupInitialUser(db)

	u1 := tests.UserListDummy[0]
	uRow1 := tests.MappingUserStore(u1, 1)
	u2 := tests.UserListDummy[1]
	uRow2 := tests.MappingUserStore(u2, 2)

	m1 := tests.MatchQueueListDummy[0]
	mRow1 := tests.MappingMatchQueueStore(m1, 1)

	t.Run("User Not Found", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.NotFoundError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()

		user := sqlmock.NewRows(tests.UserCols)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserSQL).WithArgs(1).WillReturnRows(user)
		mock.ExpectRollback()

		service.GetAllMatch(ctx)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Success Test", func(t *testing.T) {
		queue := sqlmock.NewRows(tests.MatchQueueCols).
			AddRow(mRow1...)
		user := sqlmock.NewRows(tests.UserCols).
			AddRow(uRow1...)
		user2 := sqlmock.NewRows(tests.UserCols).
			AddRow(uRow2...)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserSQL).WithArgs(1).WillReturnRows(user)
		mock.ExpectQuery(tests.SelectOneMatchQueueSQL).WithArgs(1).WillReturnRows(queue)
		mock.ExpectQuery(tests.SelectOneUserSQL).WithArgs(2).WillReturnRows(user2)
		mock.ExpectCommit()

		service.GetAllMatch(ctx)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Like Match User Test (Service)
func TestLikeMatchUser(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	service, ctx := setupInitialUser(db)

	u1 := tests.UserListDummy[0]
	uRow1 := tests.MappingUserStore(u1, 1)
	u2 := tests.UserListDummy[1]
	uRow2 := tests.MappingUserStore(u2, 2)

	m1 := tests.MatchQueueListDummy[0]
	mRow1 := tests.MappingMatchQueueStore(m1, 1)
	m1.Id = 1
	argsCtInst := tests.AnyThingArgs(7)

	t.Run("User Not Found", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.NotFoundError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()

		user := sqlmock.NewRows(tests.UserCols)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserSQL).WithArgs(1).WillReturnRows(user)
		mock.ExpectRollback()

		service.LikeMatch(ctx)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("User Match Not Found", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.NotFoundError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()

		user := sqlmock.NewRows(tests.UserCols).AddRow(uRow1...)
		queue := sqlmock.NewRows(tests.MatchQueueCols)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserSQL).WithArgs(1).WillReturnRows(user)
		mock.ExpectQuery(tests.SelectOneMatchQueueCurrentSQL).WithArgs(tests.AnyThing{}, tests.AnyThing{}).WillReturnRows(queue)
		mock.ExpectRollback()

		service.LikeMatch(ctx)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Success Test", func(t *testing.T) {
		queue := sqlmock.NewRows(tests.MatchQueueCols).
			AddRow(mRow1...)
		user := sqlmock.NewRows(tests.UserCols).
			AddRow(uRow1...)
		user2 := sqlmock.NewRows(tests.UserCols).
			AddRow(uRow2...)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserSQL).WithArgs(1).WillReturnRows(user)
		mock.ExpectQuery(tests.SelectOneMatchQueueCurrentSQL).WithArgs(tests.AnyThing{}, tests.AnyThing{}).WillReturnRows(queue)
		mock.ExpectQuery(tests.SelectOneUserSQL).WithArgs(3).WillReturnRows(user2)
		mock.ExpectExec(tests.UpdateMatchQueueSQL).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		service.LikeMatch(ctx)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Pass Match User Test (Service)
func TestPassMatchUser(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	service, ctx := setupInitialUser(db)

	u1 := tests.UserListDummy[0]
	uRow1 := tests.MappingUserStore(u1, 1)
	u2 := tests.UserListDummy[1]
	uRow2 := tests.MappingUserStore(u2, 2)

	m1 := tests.MatchQueueListDummy[0]
	mRow1 := tests.MappingMatchQueueStore(m1, 1)
	m1.Id = 1
	argsCtInst := tests.AnyThingArgs(7)

	t.Run("User Not Found", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.NotFoundError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()

		user := sqlmock.NewRows(tests.UserCols)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserSQL).WithArgs(1).WillReturnRows(user)
		mock.ExpectRollback()

		service.PassMatch(ctx)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("User Match Not Found", func(t *testing.T) {
		defer func() {
			err := recover()
			t.Logf("%+v", err)
			_, ok := err.(exceptions.NotFoundError)

			assert.True(t, ok)
			assert.NotNil(t, err)
		}()

		user := sqlmock.NewRows(tests.UserCols).AddRow(uRow1...)
		queue := sqlmock.NewRows(tests.MatchQueueCols)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserSQL).WithArgs(1).WillReturnRows(user)
		mock.ExpectQuery(tests.SelectOneMatchQueueCurrentSQL).WithArgs(tests.AnyThing{}, tests.AnyThing{}).WillReturnRows(queue)
		mock.ExpectRollback()

		service.PassMatch(ctx)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
	t.Run("Success Test", func(t *testing.T) {
		queue := sqlmock.NewRows(tests.MatchQueueCols).
			AddRow(mRow1...)
		user := sqlmock.NewRows(tests.UserCols).
			AddRow(uRow1...)
		user2 := sqlmock.NewRows(tests.UserCols).
			AddRow(uRow2...)

		mock.ExpectBegin()
		mock.ExpectQuery(tests.SelectOneUserSQL).WithArgs(1).WillReturnRows(user)
		mock.ExpectQuery(tests.SelectOneMatchQueueCurrentSQL).WithArgs(tests.AnyThing{}, tests.AnyThing{}).WillReturnRows(queue)
		mock.ExpectQuery(tests.SelectOneUserSQL).WithArgs(3).WillReturnRows(user2)
		mock.ExpectExec(tests.UpdateMatchQueueSQL).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		service.PassMatch(ctx)

		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
