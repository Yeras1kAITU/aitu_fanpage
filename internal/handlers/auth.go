package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/dto"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/middleware"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrEmailAlreadyExists {
			status = http.StatusConflict
		} else if err == service.ErrInvalidEmail || err == service.ErrWeakPassword {
			status = http.StatusBadRequest
		}
		http.Error(w, err.Error(), status)
		return
	}

	token, err := h.authService.GetTokenAuth().GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := dto.AuthResponse{
		Token: token,
		User: dto.UserProfile{
			ID:           user.ID.Hex(),
			Email:        user.Email,
			DisplayName:  user.DisplayName,
			Role:         string(user.Role),
			ProfileImage: user.ProfileImage,
			Bio:          user.Bio,
			CreatedAt:    user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, user, err := h.authService.Login(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrInvalidCredentials {
			status = http.StatusUnauthorized
		}
		http.Error(w, err.Error(), status)
		return
	}

	response := dto.AuthResponse{
		Token: token,
		User: dto.UserProfile{
			ID:           user.ID.Hex(),
			Email:        user.Email,
			DisplayName:  user.DisplayName,
			Role:         string(user.Role),
			ProfileImage: user.ProfileImage,
			Bio:          user.Bio,
			CreatedAt:    user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := dto.UserProfile{
		ID:           user.ID.Hex(),
		Email:        user.Email,
		DisplayName:  user.DisplayName,
		Role:         string(user.Role),
		ProfileImage: user.ProfileImage,
		Bio:          user.Bio,
		CreatedAt:    user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.UpdateProfile(userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.UserProfile{
		ID:           user.ID.Hex(),
		Email:        user.Email,
		DisplayName:  user.DisplayName,
		Role:         string(user.Role),
		ProfileImage: user.ProfileImage,
		Bio:          user.Bio,
		CreatedAt:    user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.authService.ChangePassword(userID, req); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "current password is incorrect" {
			status = http.StatusBadRequest
		} else if err == service.ErrWeakPassword {
			status = http.StatusBadRequest
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password updated successfully",
	})
}
