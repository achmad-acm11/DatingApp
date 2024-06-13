package migrations

import (
	"DatingApp/entities"
	"gorm.io/gorm"
)

func DoMigration(db *gorm.DB) {
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&entities.User{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&entities.UserFind{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&entities.ProductPackage{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&entities.Order{})
}
