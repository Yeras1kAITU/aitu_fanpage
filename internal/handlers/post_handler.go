package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/dto"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/service"
)

type PostHandler struct {
	service *service.PostService
}

type HandlerContainer struct {
	Post *PostHandler
}

func NewPostHandler(service *service.PostService) *PostHandler {
	return &PostHandler{service: service}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePostRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	authorID := primitive.NewObjectID()

	post, err := h.service.CreatePost(req, authorID)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	response := dto.PostResponse{
		ID:           post.ID.Hex(),
		AuthorID:     post.AuthorID.Hex(),
		AuthorName:   post.AuthorName,
		Title:        post.Title,
		Content:      post.Content,
		Description:  post.Description,
		Category:     string(post.Category),
		MediaCount:   post.MediaCount,
		LikeCount:    post.LikeCount,
		CommentCount: post.CommentCount,
		CreatedAt:    post.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	for _, media := range post.Media {
		response.Media = append(response.Media, dto.MediaItemResponse{
			URL:      media.URL,
			Type:     media.Type,
			Caption:  media.Caption,
			Position: media.Position,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	posts, err := h.service.GetFeed(limit)
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	var responses []dto.PostResponse
	for _, post := range posts {
		response := dto.PostResponse{
			ID:           post.ID.Hex(),
			AuthorID:     post.AuthorID.Hex(),
			AuthorName:   post.AuthorName,
			Title:        post.Title,
			Content:      post.Content,
			Description:  post.Description,
			Category:     string(post.Category),
			MediaCount:   post.MediaCount,
			LikeCount:    post.LikeCount,
			CommentCount: post.CommentCount,
			CreatedAt:    post.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}

		for _, media := range post.Media {
			response.Media = append(response.Media, dto.MediaItemResponse{
				URL:      media.URL,
				Type:     media.Type,
				Caption:  media.Caption,
				Position: media.Position,
			})
		}

		responses = append(responses, response)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (h *PostHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	h.GetPosts(w, r)
}

func (h *PostHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")

	postID, err := primitive.ObjectIDFromHex(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID := primitive.NewObjectID()

	err = h.service.LikePost(postID, userID)
	if err != nil {
		http.Error(w, "Failed to like post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Post liked successfully",
		"post_id": postIDStr,
	})
}
