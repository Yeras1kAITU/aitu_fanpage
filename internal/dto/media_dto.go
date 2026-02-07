package dto

type MediaUploadRequest struct {
	URL     string `json:"url"`
	Type    string `json:"type"`
	Caption string `json:"caption,omitempty"`
}

type MediaResponse struct {
	ID           string `json:"id"`
	OriginalName string `json:"original_name,omitempty"`
	FileName     string `json:"file_name"`
	FileSize     int64  `json:"file_size"`
	FileType     string `json:"file_type"`
	MediaType    string `json:"media_type"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	Dimensions   string `json:"dimensions,omitempty"`
	Duration     string `json:"duration,omitempty"`
	Checksum     string `json:"checksum,omitempty"`
	CreatedAt    string `json:"created_at"`
}

type UpdatePostMediaRequest struct {
	Action     string               `json:"action"` // add, remove, update
	Position   int                  `json:"position,omitempty"`
	Caption    string               `json:"caption,omitempty"`
	MediaItems []MediaUploadRequest `json:"media_items,omitempty"`
}

type MediaStats struct {
	TotalFiles     int                       `json:"total_files"`
	TotalSize      string                    `json:"total_size"`
	ByType         map[string]MediaTypeStats `json:"by_type"`
	LastCleanup    string                    `json:"last_cleanup"`
	StorageQuota   string                    `json:"storage_quota"`
	StorageUsed    string                    `json:"storage_used"`
	StoragePercent float64                   `json:"storage_percent"`
}

type MediaTypeStats struct {
	Count int    `json:"count"`
	Size  string `json:"size"`
}

type MediaItemResponse struct {
	URL          string `json:"url"`
	Type         string `json:"type"`
	Caption      string `json:"caption,omitempty"`
	Position     int    `json:"position"`
	FileSize     int64  `json:"file_size,omitempty"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	Dimensions   string `json:"dimensions,omitempty"`
}
