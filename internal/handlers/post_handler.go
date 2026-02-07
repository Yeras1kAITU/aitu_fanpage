package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/dto"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/middleware"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/models"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/service"
)

type PostHandler struct {
	service     *service.PostService
	fileService *service.FileService
}

type HandlerContainer struct {
	Auth    *AuthHandler
	Post    *PostHandler
	Comment *CommentHandler
	User    *UserHandler
	Admin   *AdminHandler
	Media   *MediaHandler
}

func NewPostHandler(service *service.PostService, fileService *service.FileService) *PostHandler {
	return &PostHandler{
		service:     service,
		fileService: fileService,
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	contentType := r.Header.Get("Content-Type")

	var req dto.CreatePostRequest
	var uploadedFiles []*service.UploadedFile

	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Multipart form
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Parse JSON
		postData := r.FormValue("post")
		if postData == "" {
			http.Error(w, "Post data is required", http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal([]byte(postData), &req); err != nil {
			http.Error(w, "Invalid post data: "+err.Error(), http.StatusBadRequest)
			return
		}

		files := r.MultipartForm.File["files"]
		if len(files) > 0 {
			uploadedFiles, err = h.fileService.UploadMultipleFiles(files)
			if err != nil {
				http.Error(w, "Failed to upload files: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		// JSON without files
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	if req.Category == "" {
		http.Error(w, "Category is required", http.StatusBadRequest)
		return
	}

	// Create post
	post, err := h.service.CreatePost(req, userID, uploadedFiles)
	if err != nil {
		http.Error(w, "Failed to create post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := h.mapPostToResponse(post)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== GetPosts called ===")
	fmt.Printf("URL: %s\n", r.URL.String())
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("Query: %v\n", r.URL.Query())

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	fmt.Printf("Limit: %s, Offset: %s\n", limitStr, offsetStr)

	category := r.URL.Query().Get("category")
	authorID := r.URL.Query().Get("author_id")

	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	var posts []*models.Post
	var err error

	if category != "" {
		posts, err = h.service.GetPostsByCategory(category, limit)
	} else if authorID != "" {
		authorObjID, err := primitive.ObjectIDFromHex(authorID)
		if err != nil {
			http.Error(w, "Invalid author ID", http.StatusBadRequest)
			return
		}
		posts, err = h.service.GetPostsByAuthor(authorObjID, limit)
	} else {
		posts, err = h.service.GetPosts(limit, offset)
	}

	if err != nil {
		http.Error(w, "Failed to get posts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []dto.PostResponse
	for _, post := range posts {
		responses = append(responses, h.mapPostToResponse(post))
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := primitive.ObjectIDFromHex(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.service.GetPostByID(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	response := h.mapPostToResponse(post)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	category := r.URL.Query().Get("category")

	limit := 20
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

	posts, err := h.service.GetFeed(userID, category, limit, offset)
	if err != nil {
		http.Error(w, "Failed to get feed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []dto.PostResponse
	for _, post := range posts {
		responses = append(responses, h.mapPostToResponse(post))
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) LikePost(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.LikePost(postID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "already liked") {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to like post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Post liked successfully",
		"post_id": postIDStr,
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) UnlikePost(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.UnlikePost(postID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not liked") {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to unlike post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Post unliked successfully",
		"post_id": postIDStr,
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
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

	var req dto.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	post, err := h.service.UpdatePost(postID, userID, req)
	if err != nil {
		if strings.Contains(err.Error(), "not authorized") {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to update post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := h.mapPostToResponse(post)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.DeletePost(postID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not authorized") {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to delete post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Post deleted successfully",
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetPostLikes(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := primitive.ObjectIDFromHex(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	likes, err := h.service.GetPostLikes(postID)
	if err != nil {
		http.Error(w, "Failed to get likes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(likes); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) SearchPosts(w http.ResponseWriter, r *http.Request) {
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

	posts, err := h.service.SearchPosts(query, limit)
	if err != nil {
		http.Error(w, "Failed to search posts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []dto.PostResponse
	for _, post := range posts {
		responses = append(responses, h.mapPostToResponse(post))
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetPinnedPosts(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 5
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	posts, err := h.service.GetPinnedPosts(limit)
	if err != nil {
		http.Error(w, "Failed to get pinned posts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []dto.PostResponse
	for _, post := range posts {
		responses = append(responses, h.mapPostToResponse(post))
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetFeaturedPosts(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	posts, err := h.service.GetFeaturedPosts(limit)
	if err != nil {
		http.Error(w, "Failed to get featured posts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []dto.PostResponse
	for _, post := range posts {
		responses = append(responses, h.mapPostToResponse(post))
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetPopularPosts(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	daysStr := r.URL.Query().Get("days")

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	days := 7
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	posts, err := h.service.GetPopularPosts(limit, days)
	if err != nil {
		http.Error(w, "Failed to get popular posts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []dto.PostResponse
	for _, post := range posts {
		responses = append(responses, h.mapPostToResponse(post))
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetPostsByTags(w http.ResponseWriter, r *http.Request) {
	tagsParam := r.URL.Query().Get("tags")
	if tagsParam == "" {
		http.Error(w, "Tags parameter is required", http.StatusBadRequest)
		return
	}

	tags := strings.Split(tagsParam, ",")

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	posts, err := h.service.GetPostsByTags(tags, limit)
	if err != nil {
		http.Error(w, "Failed to get posts by tags: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []dto.PostResponse
	for _, post := range posts {
		responses = append(responses, h.mapPostToResponse(post))
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetCategoriesStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetCategoriesStats()
	if err != nil {
		http.Error(w, "Failed to get categories stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	total := 0
	for _, count := range stats {
		total += count
	}

	response := dto.CategoriesStatsResponse{
		Categories: stats,
		TotalPosts: total,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) PinPost(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.PinPost(postID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not authorized") {
			status = http.StatusForbidden
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Post pinned successfully",
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) UnpinPost(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.UnpinPost(postID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not authorized") {
			status = http.StatusForbidden
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Post unpinned successfully",
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) FeaturePost(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.FeaturePost(postID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not authorized") {
			status = http.StatusForbidden
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Post featured successfully",
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) UnfeaturePost(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.UnfeaturePost(postID, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not authorized") {
			status = http.StatusForbidden
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Post unfeatured successfully",
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PostHandler) mapPostToResponse(post *models.Post) dto.PostResponse {
	response := dto.PostResponse{
		ID:              post.ID.Hex(),
		AuthorID:        post.AuthorID.Hex(),
		AuthorName:      post.AuthorName,
		Title:           post.Title,
		Content:         post.Content,
		Description:     post.Description,
		Category:        string(post.Category),
		Tags:            post.Tags,
		MediaCount:      post.MediaCount,
		LikeCount:       post.LikeCount,
		CommentCount:    post.CommentCount,
		ViewCount:       post.ViewCount,
		IsFeatured:      post.IsFeatured,
		IsPinned:        post.IsPinned,
		CreatedAt:       post.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:       post.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		PopularityScore: post.PopularityScore,
	}

	for _, media := range post.Media {
		response.Media = append(response.Media, dto.MediaItemResponse{
			URL:          media.URL,
			Type:         media.Type,
			Caption:      media.Caption,
			Position:     media.Position,
			FileSize:     media.FileSize,
			ThumbnailURL: media.ThumbnailURL,
			Dimensions:   media.Dimensions,
		})
	}

	return response
}
