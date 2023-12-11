package database

import (
	"context"

	"postgresql-blog/models"

	"gorm.io/gorm"
)

type PostgreSQLGORMRepository struct {
	db *gorm.DB
}

func NewPostgreSQLGORMRepository(db *gorm.DB) *PostgreSQLGORMRepository {
	return &PostgreSQLGORMRepository{
		db: db,
	}
}

// migrate database
func (r *PostgreSQLGORMRepository) Migrate(ctx context.Context) error {
	// table user
	err := r.db.WithContext(ctx).AutoMigrate(&models.GormUser{})
	if err != nil {
		return err
	}

	// table post
	err = r.db.WithContext(ctx).AutoMigrate(&models.GormPost{})
	if err != nil {
		return err
	}

	// table comment
	err = r.db.WithContext(ctx).AutoMigrate(&models.GormComment{})
	if err != nil {
		return err
	}

	return nil
}
