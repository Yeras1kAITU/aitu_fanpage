package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/dto"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/middleware"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/models"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/service"
)

type AdminHandler struct {
	postService    *service.PostService
	userService    *service.UserService
	commentService *service.CommentService
}

func NewAdminHandler(postService *service.PostService, userService *service.UserService, commentService *service.CommentService) *AdminHandler {
	return &AdminHandler{
		postService:    postService,
		userService:    userService,
		commentService: commentService,
	}
}

func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	users, err := h.userService.GetAllUsers(limit, offset)
	if err != nil {
		http.Error(w, "Failed to get users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []dto.AdminUserResponse
	for _, user := range users {
		lastLoginAt := ""
		if !user.LastLoginAt.IsZero() {
			lastLoginAt = user.LastLoginAt.Format("2006-01-02T15:04:05Z")
		}

		responses = append(responses, dto.AdminUserResponse{
			ID:           user.ID.Hex(),
			Email:        user.Email,
			DisplayName:  user.DisplayName,
			Role:         string(user.Role),
			IsActive:     user.IsActive,
			PostCount:    user.PostCount,
			LikeCount:    user.LikeCount,
			CommentCount: user.CommentCount,
			CreatedAt:    user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			LastLoginAt:  lastLoginAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *AdminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetUserIDFromContext(r.Context())
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

	var req dto.UpdateUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.userService.UpdateUserRole(adminID, userID, models.UserRole(req.Role)); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrPermissionDenied) {
			status = http.StatusForbidden
		} else if errors.Is(err, service.ErrInvalidRole) {
			status = http.StatusBadRequest
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "User role updated successfully",
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *AdminHandler) ToggleUserStatus(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetUserIDFromContext(r.Context())
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

	action := chi.URLParam(r, "action")
	var serviceErr error

	switch action {
	case "deactivate":
		serviceErr = h.userService.DeactivateUser(adminID, userID)
	case "activate":
		serviceErr = h.userService.ActivateUser(adminID, userID)
	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	if serviceErr != nil {
		status := http.StatusInternalServerError
		if errors.Is(serviceErr, service.ErrPermissionDenied) {
			status = http.StatusForbidden
		} else if errors.Is(serviceErr, service.ErrCannotDeactivateSelf) {
			status = http.StatusBadRequest
		}
		http.Error(w, serviceErr.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "User status updated successfully",
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetUserIDFromContext(r.Context())
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

	if err := h.userService.DeleteUser(adminID, userID); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrPermissionDenied) {
			status = http.StatusForbidden
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "User deleted successfully",
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *AdminHandler) GetSystemStats(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	admin, err := h.userService.GetUserByID(adminID)
	if err != nil || !admin.CanViewAnalytics() {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	allUsers, err := h.userService.GetAllUsers(1000, 0)
	if err != nil {
		http.Error(w, "Failed to get stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	stats := dto.SystemStats{
		TotalUsers:    len(allUsers),
		ActiveUsers:   0,
		NewUsersToday: 0,
		TotalPosts:    0,
		PostsToday:    0,
		TotalComments: 0,
		TotalLikes:    0,
		UsersByRole:   make(map[string]int),
	}

	today := time.Now().Truncate(24 * time.Hour)

	for _, user := range allUsers {
		if user.IsActive {
			stats.ActiveUsers++
		}

		if user.CreatedAt.After(today) {
			stats.NewUsersToday++
		}

		stats.TotalPosts += user.PostCount
		stats.TotalComments += user.CommentCount
		stats.TotalLikes += user.LikeCount

		role := string(user.Role)
		stats.UsersByRole[role] = stats.UsersByRole[role] + 1
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *AdminHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	admin, err := h.userService.GetUserByID(adminID)
	if err != nil || !admin.CanManageUsers() {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	users, err := h.userService.SearchUsers(query, limit)
	if err != nil {
		http.Error(w, "Failed to search users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []dto.AdminUserResponse
	for _, user := range users {
		responses = append(responses, dto.AdminUserResponse{
			ID:           user.ID.Hex(),
			Email:        user.Email,
			DisplayName:  user.DisplayName,
			Role:         string(user.Role),
			IsActive:     user.IsActive,
			PostCount:    user.PostCount,
			LikeCount:    user.LikeCount,
			CommentCount: user.CommentCount,
			CreatedAt:    user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
