package middleware

import (
	"FileRelay/internal/auth"
	"FileRelay/internal/config"
	"FileRelay/internal/model"
	"FileRelay/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	tokenService := service.NewTokenService()
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "Authorization header required"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "Invalid authorization format"))
			c.Abort()
			return
		}

		tokenStr := parts[1]

		// 1. 尝试解析为管理员 JWT
		claims, err := auth.ParseToken(tokenStr)
		if err == nil {
			c.Set("admin_id", claims.AdminID)
			c.Next()
			return
		}

		// 2. 尝试解析为 API Token (如果配置允许)
		if config.GlobalConfig.APIToken.Enabled && config.GlobalConfig.APIToken.AllowAdminAPI {
			token, err := tokenService.ValidateToken(tokenStr, model.ScopeAdmin)
			if err == nil {
				c.Set("token_id", token.ID)
				c.Set("token_scope", token.Scope)
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "Invalid or expired token"))
		c.Abort()
	}
}

func APITokenAuth(requiredScope string, optional bool) gin.HandlerFunc {
	tokenService := service.NewTokenService()
	return func(c *gin.Context) {
		handleAPITokenAuth(c, tokenService, requiredScope, optional)
	}
}

func UploadAuth() gin.HandlerFunc {
	tokenService := service.NewTokenService()
	return func(c *gin.Context) {
		// 动态获取配置
		optional := !config.GlobalConfig.Upload.RequireToken
		handleAPITokenAuth(c, tokenService, model.ScopeUpload, optional)
	}
}

func handleAPITokenAuth(c *gin.Context, tokenService *service.TokenService, requiredScope string, optional bool) {
	// 如果是可选的，直接跳过校验，满足“未打开对应的开关时不需校验”的需求
	if optional {
		c.Next()
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "Authorization header required"))
		c.Abort()
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "Invalid authorization format"))
		c.Abort()
		return
	}

	tokenStr := parts[1]

	// 1. 尝试解析为管理员 JWT
	if claims, err := auth.ParseToken(tokenStr); err == nil {
		c.Set("admin_id", claims.AdminID)
		c.Next()
		return
	}

	if !config.GlobalConfig.APIToken.Enabled {
		c.JSON(http.StatusForbidden, model.ErrorResponse(model.CodeForbidden, "API Token is disabled"))
		c.Abort()
		return
	}

	token, err := tokenService.ValidateToken(tokenStr, requiredScope)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, err.Error()))
		c.Abort()
		return
	}

	c.Set("token_id", token.ID)
	c.Set("token_scope", token.Scope)
	c.Next()
}
