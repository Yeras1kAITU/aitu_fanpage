package service

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/models"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/repository"
)

type CommentService struct {
	commentRepo repository.CommentRepository
	userRepo    repository.UserRepository
	postRepo    repository.PostRepository
}

func NewCommentService(commentRepo repository.CommentRepository, userRepo repository.UserRepository, postRepo repository.PostRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		userRepo:    userRepo,
		postRepo:    postRepo,
	}
}

func (s *CommentService) CreateComment(postID, authorID primitive.ObjectID, content string) (*models.Comment, error) {
	user, err := s.userRepo.FindByID(authorID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	comment := models.NewComment(postID, authorID, user.DisplayName, content)

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	user.IncrementCommentCount()
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	if err := s.postRepo.IncrementCommentCount(postID); err != nil {
		user.DecrementCommentCount()
		s.userRepo.Update(user)
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) GetCommentsByPostID(postID primitive.ObjectID, limit, offset int) ([]*models.Comment, error) {
	return s.commentRepo.FindByPostID(postID, limit, offset)
}

func (s *CommentService) UpdateComment(commentID, userID primitive.ObjectID, content string) (*models.Comment, error) {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if comment.AuthorID != userID && !user.CanManageComments() {
		return nil, errors.New("not authorized to edit this comment")
	}

	comment.Content = content
	if err := s.commentRepo.Update(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) DeleteComment(commentID, userID primitive.ObjectID) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if comment.AuthorID != userID && !user.CanManageComments() {
		return errors.New("not authorized to delete this comment")
	}

	if err := s.commentRepo.Delete(commentID); err != nil {
		return err
	}

	user.DecrementCommentCount()
	s.userRepo.Update(user)

	if err := s.postRepo.DecrementCommentCount(comment.PostID); err != nil {
		return err
	}

	return nil
}

func (s *CommentService) DeleteCommentByPostID(postID primitive.ObjectID) error {
	return s.commentRepo.DeleteByPostID(postID)
}
