package seeds

import (
	"gorm.io/gorm"
)

func DoSeed(db *gorm.DB) {
	SeedUsers(db)
	SeedPackages(db)
}
