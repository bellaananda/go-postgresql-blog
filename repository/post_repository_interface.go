package repository

import (
	"context"
	"postgresql-blog/models"
)

// Repository provides access to the website storage.
type PostRepository interface {
	MigratePost(ctx context.Context) error
	CreatePost(ctx context.Context, post models.Post) (*models.Post, error)
	AllPosts(ctx context.Context) ([]models.Post, error)
	GetPostByID(ctx context.Context, id int64) (*models.Post, error)
	GetPostByTitle(ctx context.Context, title string) (*models.Post, error)
	GetPostByUserID(ctx context.Context, userid int64) ([]models.Post, error)
	UpdatePost(ctx context.Context, id int64, updated models.Post) (*models.Post, error)
	DeletePost(ctx context.Context, id int64) error
}
