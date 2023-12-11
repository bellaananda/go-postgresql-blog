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

type CommentService struct {
	CommentRepo repository.CommentRepository
	db          *gorm.DB
}

func NewCommentService(commentRepo repository.CommentRepository, db *gorm.DB) *CommentService {
	return &CommentService{
		CommentRepo: commentRepo,
		db:          db,
	}
}

func (commentService *CommentService) CreateComment(ctx context.Context, comment models.Comment) (*models.Comment, error) {
	_, err := commentService.CommentRepo.GetCommentByUserIDPostID(ctx, int64(comment.UserID), int64(comment.PostID))
	if err == nil || !errors.Is(err, repository.ErrNotExist) {
		return nil, errors.New("a comment with the post already exists")
	}

	return commentService.CommentRepo.CreateComment(ctx, comment)
}

func (commentService *CommentService) GetAllComments(ctx context.Context) ([]models.Comment, error) {
	var allComments []models.GormComment
	if err := commentService.db.WithContext(ctx).Find(&allComments).Error; err != nil {
		return nil, err
	}

	var result []models.Comment
	for _, comment := range allComments {
		result = append(result, models.Comment(comment))
	}

	return result, nil
}

func (commentService *CommentService) GetCommentByID(ctx context.Context, id int64) (*models.Comment, error) {
	comment, err := commentService.CommentRepo.GetCommentByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}

	return comment, nil
}

func (commentService *CommentService) GetCommentByUserID(ctx context.Context, userid int64) ([]models.Comment, error) {
	comment, err := commentService.CommentRepo.GetCommentByUserID(ctx, userid)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}
	return comment, nil
}

func (commentService *CommentService) GetCommentByPostID(ctx context.Context, postid int64) ([]models.Comment, error) {
	comment, err := commentService.CommentRepo.GetCommentByPostID(ctx, postid)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			return nil, err
		}
		return nil, err
	}

	return comment, nil
}

func (commentService *CommentService) GetCommentByUserIDPostID(ctx context.Context, userid int64, postid int64) (*models.Comment, error) {
	comment, err := commentService.CommentRepo.GetCommentByUserIDPostID(ctx, userid, postid)
	if err != nil {
		if errors.Is(err, repository.ErrNotExist) {
			log.Printf("Post with user id '%d' and post id '%d' does not exist in the repository\n", userid, postid)
			return nil, err
		}
		log.Printf("Error getting post by user id '%d' and post id '%d': %v\n", userid, postid, err)
		return nil, err
	}
	log.Printf("Post found by user id '%d' and post id '%d': %+v\n", userid, postid, comment)
	return comment, nil
}

func (commentService *CommentService) UpdateCommentByID(ctx context.Context, comment models.Comment) (*models.Comment, error) {
	existingComment, err := commentService.CommentRepo.GetCommentByID(ctx, comment.ID)
	if err != nil {
		return nil, err
	}
	if _, err := commentService.CommentRepo.UpdateComment(ctx, comment.ID, comment); err != nil {
		return nil, err
	}
	return existingComment, nil
}

func (commentService *CommentService) DeleteCommentByID(ctx context.Context, id int64) error {
	if err := commentService.CommentRepo.DeleteComment(ctx, id); err != nil {
		return err
	}
	return nil
}
