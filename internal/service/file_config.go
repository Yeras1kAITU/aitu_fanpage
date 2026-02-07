package service

import (
	"os"
	"path/filepath"
)

type FileServiceConfig struct {
	UploadDir        string
	TempDir          string
	MaxFileSize      int64
	AllowedTypes     []string
	AllowedMIMETypes []string
	MaxFilesPerPost  int
	ServeURL         string
	ImageSizes       []ImageSizeConfig
	EnableThumbnails bool
}

type ImageSizeConfig struct {
	Name   string
	Width  int
	Height int
}

func DefaultFileServiceConfig() FileServiceConfig {
	uploadDir := getEnv("UPLOAD_DIR", "./uploads")

	return FileServiceConfig{
		UploadDir:        uploadDir,
		TempDir:          filepath.Join(uploadDir, "temp"),
		MaxFileSize:      10 * 1024 * 1024,
		AllowedTypes:     []string{".jpg", ".jpeg", ".png", ".gif", ".mp4", ".mov", ".webm", ".pdf"},
		AllowedMIMETypes: []string{"image/jpeg", "image/png", "image/gif", "video/mp4", "video/quicktime", "video/webm", "application/pdf"},
		MaxFilesPerPost:  10,
		ServeURL:         "/uploads",
		ImageSizes: []ImageSizeConfig{
			{Name: "thumb", Width: 150, Height: 150},
			{Name: "small", Width: 400, Height: 400},
			{Name: "medium", Width: 800, Height: 800},
			{Name: "large", Width: 1200, Height: 1200},
		},
		EnableThumbnails: true,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
