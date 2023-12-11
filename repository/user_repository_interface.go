package repository

import (
	"context"
	"errors"
	"postgresql-blog/models"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExist     = errors.New("row does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

// Repository provides access to the website storage.
type UserRepository interface {
	MigrateUser(ctx context.Context) error
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	AllUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsernameAndPassword(ctx context.Context, username string, password string) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, updated models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id int64) error
}
