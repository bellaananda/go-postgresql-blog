package repository

import (
	"context"
	"errors"

	"postgresql-blog/models"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func (repo *PostgreSQLGORMRepository) MigrateComment(ctx context.Context) error {
	err := repo.db.WithContext(ctx).AutoMigrate(&models.GormComment{})
	if err != nil {
		return err
	}
	return nil
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &PostgreSQLGORMRepository{db}
}

func (repo *PostgreSQLGORMRepository) CreateComment(ctx context.Context, comment models.Comment) (*models.Comment, error) {
	gormComment := models.Comment{
		UserID:      comment.UserID,
		PostID:      comment.PostID,
		Content:     comment.Content,
		IsPublished: comment.IsPublished,
		PublishedAt: comment.PublishedAt,
		CreatedAt:   comment.CreatedAt,
		UpdatedAt:   comment.UpdatedAt,
	}

	if err := repo.db.WithContext(ctx).Create(&gormComment).Error; err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	result := models.Comment(gormComment)
	return &result, nil
}

func (repo *PostgreSQLGORMRepository) AllComments(ctx context.Context) ([]models.Comment, error) {
	var allComments []models.GormComment
	if err := repo.db.WithContext(ctx).Find(&allComments).Error; err != nil {
		return nil, err
	}

	var result []models.Comment
	for _, comments := range allComments {
		result = append(result, models.Comment(comments))
	}

	return result, nil
}

func (repo *PostgreSQLGORMRepository) GetCommentByID(ctx context.Context, id int64) (*models.Comment, error) {
	var gormComment models.GormComment
	if err := repo.db.WithContext(ctx).Where("id = ?", id).First(&gormComment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	result := models.Comment(gormComment)
	return &result, nil
}

func (repo *PostgreSQLGORMRepository) GetCommentByUserID(ctx context.Context, userid int64) ([]models.Comment, error) {
	var gormComment []models.GormComment
	if err := repo.db.WithContext(ctx).Where("user_id = ?", userid).Find(&gormComment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	var result []models.Comment
	for _, comments := range gormComment {
		result = append(result, models.Comment(comments))
	}

	return result, nil
}

func (repo *PostgreSQLGORMRepository) GetCommentByPostID(ctx context.Context, postid int64) ([]models.Comment, error) {
	var gormComment []models.GormComment
	if err := repo.db.WithContext(ctx).Where("post_id = ?", postid).Find(&gormComment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	var result []models.Comment
	for _, comments := range gormComment {
		result = append(result, models.Comment(comments))
	}

	return result, nil
}

func (repo *PostgreSQLGORMRepository) GetCommentByUserIDPostID(ctx context.Context, userid int64, postid int64) (*models.Comment, error) {
	var gormComment models.GormComment
	if err := repo.db.WithContext(ctx).Where("user_id = ? AND post_id = ?", userid, postid).First(&gormComment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	result := models.Comment(gormComment)
	return &result, nil
}

func (repo *PostgreSQLGORMRepository) UpdateComment(ctx context.Context, id int64, updated models.Comment) (*models.Comment, error) {
	gormComment := models.Comment(updated)
	updateRes := repo.db.WithContext(ctx).Where("id = ?", id).Save(&gormComment)
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

func (repo *PostgreSQLGORMRepository) DeleteComment(ctx context.Context, id int64) error {
	res := repo.db.WithContext(ctx).Delete(&models.GormComment{}, id)
	if err := res.Error; err != nil {
		return err
	}

	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return nil
}
