package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostID     primitive.ObjectID `bson:"post_id" json:"post_id"`
	AuthorID   primitive.ObjectID `bson:"author_id" json:"author_id"`
	AuthorName string             `bson:"author_name" json:"author_name"`
	Content    string             `bson:"content" json:"content"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

func NewComment(postID, authorID primitive.ObjectID, authorName, content string) *Comment {
	now := time.Now()
	return &Comment{
		ID:         primitive.NewObjectID(),
		PostID:     postID,
		AuthorID:   authorID,
		AuthorName: authorName,
		Content:    content,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
