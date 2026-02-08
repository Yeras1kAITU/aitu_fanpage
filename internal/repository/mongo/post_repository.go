package mongorepo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/models"
)

type PostRepository struct {
	collection *mongo.Collection
}

func NewPostRepository(db *mongo.Database) *PostRepository {
	return &PostRepository{
		collection: db.Collection("posts"),
	}
}

func (r *PostRepository) Create(post *models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, post)
	return err
}

func (r *PostRepository) FindByID(id primitive.ObjectID) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var post models.Post
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) FindAll(limit, offset int) ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{"is_archived": false}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*models.Post
	for cursor.Next(ctx) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) IncrementLikeCount(postID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": postID}
	update := bson.M{
		"$inc": bson.M{"like_count": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *PostRepository) DecrementLikeCount(postID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": postID}
	update := bson.M{
		"$inc": bson.M{"like_count": -1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)

	if err == nil {
		fixFilter := bson.M{"_id": postID, "like_count": bson.M{"$lt": 0}}
		fixUpdate := bson.M{"$set": bson.M{"like_count": 0}}
		r.collection.UpdateOne(ctx, fixFilter, fixUpdate)
	}

	return err
}

func (r *PostRepository) FindByCategory(category string, limit int) ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetLimit(int64(limit))

	cursor, err := r.collection.Find(
		ctx,
		bson.M{"category": category, "is_archived": false},
		findOptions,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*models.Post
	for cursor.Next(ctx) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) Update(post *models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	post.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": post.ID},
		bson.M{"$set": post},
	)
	return err
}

func (r *PostRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *PostRepository) FindByAuthor(authorID primitive.ObjectID, limit int) ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetLimit(int64(limit))

	cursor, err := r.collection.Find(
		ctx,
		bson.M{"author_id": authorID, "is_archived": false},
		findOptions,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*models.Post
	for cursor.Next(ctx) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) FindPinned(limit int) ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "pinned_at", Value: -1}})
	findOptions.SetLimit(int64(limit))

	cursor, err := r.collection.Find(
		ctx,
		bson.M{"is_pinned": true, "is_archived": false},
		findOptions,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*models.Post
	for cursor.Next(ctx) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) FindFeatured(limit int) ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "featured_at", Value: -1}})
	findOptions.SetLimit(int64(limit))

	cursor, err := r.collection.Find(
		ctx,
		bson.M{"is_featured": true, "is_archived": false},
		findOptions,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*models.Post
	for cursor.Next(ctx) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) FindPopular(limit int, days int) ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	since := time.Now().AddDate(0, 0, -days)

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "popularity_score", Value: -1}})
	findOptions.SetLimit(int64(limit))

	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"created_at":  bson.M{"$gte": since},
			"is_archived": false,
		},
		findOptions,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*models.Post
	for cursor.Next(ctx) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) FindByTags(tags []string, limit int) ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetLimit(int64(limit))

	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"tags":        bson.M{"$in": tags},
			"is_archived": false,
		},
		findOptions,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*models.Post
	for cursor.Next(ctx) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) Search(query string, limit int) ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"$text":       bson.M{"$search": query},
		"is_archived": false,
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "score", Value: bson.M{"$meta": "textScore"}}})
	findOptions.SetLimit(int64(limit))
	findOptions.SetProjection(bson.M{
		"score": bson.M{"$meta": "textScore"},
	})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return r.simpleSearch(query, limit)
	}
	defer cursor.Close(ctx)

	var posts []*models.Post
	for cursor.Next(ctx) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) simpleSearch(query string, limit int) ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"content": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
			{"tags": bson.M{"$regex": query, "$options": "i"}},
		},
		"is_archived": false,
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*models.Post
	for cursor.Next(ctx) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) IncrementViewCount(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$inc": bson.M{"view_count": 1}},
	)
	return err
}

func (r *PostRepository) GetCategoriesStats() (map[string]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := []bson.M{
		{
			"$match": bson.M{"is_archived": false},
		},
		{
			"$group": bson.M{
				"_id":   "$category",
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"count": -1},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	stats := make(map[string]int)
	for cursor.Next(ctx) {
		var result struct {
			Category string `bson:"_id"`
			Count    int    `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		stats[result.Category] = result.Count
	}

	return stats, nil
}
