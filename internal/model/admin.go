package model

import (
	"time"
)

// AdminSession 管理员会话信息 (不再存库，仅用于 JWT 或 API 交互)
type AdminSession struct {
	ID        uint       `json:"id"`
	LastLogin *time.Time `json:"last_login"`
}
