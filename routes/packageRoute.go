package routes

import (
	"DatingApp/controllers"
	"DatingApp/middlewares"
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

	router.GET("packages", middlewares.JwtTokenHandler(), controller.GetAllHandler)           // ✅
	router.GET("packages/:id", middlewares.JwtTokenHandler(), controller.GetDetailHandler)    // ✅
	router.POST("packages", middlewares.JwtTokenHandler(), controller.CreateHandler)          // ✅
	router.PUT("packages/:id", middlewares.JwtTokenHandler(), controller.UpdateHandler)       // ✅
	router.DELETE("packages/:id", middlewares.JwtTokenHandler(), controller.DeleteOneHandler) // ✅

	return router
}
