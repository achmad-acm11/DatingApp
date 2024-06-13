package configs

import (
	"DatingApp/helpers"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func ConfigDB() *gorm.DB {
	if os.Getenv("APP_ENV") == "" {
		errEnv := godotenv.Load(".env")

		helpers.ErrorHandler(errEnv)
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	config := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, db_name)

	db, err := gorm.Open(mysql.Open(config), &gorm.Config{})

	helpers.ErrorHandler(err)

	return db
}
