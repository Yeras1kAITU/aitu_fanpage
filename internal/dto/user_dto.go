package dto

type PublicUserProfile struct {
	ID           string `json:"id"`
	DisplayName  string `json:"display_name"`
	Role         string `json:"role"`
	ProfileImage string `json:"profile_image,omitempty"`
	Bio          string `json:"bio,omitempty"`
	PostCount    int    `json:"post_count"`
	LikeCount    int    `json:"like_count"`
	CommentCount int    `json:"comment_count"`
	CreatedAt    string `json:"created_at"`
}
