package model

import (
	"time"
)

type FileItem struct {
	ID           string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	BatchID      string    `gorm:"index;not null;type:varchar(36)" json:"batch_id"`
	OriginalName string    `json:"original_name"`
	StoragePath  string    `json:"storage_path"`
	Size         int64     `json:"size"`
	MimeType     string    `json:"mime_type"`
	DownloadURL  string    `gorm:"-" json:"download_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
