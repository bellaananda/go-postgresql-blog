package service

import (
	"context"
	"errors"
	"log"

	// "fmt"
	// "log"
	"postgresql-blog/models"
	"postgresql-blog/repository"

	"gorm.io/gorm"
)

// func CreateUser(ctx context.Context, userRepository repository.UserRepository) {
// 	inputdummy := models.User{
// 		Name:     "Dummy Doo",
// 		Email:    "dummy@mail.com",
// 		Password: "dummy123",
// 		Username: "dummydoo",
// 	}

// 	createdUser, err := userRepository.Create(ctx, inputdummy)
// 	if errors.Is(err, repository.ErrDuplicate) {
// 		fmt.Printf("record: %+v already exists\n", inputdummy)
// 	} else if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("created record: %+v\n", createdUser)
// }

// service/user_service.go

type UserService struct {
	UserRepo repository.UserRepository
	db       *gorm.DB
}

func NewUserService(userRepo repository.UserRepository, db *gorm.DB) *UserService {
	return &UserService{
		UserRepo: userRepo,
		db:       db,
	}
}

func (userService *UserService) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	_, err := userService.UserRepo.GetUserByEmail(ctx, user.Email)
	if err == nil || !errors.Is(err, repository.ErrNotExist) {
		return nil, errors.New("user with this email already exists")
	}

	return userService.UserRepo.CreateUser(ctx, user)
}

func (userService *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var allUsers []models.GormUser
	if err := userService.db.WithContext(ctx).Find(&allUsers).Error; err != nil {
		return nil, err
	}

	var result []models.User
	for _, user := range allUsers {
		result = append(result, models.User(user))
	}

	return result, nil
}

func (userService *UserService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	user, err := userService.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}

	// log.Printf("User found by id '%d': %+v\n", id, user)
	return user, nil
}

func (userService *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := userService.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}

	return user, nil
}

func (userService *UserService) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (*models.User, error) {
	user, err := userService.UserRepo.GetUserByUsernameAndPassword(ctx, username, password)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}

	return user, nil
}

func (userService *UserService) UpdateUserByID(ctx context.Context, user models.User) (*models.User, error) {
	existingUser, err := userService.UserRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if _, err := userService.UserRepo.UpdateUser(ctx, user.ID, user); err != nil {
		log.Printf("Error updating user with ID %d: %v", user.ID, err)
		return nil, err
	}
	return existingUser, nil
}

func (userService *UserService) DeleteUserByID(ctx context.Context, id int64) error {
	if err := userService.UserRepo.DeleteUser(ctx, id); err != nil {
		log.Printf("Error deleting user with ID %d: %v", id, err)
		return err
	}
	return nil
}
