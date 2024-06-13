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

type PackageService interface {
	GetAllPackage(ctx *gin.Context) []entities.ProductPackage
	GetOnePackage(ctx *gin.Context, id int) entities.ProductPackage
	CreatePackage(ctx *gin.Context, request requests.CreatePackageRequest) entities.ProductPackage
	UpdatePackage(ctx *gin.Context, request requests.UpdatePackageRequest, id int) entities.ProductPackage
	DeleteOnePackage(ctx *gin.Context, id int)
}

type packageService struct {
	repo     repositories.PackageRepository
	db       *gorm.DB
	validate *validator.Validate
}

func NewPackageService(db *gorm.DB, validate *validator.Validate, repo repositories.PackageRepository) packageService {
	return packageService{
		repo:     repo,
		db:       db,
		validate: validate,
	}
}

func (p packageService) GetAllPackage(ctx *gin.Context) []entities.ProductPackage {
	tx := p.db.Begin()
	defer helpers.CommitOrRollback(tx)

	packages := p.repo.GetAll(ctx, tx)

	return packages
}

func (p packageService) GetOnePackage(ctx *gin.Context, id int) entities.ProductPackage {
	tx := p.db.Begin()
	defer helpers.CommitOrRollback(tx)

	productPackage := p.repo.GetOneById(ctx, tx, id)
	p.checkPackageExistWithPanic(productPackage)

	return productPackage
}

func (p packageService) CreatePackage(ctx *gin.Context, request requests.CreatePackageRequest) entities.ProductPackage {
	err := p.validate.Struct(request)
	helpers.ErrorHandlerValidator(err)

	tx := p.db.Begin()
	defer helpers.CommitOrRollback(tx)

	productPackage := p.repo.Create(ctx, tx, entities.ProductPackage{
		NamePackage: request.NamePackage,
		Amount:      request.Amount,
	})

	return productPackage
}

func (p packageService) UpdatePackage(ctx *gin.Context, request requests.UpdatePackageRequest, id int) entities.ProductPackage {
	err := p.validate.Struct(request)
	helpers.ErrorHandlerValidator(err)

	tx := p.db.Begin()
	defer helpers.CommitOrRollback(tx)

	productPackage := p.repo.GetOneById(ctx, tx, id)
	p.checkPackageExistWithPanic(productPackage)

	productPackageNew := p.repo.Create(ctx, tx, entities.ProductPackage{
		NamePackage: request.NamePackage,
		Amount:      request.Amount,
	})

	return productPackageNew
}

func (p packageService) DeleteOnePackage(ctx *gin.Context, id int) {
	tx := p.db.Begin()
	defer helpers.CommitOrRollback(tx)

	productPackage := p.repo.GetOneById(ctx, tx, id)
	p.checkPackageExistWithPanic(productPackage)

	p.repo.DeleteOne(ctx, tx, productPackage)
}

func (p packageService) checkPackageExistWithPanic(productPackage entities.ProductPackage) {
	if productPackage.Id == 0 {
		panic(exceptions.NewNotFoundError(errors.New(constants.PackageNotFound).Error()))
	}
}
