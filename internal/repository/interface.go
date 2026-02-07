package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/models"
)

type PostRepository interface {
	Create(post *models.Post) error
	FindByID(id primitive.ObjectID) (*models.Post, error)
	FindAll(limit, offset int) ([]*models.Post, error)
	IncrementLikeCount(id primitive.ObjectID) error
	FindByCategory(category string, limit int) ([]*models.Post, error)
	Update(post *models.Post) error
	Delete(id primitive.ObjectID) error
	FindByAuthor(authorID primitive.ObjectID, limit int) ([]*models.Post, error)
	FindPinned(limit int) ([]*models.Post, error)
	FindFeatured(limit int) ([]*models.Post, error)
	FindPopular(limit int, days int) ([]*models.Post, error)
	FindByTags(tags []string, limit int) ([]*models.Post, error)
	Search(query string, limit int) ([]*models.Post, error)
	IncrementViewCount(id primitive.ObjectID) error
	GetCategoriesStats() (map[string]int, error)
}

type UserRepository interface {
	FindByID(id primitive.ObjectID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id primitive.ObjectID) error
	FindAll(limit, offset int) ([]*models.User, error)
}

type CommentRepository interface {
	Create(comment *models.Comment) error
	FindByID(id primitive.ObjectID) (*models.Comment, error)
	FindByPostID(postID primitive.ObjectID, limit, offset int) ([]*models.Comment, error)
	Update(comment *models.Comment) error
	Delete(id primitive.ObjectID) error
	DeleteByPostID(postID primitive.ObjectID) error
	CountByPostID(postID primitive.ObjectID) (int64, error)
}
