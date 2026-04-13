package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Upload   UploadConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	URI  string
	Name string
}

type JWTConfig struct {
	SecretKey string
}

type UploadConfig struct {
	UploadDir        string
	TempDir          string
	MaxFileSize      int64
	AllowedTypes     []string
	AllowedMIMETypes []string
	MaxFilesPerPost  int
	ServeURL         string
	ImageSizes       []ImageSize
	EnableThumbnails bool
}

type ImageSize struct {
	Name   string
	Width  int
	Height int
}

func Load() *Config {
	port := getEnv("PORT", "8080")
	env := getEnv("ENVIRONMENT", "development")

	return &Config{
		Server: ServerConfig{
			Port: port,
			Env:  env,
		},
		Database: DatabaseConfig{
			URI:  getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Name: getEnv("DATABASE_NAME", "aitu_fanpage"),
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		},
		Upload: UploadConfig{
			UploadDir:        getEnv("UPLOAD_DIR", "./uploads"),
			TempDir:          getEnv("TEMP_DIR", "./temp"),
			MaxFileSize:      parseInt64(getEnv("MAX_FILE_SIZE", "10485760")), // 10MB
			AllowedTypes:     []string{".jpg", ".jpeg", ".png", ".gif", ".mp4", ".mov", ".webm", ".pdf"},
			AllowedMIMETypes: []string{"image/jpeg", "image/png", "image/gif", "video/mp4", "video/quicktime", "video/webm", "application/pdf"},
			MaxFilesPerPost:  parseInt(getEnv("MAX_FILES_PER_POST", "10")),
			ServeURL:         getEnv("SERVE_URL", "/uploads"),
			ImageSizes: []ImageSize{
				{Name: "thumb", Width: 150, Height: 150},
				{Name: "small", Width: 400, Height: 400},
				{Name: "medium", Width: 800, Height: 800},
				{Name: "large", Width: 1200, Height: 1200},
			},
			EnableThumbnails: parseBool(getEnv("ENABLE_THUMBNAILS", "true")),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func parseInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func parseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return b
}
