package dto

type CreatePostRequest struct {
	Title       string               `json:"title"`
	Content     string               `json:"content"`
	Description string               `json:"description"`
	Category    string               `json:"category"`
	Tags        []string             `json:"tags,omitempty"`
	Media       []MediaUploadRequest `json:"media,omitempty"`
}

type UpdatePostRequest struct {
	Title       string   `json:"title,omitempty"`
	Content     string   `json:"content,omitempty"`
	Description string   `json:"description,omitempty"`
	Category    string   `json:"category,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type PostResponse struct {
	ID              string              `json:"id"`
	AuthorID        string              `json:"author_id"`
	AuthorName      string              `json:"author_name"`
	Title           string              `json:"title"`
	Content         string              `json:"content"`
	Description     string              `json:"description"`
	Category        string              `json:"category"`
	Tags            []string            `json:"tags,omitempty"`
	Media           []MediaItemResponse `json:"media,omitempty"` // Uses MediaItemResponse from media_dto.go
	MediaCount      int                 `json:"media_count"`
	LikeCount       int                 `json:"like_count"`
	CommentCount    int                 `json:"comment_count"`
	ViewCount       int                 `json:"view_count"`
	IsFeatured      bool                `json:"is_featured"`
	IsPinned        bool                `json:"is_pinned"`
	CreatedAt       string              `json:"created_at"`
	UpdatedAt       string              `json:"updated_at"`
	PopularityScore float64             `json:"popularity_score,omitempty"`
}

type PostFilterRequest struct {
	Category string   `json:"category,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	AuthorID string   `json:"author_id,omitempty"`
	SortBy   string   `json:"sort_by,omitempty"` // recent, popular, featured
	Limit    int      `json:"limit,omitempty"`
	Offset   int      `json:"offset,omitempty"`
}

type CategoriesStatsResponse struct {
	Categories map[string]int `json:"categories"`
	TotalPosts int            `json:"total_posts"`
}

type LikeResponse struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	LikedAt  string `json:"liked_at"`
}

type SearchResponse struct {
	Posts []PostResponse `json:"posts"`
	Total int            `json:"total"`
}
