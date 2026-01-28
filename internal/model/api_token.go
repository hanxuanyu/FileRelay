package model

import (
	"time"
)

const (
	ScopeUpload = "upload" // 上传权限
	ScopePickup = "pickup" // 取件/下载权限
	ScopeAdmin  = "admin"  // 管理权限
)

type APIToken struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	Name       string     `json:"name"`
	TokenHash  string     `gorm:"uniqueIndex;not null" json:"-"`
	Scope      string     `json:"scope"`
	ExpireAt   *time.Time `json:"expire_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
	Revoked    bool       `gorm:"default:false" json:"revoked"`
	CreatedAt  time.Time  `json:"created_at"`
}
