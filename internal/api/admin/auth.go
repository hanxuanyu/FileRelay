package admin

import (
	"FileRelay/internal/auth"
	"FileRelay/internal/config"
	"FileRelay/internal/model"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

type LoginRequest struct {
	Password string `json:"password" binding:"required" example:"admin"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Login 管理员登录
// @Summary      管理员登录
// @Description  通过密码换取 JWT Token
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "登录请求"
// @Success      200  {object}  model.Response{data=LoginResponse}
// @Failure      401  {object}  model.Response
// @Router       /api/admin/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, "Invalid request"))
		return
	}

	passwordHash := config.GlobalConfig.Security.AdminPasswordHash
	if passwordHash == "" {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "Admin password hash not configured"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		slog.Warn("Failed admin login attempt", "ip", c.ClientIP())
		c.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "Incorrect password"))
		return
	}

	// 使用固定 ID 1 代表管理员（因为不再有数据库记录）
	token, err := auth.GenerateToken(1)
	if err != nil {
		slog.Error("Failed to generate admin token", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "Failed to generate token"))
		return
	}

	slog.Info("Admin logged in", "ip", c.ClientIP())
	c.JSON(http.StatusOK, model.SuccessResponse(LoginResponse{
		Token: token,
	}))
}
