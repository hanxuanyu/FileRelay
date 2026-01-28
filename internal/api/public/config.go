package public

import (
	"FileRelay/internal/config"
	"FileRelay/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct{}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{}
}

// PublicConfig 公开配置结构
type PublicConfig struct {
	Site     config.SiteConfig    `json:"site"`
	Security PublicSecurityConfig `json:"security"`
	Upload   config.UploadConfig  `json:"upload"`
	APIToken PublicAPITokenConfig `json:"api_token"`
	Storage  PublicStorageConfig  `json:"storage"`
}

type PublicSecurityConfig struct {
	PickupCodeLength int `json:"pickup_code_length"`
}

type PublicAPITokenConfig struct {
	Enabled bool `json:"enabled"`
}

type PublicStorageConfig struct {
	Type string `json:"type"`
}

// GetPublicConfig 获取非敏感配置
// @Summary 获取公共配置
// @Description 获取前端展示所需的非敏感配置数据
// @Tags         Public
// @Produce      json
// @Success 200 {object} model.Response{data=PublicConfig}
// @Router /api/config [get]
func (h *ConfigHandler) GetPublicConfig(c *gin.Context) {
	pub := PublicConfig{
		Site: config.GlobalConfig.Site,
		Security: PublicSecurityConfig{
			PickupCodeLength: config.GlobalConfig.Security.PickupCodeLength,
		},
		Upload: config.GlobalConfig.Upload,
		APIToken: PublicAPITokenConfig{
			Enabled: config.GlobalConfig.APIToken.Enabled,
		},
		Storage: PublicStorageConfig{
			Type: config.GlobalConfig.Storage.Type,
		},
	}

	c.JSON(http.StatusOK, model.SuccessResponse(pub))
}
