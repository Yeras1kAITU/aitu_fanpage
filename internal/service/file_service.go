package service

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type FileService struct {
	uploadDir    string
	maxFileSize  int64
	allowedTypes []string
}

func NewFileService(uploadDir string, maxFileSize int64, allowedTypes []string) *FileService {
	return &FileService{
		uploadDir:    uploadDir,
		maxFileSize:  maxFileSize,
		allowedTypes: allowedTypes,
	}
}

func (fs *FileService) UploadMultipleFiles(files []*multipart.FileHeader) ([]string, error) {
	var wg sync.WaitGroup
	errors := make(chan error, len(files))
	urls := make([]string, len(files))

	if len(files) > 10 {
		return nil, fmt.Errorf("maximum 10 files allowed")
	}

	for i, file := range files {
		wg.Add(1)

		go func(index int, fh *multipart.FileHeader) {
			defer wg.Done()

			if err := fs.validateFile(fh); err != nil {
				errors <- fmt.Errorf("file %d: %v", index+1, err)
				return
			}

			url, err := fs.uploadFile(fh)
			if err != nil {
				errors <- fmt.Errorf("file %d: %v", index+1, err)
				return
			}

			urls[index] = url
		}(i, file)
	}

	wg.Wait()
	close(errors)

	var errorList []string
	for err := range errors {
		if err != nil {
			errorList = append(errorList, err.Error())
		}
	}

	if len(errorList) > 0 {
		return nil, fmt.Errorf("upload errors: %s", strings.Join(errorList, "; "))
	}

	return urls, nil
}

func (fs *FileService) GenerateThumbnail(imageURL string) (string, error) {
	return imageURL + "?thumb=200x200", nil
}

func (fs *FileService) validateFile(fh *multipart.FileHeader) error {
	if fh.Size > fs.maxFileSize {
		return fmt.Errorf("file too large: %d > %d", fh.Size, fs.maxFileSize)
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	allowed := false
	for _, allowedExt := range fs.allowedTypes {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("file type not allowed: %s", ext)
	}

	return nil
}

func (fs *FileService) uploadFile(fh *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(fh.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	return fmt.Sprintf("/uploads/%s", filename), nil
}

func (fs *FileService) UploadSingleFile(file *multipart.FileHeader) (string, error) {
	if err := fs.validateFile(file); err != nil {
		return "", err
	}

	return fs.uploadFile(file)
}
