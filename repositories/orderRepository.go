package repositories

import (
	"DatingApp/entities"
	"DatingApp/helpers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderRepository interface {
	GetAll(ctx *gin.Context, db *gorm.DB) []entities.Order
	GetOneById(ctx *gin.Context, db *gorm.DB, id int) entities.Order
	GetOrderUser(ctx *gin.Context, db *gorm.DB, packageId int, userId int) entities.Order
	Create(ctx *gin.Context, db *gorm.DB, order entities.Order) entities.Order
	Update(ctx *gin.Context, db *gorm.DB, order entities.Order) entities.Order
	DeleteOne(ctx *gin.Context, db *gorm.DB, order entities.Order)
}

type orderRepository struct {
}

func NewOrderRepository() orderRepository {
	return orderRepository{}
}

func (o orderRepository) GetAll(ctx *gin.Context, db *gorm.DB) []entities.Order {
	orders := []entities.Order{}
	tx := db.WithContext(ctx)

	err := tx.Limit(500).Find(&orders).Error
	helpers.ErrorHandler(err)

	return orders
}

func (o orderRepository) GetOneById(ctx *gin.Context, db *gorm.DB, id int) entities.Order {
	order := entities.Order{}

	err := db.WithContext(ctx).Where("id = ?", id).Find(&order).Error
	helpers.ErrorHandler(err)

	return order
}

func (o orderRepository) GetOrderUser(ctx *gin.Context, db *gorm.DB, packageId int, userId int) entities.Order {
	order := entities.Order{}

	err := db.WithContext(ctx).Where("package_id = ?", packageId).Where("user_id = ?", userId).Find(&order).Error
	helpers.ErrorHandler(err)

	return order
}

func (o orderRepository) Create(ctx *gin.Context, db *gorm.DB, order entities.Order) entities.Order {
	err := db.WithContext(ctx).Create(&order).Error
	helpers.ErrorHandler(err)

	return order
}

func (o orderRepository) Update(ctx *gin.Context, db *gorm.DB, order entities.Order) entities.Order {
	err := db.WithContext(ctx).Model(&order).Updates(order).Error
	helpers.ErrorHandler(err)

	return order
}

func (o orderRepository) DeleteOne(ctx *gin.Context, db *gorm.DB, order entities.Order) {
	err := db.WithContext(ctx).Delete(&order).Error
	helpers.ErrorHandler(err)
}
