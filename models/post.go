package models

import "time"

type Post struct {
	ID          int64
	UserID      uint64
	Title       string
	Content     string
	Thumbnail   string
	IsPublished bool
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GormPost struct {
	ID          int64 `gorm:"primary_key"`
	UserID      uint64
	Title       string
	Content     string `gorm:"type:text"`
	Thumbnail   string `gorm:"type:text"`
	IsPublished bool   `gorm:"default:false"`
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Post) TableName() string {
	return "gorm_posts"
}
