package repository

import (
	"context"
	"errors"

	"postgresql-blog/models"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func (repo *PostgreSQLGORMRepository) MigratePost(ctx context.Context) error {
	err := repo.db.WithContext(ctx).AutoMigrate(&models.GormPost{})
	if err != nil {
		return err
	}
	return nil
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &PostgreSQLGORMRepository{db}
}

func (repo *PostgreSQLGORMRepository) CreatePost(ctx context.Context, post models.Post) (*models.Post, error) {
	gormPost := models.Post{
		UserID:      post.UserID,
		Title:       post.Title,
		Content:     post.Content,
		Thumbnail:   post.Thumbnail,
		IsPublished: post.IsPublished,
		PublishedAt: post.PublishedAt,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}

	if err := repo.db.WithContext(ctx).Create(&gormPost).Error; err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	result := models.Post(gormPost)
	return &result, nil
}

func (repo *PostgreSQLGORMRepository) AllPosts(ctx context.Context) ([]models.Post, error) {
	var allPosts []models.GormPost
	if err := repo.db.WithContext(ctx).Find(&allPosts).Error; err != nil {
		return nil, err
	}

	var result []models.Post
	for _, posts := range allPosts {
		result = append(result, models.Post(posts))
	}

	return result, nil
}

func (repo *PostgreSQLGORMRepository) GetPostByID(ctx context.Context, id int64) (*models.Post, error) {
	var gormPost models.GormPost
	if err := repo.db.WithContext(ctx).Where("id = ?", id).First(&gormPost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	result := models.Post(gormPost)
	return &result, nil
}

func (repo *PostgreSQLGORMRepository) GetPostByTitle(ctx context.Context, title string) (*models.Post, error) {
	var gormPost models.GormPost
	if err := repo.db.WithContext(ctx).Where("title = ?", title).First(&gormPost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	result := models.Post(gormPost)
	return &result, nil
}

func (repo *PostgreSQLGORMRepository) GetPostByUserID(ctx context.Context, userid int64) ([]models.Post, error) {
	var gormPost []models.GormPost
	if err := repo.db.WithContext(ctx).Where("user_id = ?", userid).Find(&gormPost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	var result []models.Post
	for _, posts := range gormPost {
		result = append(result, models.Post(posts))
	}

	return result, nil
}

func (repo *PostgreSQLGORMRepository) UpdatePost(ctx context.Context, id int64, updated models.Post) (*models.Post, error) {
	gormPost := models.Post(updated)
	updateRes := repo.db.WithContext(ctx).Where("id = ?", id).Save(&gormPost)
	if err := updateRes.Error; err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	rowsAffected := updateRes.RowsAffected
	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}
	return &updated, nil
}

func (repo *PostgreSQLGORMRepository) DeletePost(ctx context.Context, id int64) error {
	res := repo.db.WithContext(ctx).Delete(&models.GormPost{}, id)
	if err := res.Error; err != nil {
		return err
	}

	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return nil
}
