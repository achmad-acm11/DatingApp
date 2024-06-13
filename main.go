package main

import (
	"DatingApp/configs"
	"DatingApp/helpers"
	"DatingApp/middlewares"
	"DatingApp/migrations"
	"DatingApp/routes"
	"github.com/gin-gonic/gin"
	validator2 "github.com/go-playground/validator/v10"
	cors "github.com/itsjamie/gin-cors"
	"github.com/joho/godotenv"
	"os"
	"time"
)

func main() {
	if os.Getenv("APP_ENV") == "" {
		errEnv := godotenv.Load(".env")
		helpers.ErrorHandler(errEnv)
	}
	app_port := os.Getenv("APP_PORT")
	validator := validator2.New()
	db := configs.ConfigDB()
	migrations.DoMigration(db)

	router := gin.Default()
	router.Use(middlewares.ErrorHandler())

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE, OPTIONS",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		ValidateHeaders: false,
	}))

	api := router.Group("api/")
	api = routes.UserRoute(db, validator, api)
	api = routes.PackageRoute(db, validator, api)
	routes.OrderRoute(db, validator, api)

	router.Run(":" + app_port)
}
