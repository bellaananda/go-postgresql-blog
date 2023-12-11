package repository

import (
	"context"
	"errors"

	"postgresql-blog/models"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type PostgreSQLGORMRepository struct {
	db *gorm.DB
}

func (repo *PostgreSQLGORMRepository) MigrateUser(ctx context.Context) error {
	err := repo.db.WithContext(ctx).AutoMigrate(&models.GormUser{})
	if err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &PostgreSQLGORMRepository{db}
}

func (repo *PostgreSQLGORMRepository) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	gormUser := models.GormUser{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Username: user.Username,
	}

	if err := repo.db.WithContext(ctx).Create(&gormUser).Error; err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	result := models.User(gormUser)
	return &result, nil
}

func (repo *PostgreSQLGORMRepository) AllUsers(ctx context.Context) ([]models.User, error) {
	var allUsers []models.GormUser
	if err := repo.db.WithContext(ctx).Find(&allUsers).Error; err != nil {
		return nil, err
	}

	var result []models.User
	for _, users := range allUsers {
		result = append(result, models.User(users))
	}

	return result, nil
}

func (repo *PostgreSQLGORMRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	var gormUser models.GormUser
	if err := repo.db.WithContext(ctx).Where("id = ?", id).First(&gormUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	result := models.User(gormUser)
	return &result, nil
}

func (repo *PostgreSQLGORMRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var gormUser models.GormUser
	if err := repo.db.WithContext(ctx).Where("email = ?", email).First(&gormUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	result := models.User(gormUser)
	return &result, nil
}

func (repo *PostgreSQLGORMRepository) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (*models.User, error) {
	var gormUser models.GormUser
	if err := repo.db.WithContext(ctx).Where("username = ? AND password = ?", username, password).First(&gormUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExist
		}
		return nil, err
	}

	result := models.User(gormUser)
	return &result, nil
}

func (repo *PostgreSQLGORMRepository) UpdateUser(ctx context.Context, id int64, updated models.User) (*models.User, error) {
	gormUser := models.User(updated)
	updateRes := repo.db.WithContext(ctx).Where("id = ?", id).Save(&gormUser)
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

func (repo *PostgreSQLGORMRepository) DeleteUser(ctx context.Context, id int64) error {
	res := repo.db.WithContext(ctx).Delete(&models.GormUser{}, id)
	if err := res.Error; err != nil {
		return err
	}

	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return nil
}
