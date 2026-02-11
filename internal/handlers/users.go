package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/dto"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/middleware"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := dto.PublicUserProfile{
		ID:           user.ID.Hex(),
		DisplayName:  user.DisplayName,
		Role:         string(user.Role),
		ProfileImage: user.ProfileImage,
		Bio:          user.Bio,
		PostCount:    user.PostCount,
		LikeCount:    user.LikeCount,
		CommentCount: user.CommentCount,
		CreatedAt:    user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	currentUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userIDStr := chi.URLParam(r, "id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Check if user is requesting their own stats or is admin
	currentUser, err := h.service.GetUserByID(currentUserID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if currentUserID != userID && !currentUser.CanViewAnalytics() {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	stats, err := h.service.GetUserStats(userID)
	if err != nil {
		http.Error(w, "Failed to get user stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
