package routes

import (
	"DatingApp/controllers"
	"DatingApp/repositories"
	"DatingApp/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func UserRoute(db *gorm.DB, validate *validator.Validate, router *gin.RouterGroup) *gin.RouterGroup {
	repository := repositories.NewUserRepository()
	repositoryFind := repositories.NewUserFindRepository(repository)
	service := services.NewUserService(
		db,
		validate,
		repository,
		repositoryFind,
	)
	controller := controllers.NewUserController(service)

	router.GET("users", controller.GetAllHandler)                         // ✅
	router.GET("users/:userId", controller.GetDetailHandler)              // ✅
	router.GET("users/:userId/couple", controller.GetCoupleHandler)       // ✅
	router.GET("users/:userId/like-couple", controller.LikeCoupleHandler) // ✅
	router.GET("users/:userId/pass-couple", controller.PassCoupleHandler) // ✅
	router.POST("signup", controller.SignUpHandler)                       // ✅
	router.POST("login", controller.LoginHandler)                         // ✅
	router.DELETE("users/:userId", controller.DeleteOneHandler)           // ✅
	return router
}
