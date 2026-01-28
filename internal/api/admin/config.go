package admin

import (
	"FileRelay/internal/bootstrap"
	"FileRelay/internal/config"
	"FileRelay/internal/model"
	"FileRelay/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ConfigHandler struct{}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{}
}

// GetConfig 获取当前完整配置
// @Summary 获取完整配置
// @Description 获取系统的完整配置文件内容（仅管理员）
// @Tags         Admin
// @Security AdminAuth
// @Produce json
// @Success 200 {object} model.Response{data=config.Config}
// @Router /api/admin/config [get]
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, model.SuccessResponse(config.GlobalConfig))
}

// UpdateConfig 更新配置
// @Summary 更新配置
// @Description 更新系统的配置文件内容（仅管理员）
// @Tags         Admin
// @Security AdminAuth
// @Accept json
// @Produce json
// @Param config body config.Config true "新配置内容"
// @Success 200 {object} model.Response{data=config.Config}
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/admin/config [put]
func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	var newConfig config.Config
	if err := c.ShouldBindJSON(&newConfig); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, err.Error()))
		return
	}

	// 简单的校验，防止关键配置被改空
	if newConfig.Database.Path == "" {
		newConfig.Database.Path = config.GlobalConfig.Database.Path
	}
	if newConfig.Site.Port <= 0 || newConfig.Site.Port > 65535 {
		newConfig.Site.Port = 8080
	}

	// 如果传入了明文密码，则重新生成 hash
	if newConfig.Security.AdminPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(newConfig.Security.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "Failed to hash password: "+err.Error()))
			return
		}
		newConfig.Security.AdminPasswordHash = string(hash)
	}

	// 检查取件码长度是否变化
	pickupCodeLengthChanged := newConfig.Security.PickupCodeLength != config.GlobalConfig.Security.PickupCodeLength && newConfig.Security.PickupCodeLength > 0
	// 检查数据库配置是否变化
	dbConfigChanged := newConfig.Database != config.GlobalConfig.Database

	// 如果长度变化，同步更新现有取件码 (在可能切换数据库前，先处理旧库数据)
	if pickupCodeLengthChanged {
		batchService := service.NewBatchService()
		if err := batchService.UpdateAllPickupCodes(newConfig.Security.PickupCodeLength); err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "Failed to update existing pickup codes: "+err.Error()))
			return
		}
	}

	// 更新内存配置
	config.UpdateGlobalConfig(&newConfig)

	// 重新连接数据库并迁移数据（如果配置发生变化）
	if dbConfigChanged {
		if err := bootstrap.ReloadDB(newConfig.Database); err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "Failed to reload database: "+err.Error()))
			return
		}
	}

	// 重新初始化存储（热更新业务逻辑）
	if err := bootstrap.ReloadStorage(); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "Failed to reload storage: "+err.Error()))
		return
	}

	// 保存到文件
	if err := config.SaveConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "Failed to save config: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(config.GlobalConfig))
}
