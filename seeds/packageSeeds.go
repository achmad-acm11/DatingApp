package seeds

import (
	"DatingApp/entities"
	"gorm.io/gorm"
)

func SeedPackages(db *gorm.DB) error {
	packages := []entities.ProductPackage{
		entities.ProductPackage{
			Id:          1,
			NamePackage: "Premium",
			Amount:      0,
		},
		entities.ProductPackage{
			Id:          2,
			NamePackage: "Verified",
			Amount:      0,
		},
	}

	db.Create(&packages)

	return nil
}
