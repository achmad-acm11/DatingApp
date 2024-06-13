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

func OrderRoute(db *gorm.DB, validate *validator.Validate, router *gin.RouterGroup) *gin.RouterGroup {
	repository := repositories.NewOrderRepository()
	repositoryUser := repositories.NewUserRepository()
	repositoryPackage := repositories.NewPackageRepository()
	service := services.NewOrderService(
		db,
		validate,
		repository,
		repositoryPackage,
		repositoryUser,
	)
	controller := controllers.NewOrderController(service)

	router.GET("orders/:id", middlewares.JwtTokenHandler(), controller.GetDetailHandler)    // ✅
	router.POST("orders", middlewares.JwtTokenHandler(), controller.CreateHandler)          // ✅
	router.DELETE("orders/:id", middlewares.JwtTokenHandler(), controller.DeleteOneHandler) // ✅

	return router
}
