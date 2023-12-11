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

type PostService struct {
	PostRepo repository.PostRepository
	db       *gorm.DB
}

func NewPostService(postRepo repository.PostRepository, db *gorm.DB) *PostService {
	return &PostService{
		PostRepo: postRepo,
		db:       db,
	}
}

func (postService *PostService) CreatePost(ctx context.Context, post models.Post) (*models.Post, error) {
	_, err := postService.PostRepo.GetPostByTitle(ctx, post.Title)
	if err == nil || !errors.Is(err, repository.ErrNotExist) {
		return nil, errors.New("a post with this title already exists")
	}

	return postService.PostRepo.CreatePost(ctx, post)
}

func (postService *PostService) GetAllPosts(ctx context.Context) ([]models.Post, error) {
	var allPosts []models.GormPost
	if err := postService.db.WithContext(ctx).Find(&allPosts).Error; err != nil {
		return nil, err
	}

	var result []models.Post
	for _, post := range allPosts {
		result = append(result, models.Post(post))
	}

	return result, nil
}

func (postService *PostService) GetPostByID(ctx context.Context, id int64) (*models.Post, error) {
	post, err := postService.PostRepo.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}

	return post, nil
}

func (postService *PostService) GetPostByTitle(ctx context.Context, title string) (*models.Post, error) {
	post, err := postService.PostRepo.GetPostByTitle(ctx, title)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}

	return post, nil
}

func (postService *PostService) GetPostByUserID(ctx context.Context, userid int64) ([]models.Post, error) {
	post, err := postService.PostRepo.GetPostByUserID(ctx, userid)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}
	return post, nil
}

func (postService *PostService) UpdatePostByID(ctx context.Context, post models.Post) (*models.Post, error) {
	existingPost, err := postService.PostRepo.GetPostByID(ctx, post.ID)
	if err != nil {
		return nil, err
	}
	if _, err := postService.PostRepo.UpdatePost(ctx, post.ID, post); err != nil {
		log.Printf("Error updating post with ID %d: %v", post.ID, err)
		return nil, err
	}
	return existingPost, nil
}

func (postService *PostService) DeletePostByID(ctx context.Context, id int64) error {
	if err := postService.PostRepo.DeletePost(ctx, id); err != nil {
		log.Printf("Error deleting post with ID %d: %v", id, err)
		return err
	}
	return nil
}
