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
}

type UserRepository interface {
	FindByID(id primitive.ObjectID) (*models.User, error)
}
