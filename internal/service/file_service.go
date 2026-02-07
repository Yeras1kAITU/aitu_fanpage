package service

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/h2non/filetype"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/config"
)

type MediaType string

const (
	MediaTypeImage    MediaType = "image"
	MediaTypeVideo    MediaType = "video"
	MediaTypeDocument MediaType = "document"
	MediaTypeAudio    MediaType = "audio"
	MediaTypeOther    MediaType = "other"
)

type UploadedFile struct {
	ID           string
	OriginalName string
	FileName     string
	FileSize     int64
	FileType     string
	MediaType    MediaType
	URL          string
	ThumbnailURL string
	Dimensions   string // for images: "widthxheight"
	Duration     string // for videos: duration in seconds
	Checksum     string
	CreatedAt    string
}

type FileService struct {
	uploadDir        string
	tempDir          string
	maxFileSize      int64
	allowedTypes     []string
	allowedMIMEs     []string
	maxFilesPerPost  int
	serveURL         string
	imageSizes       []config.ImageSize
	enableThumbnails bool
	mu               sync.RWMutex
}

func NewFileService(cfg config.UploadConfig) *FileService {
	uploadDir := cfg.UploadDir
	if uploadDir == "./uploads" {
		uploadDir = "/data/uploads"
	}

	tempDir := cfg.TempDir
	if tempDir == "./temp" {
		tempDir = "/data/temp"
	}

	os.MkdirAll(uploadDir, 0755)
	os.MkdirAll(tempDir, 0755)

	// Create subdirectories
	for _, mediaType := range []string{"images", "videos", "documents", "temp"} {
		os.MkdirAll(filepath.Join(uploadDir, mediaType), 0755)
	}

	return &FileService{
		uploadDir:        uploadDir,
		tempDir:          tempDir,
		maxFileSize:      cfg.MaxFileSize,
		allowedTypes:     cfg.AllowedTypes,
		allowedMIMEs:     cfg.AllowedMIMETypes,
		maxFilesPerPost:  cfg.MaxFilesPerPost,
		serveURL:         cfg.ServeURL,
		imageSizes:       cfg.ImageSizes,
		enableThumbnails: cfg.EnableThumbnails,
	}
}

func (fs *FileService) GetUploadDir() string {
	return fs.uploadDir
}

