package admin

import (
	"FileRelay/internal/bootstrap"
	"FileRelay/internal/model"
	"FileRelay/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	tokenService *service.TokenService
}

func NewTokenHandler() *TokenHandler {
	return &TokenHandler{
		tokenService: service.NewTokenService(),
	}
}

type CreateTokenRequest struct {
	Name     string     `json:"name" binding:"required" example:"Test Token"`
	Scope    string     `json:"scope" example:"upload,pickup" enums:"upload,pickup,admin"`
	ExpireAt *time.Time `json:"expire_at"`
}

type CreateTokenResponse struct {
	Token string          `json:"token"`
	Data  *model.APIToken `json:"data"`
}

// ListTokens 获取 API Token 列表
// @Summary      获取 API Token 列表
// @Description  获取系统中所有 API Token 的详细信息（不包含哈希）
// @Tags         Admin
// @Security     AdminAuth
// @Produce      json
// @Success      200  {object}  model.Response{data=[]model.APIToken}
// @Failure      401  {object}  model.Response
// @Router       /api/admin/api-tokens [get]
func (h *TokenHandler) ListTokens(c *gin.Context) {
	var tokens []model.APIToken
	if err := bootstrap.DB.Find(&tokens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(tokens))
}

// CreateToken 创建 API Token
// @Summary      创建 API Token
// @Description  创建一个新的 API Token，返回原始 Token（仅显示一次）
// @Tags         Admin
// @Security     AdminAuth
// @Accept       json
// @Produce      json
// @Param        request body CreateTokenRequest true "Token 信息"
// @Success      201  {object}  model.Response{data=CreateTokenResponse}
// @Failure      400  {object}  model.Response
// @Router       /api/admin/api-tokens [post]
func (h *TokenHandler) CreateToken(c *gin.Context) {
	var input CreateTokenRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, err.Error()))
		return
	}

	rawToken, token, err := h.tokenService.CreateToken(input.Name, input.Scope, input.ExpireAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse(CreateTokenResponse{
		Token: rawToken,
		Data:  token,
	}))
}

// RevokeToken 撤销 API Token
// @Summary      撤销 API Token
// @Description  将 API Token 标记为已撤销，使其失效但保留记录
// @Tags         Admin
// @Security     AdminAuth
// @Param        id   path      int     true  "Token ID"
// @Produce      json
// @Success      200  {object}  model.Response
// @Failure      500  {object}  model.Response
// @Router       /api/admin/api-tokens/{id}/revoke [post]
func (h *TokenHandler) RevokeToken(c *gin.Context) {
	id := c.Param("id")
	if err := bootstrap.DB.Model(&model.APIToken{}).Where("id = ?", id).Update("revoked", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(map[string]interface{}{}))
}

// RecoverToken 恢复 API Token
// @Summary      恢复 API Token
// @Description  将已撤销的 API Token 恢复为有效状态
// @Tags         Admin
// @Security     AdminAuth
// @Param        id   path      int     true  "Token ID"
// @Produce      json
// @Success      200  {object}  model.Response
// @Failure      500  {object}  model.Response
// @Router       /api/admin/api-tokens/{id}/recover [post]
func (h *TokenHandler) RecoverToken(c *gin.Context) {
	id := c.Param("id")
	if err := bootstrap.DB.Model(&model.APIToken{}).Where("id = ?", id).Update("revoked", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(map[string]interface{}{}))
}

// DeleteToken 删除 API Token
// @Summary      删除 API Token
// @Description  根据 ID 永久删除 API Token
// @Tags         Admin
// @Security     AdminAuth
// @Param        id   path      int     true  "Token ID"
// @Produce      json
// @Success      200  {object}  model.Response
// @Failure      500  {object}  model.Response
// @Router       /api/admin/api-tokens/{id} [delete]
func (h *TokenHandler) DeleteToken(c *gin.Context) {
	id := c.Param("id")
	if err := bootstrap.DB.Delete(&model.APIToken{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(map[string]interface{}{}))
}
