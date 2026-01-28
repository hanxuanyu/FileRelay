package service

import (
	"FileRelay/internal/bootstrap"
	"FileRelay/internal/model"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenService struct {
	db *gorm.DB
}

func NewTokenService() *TokenService {
	return &TokenService{db: bootstrap.DB}
}

func (s *TokenService) CreateToken(name string, scope string, expireAt *time.Time) (string, *model.APIToken, error) {
	rawToken := uuid.New().String()
	hash := s.hashToken(rawToken)

	token := &model.APIToken{
		Name:      name,
		TokenHash: hash,
		Scope:     scope,
		ExpireAt:  expireAt,
	}

	if err := s.db.Create(token).Error; err != nil {
		return "", nil, err
	}

	return rawToken, token, nil
}

func (s *TokenService) ValidateToken(rawToken string, requiredScope string) (*model.APIToken, error) {
	hash := s.hashToken(rawToken)
	var token model.APIToken
	if err := s.db.Where("token_hash = ? AND revoked = ?", hash, false).First(&token).Error; err != nil {
		return nil, errors.New("invalid token")
	}

	if token.ExpireAt != nil && time.Now().After(*token.ExpireAt) {
		return nil, errors.New("token expired")
	}

	// 检查 Scope (简单包含判断)
	// 在实际应用中可以实现更复杂的逻辑
	if requiredScope != "" && !s.checkScope(token.Scope, requiredScope) {
		return nil, errors.New("insufficient scope")
	}

	// 更新最后使用时间
	now := time.Now()
	s.db.Model(&token).Update("last_used_at", &now)

	return &token, nil
}

func (s *TokenService) hashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *TokenService) checkScope(tokenScope, requiredScope string) bool {
	if requiredScope == "" {
		return true
	}
	scopes := strings.Split(tokenScope, ",")
	for _, s := range scopes {
		if strings.TrimSpace(s) == requiredScope {
			return true
		}
	}
	return false
}
