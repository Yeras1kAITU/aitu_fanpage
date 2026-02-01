package mongorepo

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/models"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/repository"
)

type postRepository struct {
	collection *mongo.Collection
	mu         sync.RWMutex
}

func NewPostRepository(db *mongo.Database) repository.PostRepository {
	return &postRepository{
		collection: db.Collection("posts"),
	}
}

func (r *postRepository) Create(post *models.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.collection.InsertOne(context.Background(), post)
	return err
}

func (r *postRepository) FindByID(id primitive.ObjectID) (*models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var post models.Post
	filter := bson.M{"_id": id}

	err := r.collection.FindOne(context.Background(), filter).Decode(&post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}
