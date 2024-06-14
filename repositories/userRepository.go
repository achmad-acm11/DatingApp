package repositories

import (
	"DatingApp/entities"
	"DatingApp/helpers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

type UserRepository interface {
	GetAll(ctx *gin.Context, db *gorm.DB) []entities.User
	GetAllTenMatch(ctx *gin.Context, db *gorm.DB, userId int, gender string, userIds []string) []entities.User
	GetOneById(ctx *gin.Context, db *gorm.DB, id int) entities.User
	GetOneCouple(ctx *gin.Context, db *gorm.DB, userId int, gender string) entities.User
	Create(ctx *gin.Context, db *gorm.DB, user entities.User) entities.User
	GetOneByEmail(ctx *gin.Context, db *gorm.DB, email string) entities.User
	Update(ctx *gin.Context, db *gorm.DB, user entities.User) entities.User
	DeleteOne(ctx *gin.Context, db *gorm.DB, user entities.User)
}

type userRepository struct {
}

func NewUserRepository() userRepository {
	return userRepository{}
}

func (u userRepository) GetAll(ctx *gin.Context, db *gorm.DB) []entities.User {
	users := []entities.User{}
	tx := db.WithContext(ctx)

	err := tx.Limit(500).Find(&users).Error
	helpers.ErrorHandler(err)

	return users
}

func (u userRepository) Create(ctx *gin.Context, db *gorm.DB, user entities.User) entities.User {
	err := db.WithContext(ctx).Create(&user).Error
	helpers.ErrorHandler(err)

	return user
}

func (u userRepository) GetOneByEmail(ctx *gin.Context, db *gorm.DB, email string) entities.User {

	user := entities.User{}

	err := db.WithContext(ctx).Where("email = ?", email).Find(&user).Error
	helpers.ErrorHandler(err)

	return user
}

func (u userRepository) GetOneCouple(ctx *gin.Context, db *gorm.DB, userId int, gender string) entities.User {
	user := entities.User{}

	err := db.WithContext(ctx).Where("id != ?", userId).Where("gender != ?", gender).Find(&user).Error
	helpers.ErrorHandler(err)

	return user
}

func (u userRepository) GetAllTenMatch(ctx *gin.Context, db *gorm.DB, userId int, gender string, userIds []string) []entities.User {
	users := []entities.User{}

	tx := db.WithContext(ctx).Where("id != ?", userId).Where("gender != ?", gender)
	if len(userIds) > 0 {
		var userIdsInt []int
		for _, val := range userIds {
			userId, _ = strconv.Atoi(val)
			userIdsInt = append(userIdsInt, userId)
		}
		tx = tx.Where("id NOT IN ?", userIdsInt)
	}
	err := tx.Limit(10).Find(&users).Error
	helpers.ErrorHandler(err)

	return users
}

func (u userRepository) GetOneById(ctx *gin.Context, db *gorm.DB, id int) entities.User {
	user := entities.User{}

	err := db.WithContext(ctx).Where("id = ?", id).Find(&user).Error
	helpers.ErrorHandler(err)

	return user
}

func (u userRepository) Update(ctx *gin.Context, db *gorm.DB, user entities.User) entities.User {
	err := db.WithContext(ctx).Model(&user).Updates(user).Error
	helpers.ErrorHandler(err)

	return user
}

func (u userRepository) DeleteOne(ctx *gin.Context, db *gorm.DB, user entities.User) {
	err := db.WithContext(ctx).Delete(&user).Error
	helpers.ErrorHandler(err)
}
