package dto

type CreatePostRequest struct {
	Title       string               `json:"title"`
	Content     string               `json:"content"`
	Description string               `json:"description"`
	Category    string               `json:"category"`
	Media       []MediaUploadRequest `json:"media,omitempty"`
}

type MediaUploadRequest struct {
	URL     string `json:"url"`
	Type    string `json:"type"`
	Caption string `json:"caption,omitempty"`
}

type PostResponse struct {
	ID           string              `json:"id"`
	AuthorID     string              `json:"author_id"`
	AuthorName   string              `json:"author_name"`
	Title        string              `json:"title"`
	Content      string              `json:"content"`
	Description  string              `json:"description"`
	Category     string              `json:"category"`
	Media        []MediaItemResponse `json:"media,omitempty"`
	MediaCount   int                 `json:"media_count"`
	LikeCount    int                 `json:"like_count"`
	CommentCount int                 `json:"comment_count"`
	CreatedAt    string              `json:"created_at"`
}

type MediaItemResponse struct {
	URL      string `json:"url"`
	Type     string `json:"type"`
	Caption  string `json:"caption,omitempty"`
	Position int    `json:"position"`
}
