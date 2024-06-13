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

func UserRoute(db *gorm.DB, validate *validator.Validate, router *gin.RouterGroup) *gin.RouterGroup {
	repository := repositories.NewUserRepository()
	repositoryQueue := repositories.NewMatchQueueRepository(repository)
	service := services.NewUserService(
		db,
		validate,
		repository,
		repositoryQueue,
	)
	controller := controllers.NewUserController(service)

	router.GET("users", middlewares.JwtTokenHandler(), controller.GetAllHandler)                     // ✅
	router.GET("users/my-profile", middlewares.JwtTokenHandler(), controller.GetMyProfileHandler)    // ✅
	router.GET("users/profile/:userId", middlewares.JwtTokenHandler(), controller.GetProfileHandler) // ✅
	router.GET("users/matches", middlewares.JwtTokenHandler(), controller.GetMatchesHandler)         // ✅
	router.GET("users/like-match", middlewares.JwtTokenHandler(), controller.LikeMatchHandler)       // ✅
	router.GET("users/pass-match", middlewares.JwtTokenHandler(), controller.PassMatchHandler)       // ✅
	router.POST("signup", controller.SignUpHandler)                                                  // ✅
	router.POST("login", controller.LoginHandler)                                                    // ✅
	router.DELETE("users/:userId", middlewares.JwtTokenHandler(), controller.DeleteOneHandler)       // ✅
	return router
}
