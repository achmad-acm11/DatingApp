package repositories

import (
	"DatingApp/entities"
	"DatingApp/helpers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type UserFindRepository interface {
	GetAll(ctx *gin.Context, db *gorm.DB) []entities.UserFind
	GetOneByUserId(ctx *gin.Context, db *gorm.DB, userId int) entities.UserFind
	GetOneCurrentFind(ctx *gin.Context, db *gorm.DB, userId int, date time.Time) entities.UserFind
	Create(ctx *gin.Context, db *gorm.DB, user entities.User) entities.UserFind
	Update(ctx *gin.Context, db *gorm.DB, userFind entities.UserFind) entities.UserFind
	Reset(ctx *gin.Context, db *gorm.DB, userFind entities.UserFind, user entities.User) entities.UserFind
	DeleteOne(ctx *gin.Context, db *gorm.DB, userFind entities.UserFind)
}

type userFindRepository struct {
	repoUser UserRepository
}

func NewUserFindRepository(repoUser UserRepository) userFindRepository {
	return userFindRepository{
		repoUser: repoUser,
	}
}

func (f userFindRepository) GetAll(ctx *gin.Context, db *gorm.DB) []entities.UserFind {
	userFinds := []entities.UserFind{}
	tx := db.WithContext(ctx)

	err := tx.Limit(500).Find(&userFinds).Error
	helpers.ErrorHandler(err)

	return userFinds
}

func (f userFindRepository) GetOneByUserId(ctx *gin.Context, db *gorm.DB, userId int) entities.UserFind {
	userFind := entities.UserFind{}

	err := db.WithContext(ctx).Where("user_id = ?", userId).Find(&userFind).Error
	helpers.ErrorHandler(err)

	return userFind
}

func (f userFindRepository) GetOneCurrentFind(ctx *gin.Context, db *gorm.DB, userId int, date time.Time) entities.UserFind {
	userFind := entities.UserFind{}

	err := db.WithContext(ctx).Where("user_id = ?", userId).Where("date = ?", date).Find(&userFind).Error
	helpers.ErrorHandler(err)

	return userFind
}

func (f userFindRepository) Create(ctx *gin.Context, db *gorm.DB, user entities.User) entities.UserFind {
	timeNow := helpers.GetLocalDateNow()

	userQueue := f.getUserQueue(ctx, db, user.Id, user.Gender)

	userFind := entities.UserFind{
		UserId:       user.Id,
		PassCount:    0,
		LikeCount:    0,
		CurrentState: 1,
		UserQueue:    strings.Join(userQueue, "|"),
		Date:         timeNow,
	}

	err := db.WithContext(ctx).Create(&userFind).Error
	helpers.ErrorHandler(err)

	return userFind
}

func (f userFindRepository) Update(ctx *gin.Context, db *gorm.DB, userFind entities.UserFind) entities.UserFind {
	err := db.WithContext(ctx).Save(&userFind).Error
	helpers.ErrorHandler(err)

	return userFind
}

func (f userFindRepository) Reset(ctx *gin.Context, db *gorm.DB, userFind entities.UserFind, user entities.User) entities.UserFind {
	timeNow := helpers.GetLocalDateNow()

	userQueue := f.getUserQueue(ctx, db, user.Id, user.Gender)

	userFind.CurrentState = 1
	userFind.PassCount = 0
	userFind.LikeCount = 0
	userFind.UserQueue = strings.Join(userQueue, "|")
	userFind.Date = timeNow

	userFind = f.Update(ctx, db, userFind)

	return userFind
}

func (f userFindRepository) DeleteOne(ctx *gin.Context, db *gorm.DB, userFind entities.UserFind) {
	err := db.WithContext(ctx).Delete(&userFind).Error
	helpers.ErrorHandler(err)
}

func (f userFindRepository) getUserQueue(ctx *gin.Context, db *gorm.DB, userId int, gender string) []string {
	usersCouple := f.repoUser.GetAllTenCouple(ctx, db, userId, gender)
	var userQueue []string
	for _, userCouple := range usersCouple {
		userQueue = append(userQueue, strconv.Itoa(userCouple.Id))
	}
	return userQueue
}
