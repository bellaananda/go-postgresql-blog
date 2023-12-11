package models

import "time"

type Comment struct {
	ID          int64
	UserID      uint64
	PostID      uint64
	Content     string
	IsPublished bool
	PublishedAt time.Time
	DeletedAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GormComment struct {
	ID          int64 `gorm:"primary_key"`
	UserID      uint64
	PostID      uint64
	Content     string `gorm:"type:text"`
	IsPublished bool   `gorm:"default:false"`
	PublishedAt time.Time
	DeletedAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Comment) TableName() string {
	return "gorm_comments"
}
