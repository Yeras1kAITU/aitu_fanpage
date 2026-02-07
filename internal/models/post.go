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
	CategoryAcademic  PostCategory = "academic"
	CategorySocial    PostCategory = "social"
	CategorySports    PostCategory = "sports"
)

type MediaItem struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	URL          string             `bson:"url" json:"url"`
	Type         string             `bson:"type" json:"type"`
	Caption      string             `bson:"caption,omitempty" json:"caption"`
	Position     int                `bson:"position" json:"position"`
	FileSize     int64              `bson:"file_size,omitempty" json:"file_size,omitempty"`
	Checksum     string             `bson:"checksum,omitempty" json:"checksum,omitempty"`
	ThumbnailURL string             `bson:"thumbnail_url,omitempty" json:"thumbnail_url,omitempty"`
	Dimensions   string             `bson:"dimensions,omitempty" json:"dimensions,omitempty"`
	Duration     string             `bson:"duration,omitempty" json:"duration,omitempty"`
	CreatedAt    time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

type Post struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthorID        primitive.ObjectID `bson:"author_id" json:"author_id"`
	AuthorName      string             `bson:"author_name" json:"author_name"`
	Title           string             `bson:"title" json:"title"`
	Content         string             `bson:"content" json:"content"`
	Description     string             `bson:"description" json:"description"`
	Category        PostCategory       `bson:"category" json:"category"`
	Tags            []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	Media           []MediaItem        `bson:"media,omitempty" json:"media,omitempty"`
	MediaCount      int                `bson:"media_count" json:"media_count"`
	LikeCount       int                `bson:"like_count" json:"like_count"`
	CommentCount    int                `bson:"comment_count" json:"comment_count"`
	ViewCount       int                `bson:"view_count" json:"view_count"`
	IsFeatured      bool               `bson:"is_featured" json:"is_featured"`
	IsPinned        bool               `bson:"is_pinned" json:"is_pinned"`
	IsArchived      bool               `bson:"is_archived" json:"is_archived"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	FeaturedAt      *time.Time         `bson:"featured_at,omitempty" json:"featured_at,omitempty"`
	PinnedAt        *time.Time         `bson:"pinned_at,omitempty" json:"pinned_at,omitempty"`
	PopularityScore float64            `bson:"popularity_score" json:"popularity_score"`
}

func NewPost(title, content, description string, category PostCategory, authorID primitive.ObjectID, authorName string) *Post {
	now := time.Now()
	return &Post{
		ID:              primitive.NewObjectID(),
		AuthorID:        authorID,
		AuthorName:      authorName,
		Title:           title,
		Content:         content,
		Description:     description,
		Category:        category,
		Tags:            []string{},
		LikeCount:       0,
		CommentCount:    0,
		ViewCount:       0,
		IsFeatured:      false,
		IsPinned:        false,
		IsArchived:      false,
		CreatedAt:       now,
		UpdatedAt:       now,
		PopularityScore: 0,
	}
}

func (p *Post) AddMedia(url, mediaType, caption string, fileSize int64, thumbnailURL, dimensions, checksum string) {
	media := MediaItem{
		ID:           primitive.NewObjectID(),
		URL:          url,
		Type:         mediaType,
		Caption:      caption,
		Position:     len(p.Media),
		FileSize:     fileSize,
		Checksum:     checksum,
		ThumbnailURL: thumbnailURL,
		Dimensions:   dimensions,
	}
	p.Media = append(p.Media, media)
	p.MediaCount = len(p.Media)
}

func (p *Post) AddTags(tags ...string) {
	for _, tag := range tags {
		found := false
		for _, existingTag := range p.Tags {
			if existingTag == tag {
				found = true
				break
			}
		}
		if !found {
			p.Tags = append(p.Tags, tag)
		}
	}
}

func (p *Post) CalculatePopularityScore() {
	hoursSinceCreation := time.Since(p.CreatedAt).Hours()
	if hoursSinceCreation < 1 {
		hoursSinceCreation = 1
	}

	likeWeight := 2.0
	commentWeight := 3.0
	viewWeight := 0.1
	recencyWeight := 0.5

	p.PopularityScore = (float64(p.LikeCount)*likeWeight +
		float64(p.CommentCount)*commentWeight +
		float64(p.ViewCount)*viewWeight) /
		(hoursSinceCreation * recencyWeight)
}

func (p *Post) IncrementViewCount() {
	p.ViewCount++
	p.CalculatePopularityScore()
	p.UpdatedAt = time.Now()
}

func (p *Post) Pin() {
	now := time.Now()
	p.IsPinned = true
	p.PinnedAt = &now
	p.UpdatedAt = now
}

func (p *Post) Unpin() {
	p.IsPinned = false
	p.PinnedAt = nil
	p.UpdatedAt = time.Now()
}

func (p *Post) Feature() {
	now := time.Now()
	p.IsFeatured = true
	p.FeaturedAt = &now
	p.UpdatedAt = now
}

func (p *Post) Unfeature() {
	p.IsFeatured = false
	p.FeaturedAt = nil
	p.UpdatedAt = time.Now()
}

func (p *Post) Archive() {
	p.IsArchived = true
	p.UpdatedAt = time.Now()
}

func (p *Post) Unarchive() {
	p.IsArchived = false
	p.UpdatedAt = time.Now()
}

func (p *Post) RemoveMedia(position int) bool {
	if position < 0 || position >= len(p.Media) {
		return false
	}

	p.Media = append(p.Media[:position], p.Media[position+1:]...)

	// Update positions
	for i := position; i < len(p.Media); i++ {
		p.Media[i].Position = i
	}

	p.MediaCount = len(p.Media)
	return true
}

func (p *Post) UpdateMediaCaption(position int, caption string) bool {
	if position < 0 || position >= len(p.Media) {
		return false
	}

	p.Media[position].Caption = caption
	return true
}

func (p *Post) GetFirstImageURL() string {
	for _, media := range p.Media {
		if media.Type == "image" {
			return media.URL
		}
	}
	return ""
}

func (p *Post) GetThumbnailURL() string {
	for _, media := range p.Media {
		if media.ThumbnailURL != "" {
			return media.ThumbnailURL
		}
		if media.Type == "image" {
			return media.URL
		}
	}
	return ""
}
