package repository

import (
	"context"
	"postgresql-blog/models"
)

// Repository provides access to the website storage.
type CommentRepository interface {
	MigrateComment(ctx context.Context) error
	CreateComment(ctx context.Context, comment models.Comment) (*models.Comment, error)
	AllComments(ctx context.Context) ([]models.Comment, error)
	GetCommentByID(ctx context.Context, id int64) (*models.Comment, error)
	GetCommentByUserID(ctx context.Context, userid int64) ([]models.Comment, error)
	GetCommentByPostID(ctx context.Context, postid int64) ([]models.Comment, error)
	GetCommentByUserIDPostID(ctx context.Context, userid int64, postid int64) (*models.Comment, error)
	UpdateComment(ctx context.Context, id int64, updated models.Comment) (*models.Comment, error)
	DeleteComment(ctx context.Context, id int64) error
}
