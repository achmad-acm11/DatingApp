package entities

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model  `json:"-"`
	Id          int            `gorm:"column:id;primaryKey;autoIncrement;type:int(11)" json:"id"`
	UserId      int            `gorm:"column:user_id;type:int(11);not null" json:"user_id"`
	PackageId   int            `gorm:"column:package_id;type:int(11);not null" json:"package_id"`
	NamePackage string         `gorm:"column:name_package;type:varchar(255);not null" json:"name_package"`
	Amount      int            `gorm:"column:amount;type:int(11);default:0" json:"amount"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamp null;default:current_timestamp();->" json:"-"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamp null;default:null on update current_timestamp();->" json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp null;default:null;->" json:"-"`
}

func (o Order) TableName() string {
	return "orders"
}
