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
	stdLog   *helpers.StandartLog
}

func NewPackageService(db *gorm.DB, validate *validator.Validate, repo repositories.PackageRepository) packageService {
	return packageService{
		repo:     repo,
		db:       db,
		validate: validate,
		stdLog:   helpers.NewStandardLog(constants.Package, constants.Service),
	}
}

func (p packageService) GetAllPackage(ctx *gin.Context) []entities.ProductPackage {
	p.stdLog.NameFunc = "GetAllPackage"
	p.stdLog.StartFunction(nil)

	tx := p.db.Begin()
	defer helpers.CommitOrRollback(tx)

	packages := p.repo.GetAll(ctx, tx)

	p.stdLog.NameFunc = "GetAllPackage"
	p.stdLog.EndFunction(nil)

	return packages
}

func (p packageService) GetOnePackage(ctx *gin.Context, id int) entities.ProductPackage {
	p.stdLog.NameFunc = "GetOnePackage"
	p.stdLog.StartFunction(id)

	tx := p.db.Begin()
	defer helpers.CommitOrRollback(tx)

	productPackage := p.repo.GetOneById(ctx, tx, id)
	p.checkPackageExistWithPanic(productPackage)

	p.stdLog.NameFunc = "GetOnePackage"
	p.stdLog.EndFunction(productPackage)

	return productPackage
}

func (p packageService) CreatePackage(ctx *gin.Context, request requests.CreatePackageRequest) entities.ProductPackage {
	p.stdLog.NameFunc = "CreatePackage"
	p.stdLog.StartFunction(request)

	err := p.validate.Struct(request)
	helpers.ErrorHandlerValidator(err)

	tx := p.db.Begin()
	defer helpers.CommitOrRollback(tx)

	productPackage := p.repo.Create(ctx, tx, entities.ProductPackage{
		NamePackage: request.NamePackage,
		Amount:      request.Amount,
	})

	p.stdLog.NameFunc = "CreatePackage"
	p.stdLog.EndFunction(productPackage)

	return productPackage
}

func (p packageService) UpdatePackage(ctx *gin.Context, request requests.UpdatePackageRequest, id int) entities.ProductPackage {
	p.stdLog.NameFunc = "UpdatePackage"
	p.stdLog.StartFunction(request)

	err := p.validate.Struct(request)
	helpers.ErrorHandlerValidator(err)

	tx := p.db.Begin()
	defer helpers.CommitOrRollback(tx)

	productPackage := p.repo.GetOneById(ctx, tx, id)
	p.checkPackageExistWithPanic(productPackage)

	productPackage.NamePackage = request.NamePackage
	productPackage.Amount = request.Amount

	productPackageNew := p.repo.Update(ctx, tx, productPackage)

	p.stdLog.NameFunc = "UpdatePackage"
	p.stdLog.EndFunction(productPackage)

	return productPackageNew
}

func (p packageService) DeleteOnePackage(ctx *gin.Context, id int) {
	p.stdLog.NameFunc = "DeleteOnePackage"
	p.stdLog.StartFunction(id)

	tx := p.db.Begin()
	defer helpers.CommitOrRollback(tx)

	productPackage := p.repo.GetOneById(ctx, tx, id)
	p.checkPackageExistWithPanic(productPackage)

	p.repo.DeleteOne(ctx, tx, productPackage)

	p.stdLog.NameFunc = "DeleteOnePackage"
	p.stdLog.EndFunction(nil)
}

func (p packageService) checkPackageExistWithPanic(productPackage entities.ProductPackage) {
	if productPackage.Id == 0 {
		panic(exceptions.NewNotFoundError(errors.New(constants.PackageNotFound).Error()))
	}
}
