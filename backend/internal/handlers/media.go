package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/dto"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/middleware"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/service"
)

type MediaHandler struct {
	fileService *service.FileService
	postService *service.PostService
	userService *service.UserService
}

func NewMediaHandler(fileService *service.FileService, postService *service.PostService) *MediaHandler {
	return &MediaHandler{
		fileService: fileService,
		postService: postService,
	}
}

func (h *MediaHandler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Multipart form
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	uploadedFiles, err := h.fileService.UploadMultipleFiles(files)
	if err != nil {
		http.Error(w, "Failed to upload files: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []dto.MediaResponse
	for _, uploadedFile := range uploadedFiles {
		responses = append(responses, dto.MediaResponse{
			ID:           uploadedFile.ID,
			OriginalName: uploadedFile.OriginalName,
			FileName:     uploadedFile.FileName,
			FileSize:     uploadedFile.FileSize,
			FileType:     uploadedFile.FileType,
			MediaType:    string(uploadedFile.MediaType),
			URL:          uploadedFile.URL,
			ThumbnailURL: uploadedFile.ThumbnailURL,
			Dimensions:   uploadedFile.Dimensions,
			Checksum:     uploadedFile.Checksum,
			CreatedAt:    uploadedFile.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responses)
}

func (h *MediaHandler) DeleteMedia(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is admin
	user, err := h.postService.GetUserByID(userID)
	if err != nil || !user.IsAdmin() {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	mediaURL := chi.URLParam(r, "url")
	if mediaURL == "" {
		http.Error(w, "Media URL is required", http.StatusBadRequest)
		return
	}

	mediaURL, err = url.QueryUnescape(mediaURL)
	if err != nil {
		http.Error(w, "Invalid media URL", http.StatusBadRequest)
		return
	}

	err = h.fileService.DeleteFile(mediaURL)
	if err != nil {
		http.Error(w, "Failed to delete media: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Media deleted successfully",
	})
}

func (h *MediaHandler) GetMediaInfo(w http.ResponseWriter, r *http.Request) {
	mediaURL := chi.URLParam(r, "url")
	if mediaURL == "" {
		http.Error(w, "Media URL is required", http.StatusBadRequest)
		return
	}

	mediaURL, err := url.QueryUnescape(mediaURL)
	if err != nil {
		http.Error(w, "Invalid media URL", http.StatusBadRequest)
		return
	}

	fileInfo, err := h.fileService.GetFileInfo(mediaURL)
	if err != nil {
		http.Error(w, "Failed to get media info: "+err.Error(), http.StatusNotFound)
		return
	}

	response := dto.MediaResponse{
		ID:        fileInfo.ID,
		FileName:  fileInfo.FileName,
		FileSize:  fileInfo.FileSize,
		URL:       fileInfo.URL,
		CreatedAt: fileInfo.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *MediaHandler) ServeMedia(w http.ResponseWriter, r *http.Request) {
	mediaPath := chi.URLParam(r, "*")
	if mediaPath == "" {
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	// Security check
	if strings.Contains(mediaPath, "..") {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(h.fileService.GetUploadDir(), mediaPath)

	// Check if file exists
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year cache
	http.ServeFile(w, r, fullPath)
}

func (h *MediaHandler) GetUploadDir() string {
	return "./uploads"
}
