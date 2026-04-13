package dto

type CategoryStatsResponse struct {
	Categories map[string]CategoryStat `json:"categories"`
	TotalPosts int                     `json:"total_posts"`
	TotalLikes int                     `json:"total_likes"`
	AvgLikes   float64                 `json:"avg_likes_per_post"`
}

type CategoryStat struct {
	Count         int     `json:"count"`
	TotalLikes    int     `json:"total_likes"`
	AvgLikes      float64 `json:"avg_likes"`
	TotalComments int     `json:"total_comments"`
	AvgComments   float64 `json:"avg_comments"`
	Percentage    float64 `json:"percentage"`
}
