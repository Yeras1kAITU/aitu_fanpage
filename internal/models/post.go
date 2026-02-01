package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostCategory string

const (
	CategoryMeme      PostCategory = "meme"
	CategoryEvent     PostCategory = "event"
	CategoryNews      PostCategory = "news"
	CategoryQuestion  PostCategory = "question"
	CategoryLostFound PostCategory = "lost_found"
)

type MediaItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	URL      string             `bson:"url" json:"url"`
	Type     string             `bson:"type" json:"type"`
	Caption  string             `bson:"caption,omitempty" json:"caption"`
	Position int                `bson:"position" json:"position"`
}

type Post struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthorID     primitive.ObjectID `bson:"author_id" json:"author_id"`
	AuthorName   string             `bson:"author_name" json:"author_name"`
	Title        string             `bson:"title" json:"title"`
	Content      string             `bson:"content" json:"content"`
	Description  string             `bson:"description" json:"description"`
	Category     PostCategory       `bson:"category" json:"category"`
	Media        []MediaItem        `bson:"media,omitempty" json:"media,omitempty"`
	MediaCount   int                `bson:"media_count" json:"media_count"`
	LikeCount    int                `bson:"like_count" json:"like_count"`
	CommentCount int                `bson:"comment_count" json:"comment_count"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

func NewPost(title, content, description string, category PostCategory, authorID primitive.ObjectID, authorName string) *Post {
	now := time.Now()
	return &Post{
		ID:           primitive.NewObjectID(),
		AuthorID:     authorID,
		AuthorName:   authorName,
		Title:        title,
		Content:      content,
		Description:  description,
		Category:     category,
		LikeCount:    0,
		CommentCount: 0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (p *Post) AddMedia(url, mediaType, caption string) {
	media := MediaItem{
		ID:       primitive.NewObjectID(),
		URL:      url,
		Type:     mediaType,
		Caption:  caption,
		Position: len(p.Media),
	}
	p.Media = append(p.Media, media)
	p.MediaCount = len(p.Media)
}
