package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func RunDatabase() (*gorm.DB, error) {
	dbstring := "postgres://postgres:database@localhost:5432/go-blog"
	gormDB, err := gorm.Open(postgres.Open(dbstring), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}
