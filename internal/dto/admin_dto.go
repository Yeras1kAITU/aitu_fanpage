package dto

type AdminUserResponse struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	DisplayName  string `json:"display_name"`
	Role         string `json:"role"`
	IsActive     bool   `json:"is_active"`
	PostCount    int    `json:"post_count"`
	LikeCount    int    `json:"like_count"`
	CommentCount int    `json:"comment_count"`
	CreatedAt    string `json:"created_at"`
	LastLoginAt  string `json:"last_login_at,omitempty"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin student alumni moderator"`
}

type SystemStats struct {
	TotalUsers    int            `json:"total_users"`
	ActiveUsers   int            `json:"active_users"`
	NewUsersToday int            `json:"new_users_today"`
	TotalPosts    int            `json:"total_posts"`
	PostsToday    int            `json:"posts_today"`
	TotalComments int            `json:"total_comments"`
	TotalLikes    int            `json:"total_likes"`
	UsersByRole   map[string]int `json:"users_by_role"`
}
