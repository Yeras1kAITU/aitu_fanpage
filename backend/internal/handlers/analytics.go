package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/dto"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/service"
)

type AnalyticsHandler struct {
	postService *service.PostService
}

func NewAnalyticsHandler(postService *service.PostService) *AnalyticsHandler {
	return &AnalyticsHandler{
		postService: postService,
	}
}

func (h *AnalyticsHandler) GetCategoriesStatsAggregated(w http.ResponseWriter, r *http.Request) {

	stats, err := h.postService.GetCategoriesStatsAggregated()
	if err != nil {
		http.Error(w, "Failed to get category stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	totalPosts := 0
	totalLikes := 0

	for _, stat := range stats {
		totalPosts += stat.Count
		totalLikes += stat.TotalLikes
	}

	response := dto.CategoryStatsResponse{
		Categories: make(map[string]dto.CategoryStat),
		TotalPosts: totalPosts,
		TotalLikes: totalLikes,
		AvgLikes:   0,
	}

	if totalPosts > 0 {
		response.AvgLikes = float64(totalLikes) / float64(totalPosts)
	}

	for category, stat := range stats {
		percentage := 0.0
		if totalPosts > 0 {
			percentage = (float64(stat.Count) / float64(totalPosts)) * 100
		}

		response.Categories[category] = dto.CategoryStat{
			Count:         stat.Count,
			TotalLikes:    stat.TotalLikes,
			AvgLikes:      stat.AvgLikes,
			TotalComments: stat.TotalComments,
			AvgComments:   stat.AvgComments,
			Percentage:    percentage,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
