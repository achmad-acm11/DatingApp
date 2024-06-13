package entities

import (
	"gorm.io/gorm"
	"time"
)

type MatchQueue struct {
	gorm.Model   `json:"-"`
	Id           int            `gorm:"column:id;primaryKey;autoIncrement;type:int(11)" json:"id"`
	UserId       int            `gorm:"column:user_id;type:int(11);not null" json:"user_id"`
	PassCount    int            `gorm:"column:pass_count;type:int(11);default:0" json:"pass_count"`
	LikeCount    int            `gorm:"column:like_count;type:int(11);default:0" json:"like_count"`
	CurrentState int            `gorm:"column:current_state;type:int(11);default:0" json:"current_state"`
	UserQueue    string         `gorm:"column:user_queue;type:text" json:"user_queue"`
	Date         time.Time      `gorm:"column:date;type:datetime;default:null" json:"date"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:timestamp null;default:current_timestamp();->" json:"-"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:timestamp null;default:null on update current_timestamp();->" json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp null;default:null;->" json:"-"`
}

func (u MatchQueue) TableName() string {
	return "match_queues"
}
