package repository

import (
	"DatingApp/helpers"
	"DatingApp/repositories"
	"DatingApp/tests"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"regexp"
	"testing"
)

func setupInitialMatchQueue() (repositories.MatchQueueRepository, *gin.Context) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	repoUser := repositories.NewUserRepository()
	repository := repositories.NewMatchQueueRepository(repoUser)

	return repository, c
}

// Get All Match Queue Test (Repository)
func TestGetAllMatchQueue(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	repo, ctx := setupInitialMatchQueue()

	m1 := tests.MatchQueueListDummy[0]
	mRow1 := tests.MappingMatchQueueStore(m1, 1)

	t.Run("Success Test", func(t *testing.T) {
		queues := sqlmock.NewRows(tests.MatchQueueCols).
			AddRow(mRow1...)

		mock.ExpectQuery(tests.SelectAllMatchQueueSQL).WillReturnRows(queues)

		data := repo.GetAll(ctx, db)

		t.Logf("%+v", data)

		assert.Len(t, data, 1)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Create Match Queue Test (Repository)
func TestCreateMatchQueue(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()
	repo, ctx := setupInitialMatchQueue()

	u1 := tests.UserListDummy[0]
	u1.Id = 1
	u2 := tests.UserListDummy[1]
	uRow2 := tests.MappingUserStore(u2, 2)
	argsCtInst := tests.AnyThingArgs(6)

	t.Run("Success Test", func(t *testing.T) {
		users := sqlmock.NewRows(tests.UserCols).
			AddRow(uRow2...)

		mock.ExpectQuery(regexp.QuoteMeta(tests.SelectAllTenMatchSQL)).WillReturnRows(users)
		mock.ExpectBegin()
		mock.ExpectExec(tests.InsertMatchQueueSQL).
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

// Update Match Queue Test (Repository)
func TestUpdateMatchQueue(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()
	repo, ctx := setupInitialMatchQueue()

	m1 := tests.MatchQueueListDummy[0]
	m1.Id = 1
	argsCtInst := tests.AnyThingArgs(7)

	t.Run("Success Test", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(tests.UpdateMatchQueueSQL).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		data := repo.Update(ctx, db, m1)

		t.Logf("%+v", data)

		assert.NotNil(t, data)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Delete One Match Queue Test (Repository)
func TestDeleteOneMatchQueue(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	repo, ctx := setupInitialMatchQueue()

	m1 := tests.MatchQueueListDummy[0]
	m1.Id = 1
	argsCtInst := tests.AnyThingArgs(2)

	t.Run("Success Test", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(tests.DeleteOneMatchQueueSQL)).
			WithArgs(argsCtInst...).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo.DeleteOne(ctx, db, m1)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Get One Match Queue User Test (Repository)
func TestGetOneMatchQueueUser(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	repo, ctx := setupInitialMatchQueue()

	m1 := tests.MatchQueueListDummy[0]
	mRow1 := tests.MappingMatchQueueStore(m1, 1)

	t.Run("Success Test", func(t *testing.T) {
		queue := sqlmock.NewRows(tests.MatchQueueCols).
			AddRow(mRow1...)

		mock.ExpectQuery(tests.SelectOneMatchQueueSQL).WithArgs(1).WillReturnRows(queue)

		data := repo.GetOneByUserId(ctx, db, 1)
		assert.Equal(t, 1, data.Id)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

// Get One Match Queue Current Test (Repository)
func TestGetOneMatchQueueCurrent(t *testing.T) {
	sqlDB, db, mock := tests.DbMock(t)
	defer sqlDB.Close()

	repo, ctx := setupInitialMatchQueue()

	m1 := tests.MatchQueueListDummy[0]
	mRow1 := tests.MappingMatchQueueStore(m1, 1)

	t.Run("Success Test", func(t *testing.T) {
		queue := sqlmock.NewRows(tests.MatchQueueCols).
			AddRow(mRow1...)

		mock.ExpectQuery(tests.SelectOneMatchQueueCurrentSQL).WillReturnRows(queue)

		data := repo.GetOneCurrentQueue(ctx, db, 1, helpers.GetLocalDateNow())
		assert.Equal(t, 1, data.Id)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
