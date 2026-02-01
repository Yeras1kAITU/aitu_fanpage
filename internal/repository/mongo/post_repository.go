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

func (r *postRepository) FindAll(limit, offset int) ([]*models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var posts []*models.Post
	if err := cursor.All(context.Background(), &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepository) IncrementLikeCount(id primitive.ObjectID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"like_count": 1}}

	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *postRepository) FindByCategory(category string, limit int) ([]*models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filter := bson.M{"category": category}
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var posts []*models.Post
	if err := cursor.All(context.Background(), &posts); err != nil {
		return nil, err
	}

	return posts, nil
}
