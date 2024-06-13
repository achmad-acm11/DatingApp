package routes

import (
	"DatingApp/controllers"
	"DatingApp/repositories"
	"DatingApp/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func PackageRoute(db *gorm.DB, validate *validator.Validate, router *gin.RouterGroup) *gin.RouterGroup {
	repository := repositories.NewPackageRepository()
	service := services.NewPackageService(
		db,
		validate,
		repository,
	)
	controller := controllers.NewPackageController(service)

	router.GET("packages", controller.GetAllHandler)           // ✅
	router.GET("packages/:id", controller.GetDetailHandler)    // ✅
	router.POST("packages", controller.CreateHandler)          // ✅
	router.PUT("packages/:id", controller.UpdateHandler)       // ✅
	router.DELETE("packages/:id", controller.DeleteOneHandler) // ✅

	return router
}
