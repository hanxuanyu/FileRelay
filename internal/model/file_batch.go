package model

import (
	"time"

	"gorm.io/gorm"
)

type FileBatch struct {
	ID            string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	PickupCode    string         `gorm:"uniqueIndex;not null" json:"pickup_code"`
	Remark        string         `json:"remark"`
	ExpireType    string         `json:"expire_type"` // time / download / permanent
	ExpireAt      *time.Time     `json:"expire_at"`
	MaxDownloads  int            `json:"max_downloads"`
	DownloadCount int            `gorm:"default:0" json:"download_count"`
	Status        string         `gorm:"default:'active'" json:"status"` // active / expired / deleted
	Type          string         `gorm:"default:'file'" json:"type"`     // file / text
	Content       string         `json:"content,omitempty"`
	FileItems     []FileItem     `gorm:"foreignKey:BatchID" json:"file_items,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
