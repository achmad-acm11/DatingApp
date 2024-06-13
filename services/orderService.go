package services

import (
	constants "DatingApp/constansts"
	"DatingApp/entities"
	"DatingApp/exceptions"
	"DatingApp/helpers"
	"DatingApp/repositories"
	"DatingApp/requests"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type OrderService interface {
	GetOneOrder(ctx *gin.Context, id int) entities.Order
	CreateOrder(ctx *gin.Context, request requests.CreateOrderRequest) entities.Order
	DeleteOneOrder(ctx *gin.Context, id int)
}

type orderService struct {
	repo        repositories.OrderRepository
	repoPackage repositories.PackageRepository
	repoUser    repositories.UserRepository
	db          *gorm.DB
	validate    *validator.Validate
}

func NewOrderService(db *gorm.DB, validate *validator.Validate, repo repositories.OrderRepository, repoPackage repositories.PackageRepository, repoUser repositories.UserRepository) orderService {
	return orderService{
		repo:        repo,
		repoPackage: repoPackage,
		repoUser:    repoUser,
		db:          db,
		validate:    validate,
	}
}

func (o orderService) GetOneOrder(ctx *gin.Context, id int) entities.Order {
	tx := o.db.Begin()
	defer helpers.CommitOrRollback(tx)

	order := o.repo.GetOneById(ctx, tx, id)
	o.checkOrderExistWithPanic(order)

	return order
}

func (o orderService) CreateOrder(ctx *gin.Context, request requests.CreateOrderRequest) entities.Order {
	err := o.validate.Struct(request)
	helpers.ErrorHandlerValidator(err)

	tx := o.db.Begin()
	defer helpers.CommitOrRollback(tx)

	user := o.repoUser.GetOneById(ctx, tx, request.UserId)
	o.checkUserExistWithPanic(user)

	productPackage := o.repoPackage.GetOneById(ctx, tx, request.PackageId)
	o.checkPackageExistWithPanic(productPackage)

	if request.Amount != productPackage.Amount {
		panic(exceptions.NewBadRequestError(errors.New(constants.OrderInsufficient).Error()))
	}

	orderCheck := o.repo.GetOrderUser(ctx, tx, productPackage.Id, user.Id)
	if orderCheck.Id != 0 {
		panic(exceptions.NewConflictError(errors.New(constants.PackagePurchased).Error()))
	}

	order := o.repo.Create(ctx, tx, entities.Order{
		UserId:      user.Id,
		PackageId:   productPackage.Id,
		NamePackage: productPackage.NamePackage,
		Amount:      productPackage.Amount,
	})

	if user.IsPremium == false {
		user.IsPremium = true
		o.repoUser.Update(ctx, tx, user)
	}

	return order
}

func (o orderService) DeleteOneOrder(ctx *gin.Context, id int) {
	tx := o.db.Begin()
	defer helpers.CommitOrRollback(tx)

	order := o.repo.GetOneById(ctx, tx, id)
	o.checkOrderExistWithPanic(order)

	o.repo.DeleteOne(ctx, tx, order)
}

func (o orderService) checkOrderExistWithPanic(order entities.Order) {
	if order.Id == 0 {
		panic(exceptions.NewNotFoundError(errors.New(constants.OrderNotFound).Error()))
	}
}

func (o orderService) checkPackageExistWithPanic(productPackage entities.ProductPackage) {
	if productPackage.Id == 0 {
		panic(exceptions.NewNotFoundError(errors.New(constants.PackageNotFound).Error()))
	}
}

func (o orderService) checkUserExistWithPanic(user entities.User) {
	if user.Id == 0 {
		panic(exceptions.NewNotFoundError(errors.New(constants.UserNotFound).Error()))
	}
}
