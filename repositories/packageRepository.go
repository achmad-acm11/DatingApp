package repositories

import (
	"DatingApp/entities"
	"DatingApp/helpers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PackageRepository interface {
	GetAll(ctx *gin.Context, db *gorm.DB) []entities.ProductPackage
	GetOneById(ctx *gin.Context, db *gorm.DB, id int) entities.ProductPackage
	Create(ctx *gin.Context, db *gorm.DB, productPackage entities.ProductPackage) entities.ProductPackage
	Update(ctx *gin.Context, db *gorm.DB, productPackage entities.ProductPackage) entities.ProductPackage
	DeleteOne(ctx *gin.Context, db *gorm.DB, productPackage entities.ProductPackage)
}

type packageRepository struct {
}

func NewPackageRepository() packageRepository {
	return packageRepository{}
}

func (p packageRepository) GetAll(ctx *gin.Context, db *gorm.DB) []entities.ProductPackage {
	packages := []entities.ProductPackage{}
	tx := db.WithContext(ctx)

	err := tx.Limit(500).Find(&packages).Error
	helpers.ErrorHandler(err)

	return packages
}

func (p packageRepository) GetOneById(ctx *gin.Context, db *gorm.DB, id int) entities.ProductPackage {
	productPackage := entities.ProductPackage{}

	err := db.WithContext(ctx).Where("id = ?", id).Find(&productPackage).Error
	helpers.ErrorHandler(err)

	return productPackage
}

func (p packageRepository) Create(ctx *gin.Context, db *gorm.DB, productPackage entities.ProductPackage) entities.ProductPackage {
	err := db.WithContext(ctx).Create(&productPackage).Error
	helpers.ErrorHandler(err)

	return productPackage
}

func (p packageRepository) Update(ctx *gin.Context, db *gorm.DB, productPackage entities.ProductPackage) entities.ProductPackage {
	err := db.WithContext(ctx).Model(&productPackage).Updates(productPackage).Error
	helpers.ErrorHandler(err)

	return productPackage
}

func (p packageRepository) DeleteOne(ctx *gin.Context, db *gorm.DB, productPackage entities.ProductPackage) {
	err := db.WithContext(ctx).Delete(&productPackage).Error
	helpers.ErrorHandler(err)
}
