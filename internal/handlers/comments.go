package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/dto"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/middleware"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/service"
)

type CommentHandler struct {
	service *service.CommentService
}

func NewCommentHandler(service *service.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postIDStr := chi.URLParam(r, "id")
	postID, err := primitive.ObjectIDFromHex(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var req dto.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Comment content is required", http.StatusBadRequest)
		return
	}

	comment, err := h.service.CreateComment(postID, userID, req.Content)
	if err != nil {
		http.Error(w, "Failed to create comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.CommentResponse{
		ID:         comment.ID.Hex(),
		PostID:     comment.PostID.Hex(),
		AuthorID:   comment.AuthorID.Hex(),
		AuthorName: comment.AuthorName,
		Content:    comment.Content,
		CreatedAt:  comment.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:  comment.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := primitive.ObjectIDFromHex(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

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

	comments, err := h.service.GetCommentsByPostID(postID, limit, offset)
	if err != nil {
		http.Error(w, "Failed to get comments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []dto.CommentResponse
	for _, comment := range comments {
		responses = append(responses, dto.CommentResponse{
			ID:         comment.ID.Hex(),
			PostID:     comment.PostID.Hex(),
			AuthorID:   comment.AuthorID.Hex(),
			AuthorName: comment.AuthorName,
			Content:    comment.Content,
			CreatedAt:  comment.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:  comment.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	commentIDStr := chi.URLParam(r, "id")
	commentID, err := primitive.ObjectIDFromHex(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Comment content is required", http.StatusBadRequest)
		return
	}

	comment, err := h.service.UpdateComment(commentID, userID, req.Content)
	if err != nil {
		if err.Error() == "not authorized to edit this comment" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to update comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.CommentResponse{
		ID:         comment.ID.Hex(),
		PostID:     comment.PostID.Hex(),
		AuthorID:   comment.AuthorID.Hex(),
		AuthorName: comment.AuthorName,
		Content:    comment.Content,
		CreatedAt:  comment.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:  comment.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	commentIDStr := chi.URLParam(r, "id")
	commentID, err := primitive.ObjectIDFromHex(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteComment(commentID, userID)
	if err != nil {
		if err.Error() == "not authorized to delete this comment" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to delete comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Comment deleted successfully",
	})
}

func (h *CommentHandler) GetCommentCount(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := primitive.ObjectIDFromHex(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	count, err := h.service.GetCommentCount(postID)
	if err != nil {
		http.Error(w, "Failed to get comment count: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"post_id":       postIDStr,
		"comment_count": count,
	})
}