func (fs *FileService) UploadFile(fileHeader *multipart.FileHeader) (*UploadedFile, error) {
	if err := fs.validateFile(fileHeader); err != nil {
		return nil, err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	file.Seek(0, 0)

	// Generate checksum
	checksum := fs.calculateChecksum(fileBytes)

	// Determine file type
	fileType, mediaType := fs.detectFileType(fileBytes, fileHeader.Filename)
	if fileType == "" {
		fileType = strings.ToLower(filepath.Ext(fileHeader.Filename))
	}

	// Generate unique filename
	fileExt := filepath.Ext(fileHeader.Filename)
	if fileExt == "" && fileType != "" {
		fileExt = fileType
	}

	uniqueID := uuid.New().String()
	fileName := uniqueID + fileExt

	storageDir := fs.getStorageDir(mediaType)
	fullPath := filepath.Join(storageDir, fileName)

	// Save original file
	if err := fs.saveFile(fileBytes, fullPath); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	thumbnailURL := ""
	if mediaType == MediaTypeImage && fs.enableThumbnails {
		thumbnailURL, err = fs.createThumbnail(fileBytes, fileName, fileType)
		if err != nil {
			fmt.Printf("Failed to create thumbnail: %v\n", err)
			// Continue even if thumbnail creation fails
		}
	}

	dimensions := ""
	if mediaType == MediaTypeImage {
		if img, _, err := image.Decode(bytes.NewReader(fileBytes)); err == nil {
			bounds := img.Bounds()
			dimensions = fmt.Sprintf("%dx%d", bounds.Dx(), bounds.Dy())
		}
	}

	uploadedFile := &UploadedFile{
		ID:           uniqueID,
		OriginalName: fileHeader.Filename,
		FileName:     fileName,
		FileSize:     fileHeader.Size,
		FileType:     fileType,
		MediaType:    mediaType,
		URL:          fs.getServeURL(mediaType, fileName),
		ThumbnailURL: thumbnailURL,
		Dimensions:   dimensions,
		Checksum:     checksum,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}

	return uploadedFile, nil
}

func (fs *FileService) UploadMultipleFiles(files []*multipart.FileHeader) ([]*UploadedFile, error) {
	if len(files) > fs.maxFilesPerPost {
		return nil, fmt.Errorf("maximum %d files allowed per post", fs.maxFilesPerPost)
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(files))
	results := make(chan *UploadedFile, len(files))

	for _, file := range files {
		wg.Add(1)
		go func(fh *multipart.FileHeader) {
			defer wg.Done()

			uploadedFile, err := fs.UploadFile(fh)
			if err != nil {
				errors <- err
				return
			}
			results <- uploadedFile
		}(file)
	}

	wg.Wait()
	close(errors)
	close(results)

	var errorList []string
	for err := range errors {
		if err != nil {
			errorList = append(errorList, err.Error())
		}
	}

	if len(errorList) > 0 {
		return nil, fmt.Errorf("upload errors: %s", strings.Join(errorList, "; "))
	}

	var uploadedFiles []*UploadedFile
	for result := range results {
		uploadedFiles = append(uploadedFiles, result)
	}

	return uploadedFiles, nil
}

func (fs *FileService) validateFile(fileHeader *multipart.FileHeader) error {
	if fileHeader.Size > fs.maxFileSize {
		return fmt.Errorf("file too large: %d bytes (max: %d bytes)", fileHeader.Size, fs.maxFileSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	extAllowed := false
	for _, allowedExt := range fs.allowedTypes {
		if strings.ToLower(allowedExt) == ext {
			extAllowed = true
			break
		}
	}

	if !extAllowed && len(fs.allowedTypes) > 0 {
		return fmt.Errorf("file type not allowed: %s", ext)
	}

	return nil
}

func (fs *FileService) detectFileType(fileBytes []byte, filename string) (string, MediaType) {
	kind, err := filetype.Match(fileBytes)
	if err == nil && kind != filetype.Unknown {
		mediaType := fs.mimeToMediaType(kind.MIME.Value)
		if mediaType != MediaTypeOther {
			return kind.Extension, mediaType
		}
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return "", MediaTypeOther
	}

	if strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}

	mediaType := fs.extToMediaType(ext)
	return ext, mediaType
}

func (fs *FileService) mimeToMediaType(mime string) MediaType {
	switch {
	case strings.HasPrefix(mime, "image/"):
		return MediaTypeImage
	case strings.HasPrefix(mime, "video/"):
		return MediaTypeVideo
	case strings.HasPrefix(mime, "audio/"):
		return MediaTypeAudio
	case mime == "application/pdf" ||
		strings.Contains(mime, "document") ||
		strings.Contains(mime, "msword") ||
		strings.Contains(mime, "openxmlformats"):
		return MediaTypeDocument
	default:
		return MediaTypeOther
	}
}

func (fs *FileService) extToMediaType(ext string) MediaType {
	imageExts := []string{"jpg", "jpeg", "png", "gif", "bmp", "webp", "svg"}
	videoExts := []string{"mp4", "mov", "avi", "mkv", "webm", "flv", "wmv"}
	audioExts := []string{"mp3", "wav", "ogg", "flac", "m4a"}
	documentExts := []string{"pdf", "doc", "docx", "txt", "rtf", "odt", "xls", "xlsx", "ppt", "pptx"}

	for _, imgExt := range imageExts {
		if ext == imgExt {
			return MediaTypeImage
		}
	}

	for _, vidExt := range videoExts {
		if ext == vidExt {
			return MediaTypeVideo
		}
	}

	for _, audExt := range audioExts {
		if ext == audExt {
			return MediaTypeAudio
		}
	}

	for _, docExt := range documentExts {
		if ext == docExt {
			return MediaTypeDocument
		}
	}

	return MediaTypeOther
}

func (fs *FileService) calculateChecksum(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

func (fs *FileService) getStorageDir(mediaType MediaType) string {
	switch mediaType {
	case MediaTypeImage:
		return filepath.Join(fs.uploadDir, "images")
	case MediaTypeVideo:
		return filepath.Join(fs.uploadDir, "videos")
	case MediaTypeDocument:
		return filepath.Join(fs.uploadDir, "documents")
	case MediaTypeAudio:
		return filepath.Join(fs.uploadDir, "audio")
	default:
		return filepath.Join(fs.uploadDir, "other")
	}
}

func (fs *FileService) getServeURL(mediaType MediaType, fileName string) string {
	mediaTypeDir := strings.ToLower(string(mediaType)) + "s"
	return fmt.Sprintf("%s/%s/%s", fs.serveURL, mediaTypeDir, fileName)
}

func (fs *FileService) saveFile(data []byte, path string) error {
	return os.WriteFile(path, data, 0644)
}

func (fs *FileService) createThumbnail(data []byte, fileName, fileType string) (string, error) {
	// Decode image
	var img image.Image
	var err error

	switch strings.ToLower(fileType) {
	case "jpg", "jpeg":
		img, err = jpeg.Decode(bytes.NewReader(data))
	case "png":
		img, err = png.Decode(bytes.NewReader(data))
	default:
		// For other image types, try generic decode
		img, _, err = image.Decode(bytes.NewReader(data))
	}

	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v", err)
	}

	// Create thumbnails for each size
	for _, size := range fs.imageSizes {
		// Resize image
		resized := imaging.Resize(img, size.Width, size.Height, imaging.Lanczos)

		// Create thumbnail filename
		ext := filepath.Ext(fileName)
		nameWithoutExt := strings.TrimSuffix(fileName, ext)
		thumbFileName := fmt.Sprintf("%s_%s%s", nameWithoutExt, size.Name, ext)
		thumbPath := filepath.Join(fs.uploadDir, "images", "thumbnails", thumbFileName)

		os.MkdirAll(filepath.Dir(thumbPath), 0755)

		// Save
		switch strings.ToLower(fileType) {
		case "jpg", "jpeg":
			err = imaging.Save(resized, thumbPath, imaging.JPEGQuality(85))
		case "png":
			err = imaging.Save(resized, thumbPath, imaging.PNGCompressionLevel(png.DefaultCompression))
		default:
			err = imaging.Save(resized, thumbPath)
		}

		if err != nil {
			return "", fmt.Errorf("failed to save thumbnail: %v", err)
		}

		if size.Name == "thumb" {
			return fmt.Sprintf("%s/images/thumbnails/%s", fs.serveURL, thumbFileName), nil
		}
	}

	return "", nil
}

func (fs *FileService) DeleteFile(fileURL string) error {
	// Extract filename from URL
	parts := strings.Split(fileURL, "/")
	if len(parts) == 0 {
		return fmt.Errorf("invalid file URL")
	}

	fileName := parts[len(parts)-1]
	mediaTypeDir := parts[len(parts)-2]

	fullPath := filepath.Join(fs.uploadDir, mediaTypeDir, fileName)

	// Delete main file
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	// Delete thumbnails if it's an image
	if mediaTypeDir == "images" {
		// Delete all thumbnail sizes
		for _, size := range fs.imageSizes {
			ext := filepath.Ext(fileName)
			nameWithoutExt := strings.TrimSuffix(fileName, ext)
			thumbFileName := fmt.Sprintf("%s_%s%s", nameWithoutExt, size.Name, ext)
			thumbPath := filepath.Join(fs.uploadDir, "images", "thumbnails", thumbFileName)

			os.Remove(thumbPath)
		}
	}

	return nil
}

func (fs *FileService) GetFileInfo(fileURL string) (*UploadedFile, error) {
	parts := strings.Split(fileURL, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid file URL")
	}

	fileName := parts[len(parts)-1]
	mediaTypeDir := parts[len(parts)-2]
	fullPath := filepath.Join(fs.uploadDir, mediaTypeDir, fileName)

	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %v", err)
	}

	// Extract ID from filename
	var fileID string
	if strings.Contains(fileName, "_") {
		parts := strings.Split(fileName, "_")
		if len(parts) > 0 {
			ext := filepath.Ext(parts[0])
			fileID = strings.TrimSuffix(parts[0], ext)
		}
	}

	if fileID == "" {
		// Use filename without extension as ID
		ext := filepath.Ext(fileName)
		fileID = strings.TrimSuffix(fileName, ext)
	}

	return &UploadedFile{
		ID:        fileID,
		FileName:  fileName,
		FileSize:  fileInfo.Size(),
		URL:       fileURL,
		CreatedAt: fileInfo.ModTime().Format(time.RFC3339),
	}, nil
}

func (fs *FileService) CleanupTempFiles(maxAge time.Duration) error {
	files, err := os.ReadDir(fs.tempDir)
	if err != nil {
		return err
	}

	cutoff := time.Now().Add(-maxAge)

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			fullPath := filepath.Join(fs.tempDir, file.Name())
			os.Remove(fullPath)
		}
	}

	return nil
}
