package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model  `json:"-"`
	Id          int            `gorm:"column:id;primaryKey;autoIncrement;type:int(11)" json:"id"`
	Name        string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Email       string         `gorm:"column:email;type:varchar(255);uniqueIndex;not null" json:"email"`
	Gender      string         `gorm:"column:gender;type:enum('male', 'female');default:male;not null" json:"gender"`
	PhoneNumber string         `gorm:"column:phone_number;type:varchar(255);not null" json:"phone_number"`
	Password    string         `gorm:"column:password;type:varchar(255);not null" json:"-"`
	IsPremium   bool           `gorm:"column:is_premium;type:tinyint(1);default:0" json:"is_premium"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamp null;default:current_timestamp();->" json:"-"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamp null;default:null on update current_timestamp();->" json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp null;default:null;->" json:"-"`
}

func (u User) TableName() string {
	return "users"
}
