package admin

import (
	"FileRelay/internal/bootstrap"
	"FileRelay/internal/model"
	"FileRelay/internal/service"
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BatchHandler struct {
	batchService *service.BatchService
}

func NewBatchHandler() *BatchHandler {
	return &BatchHandler{
		batchService: service.NewBatchService(),
	}
}

type ListBatchesResponse struct {
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Data     []model.FileBatch `json:"data"`
}

type UpdateBatchRequest struct {
	Remark        *string    `json:"remark"`
	ExpireType    *string    `json:"expire_type"`
	ExpireAt      *time.Time `json:"expire_at"`
	MaxDownloads  *int       `json:"max_downloads"`
	DownloadCount *int       `json:"download_count"`
	Status        *string    `json:"status"`
}

// ListBatches 获取批次列表
// @Summary      获取批次列表
// @Description  分页查询所有文件批次，支持按状态过滤和取件码模糊搜索
// @Tags         Admin
// @Security     AdminAuth
// @Param        page         query     int     false  "页码 (默认 1)"
// @Param        page_size    query     int     false  "每页数量 (默认 20)"
// @Param        status       query     string  false  "状态 (active/expired/deleted)"
// @Param        pickup_code  query     string  false  "取件码 (模糊搜索)"
// @Produce      json
// @Success      200  {object}  model.Response{data=ListBatchesResponse}
// @Failure      401  {object}  model.Response
// @Router       /api/admin/batches [get]
func (h *BatchHandler) ListBatches(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	status := c.Query("status")
	pickupCode := c.Query("pickup_code")

	query := bootstrap.DB.Model(&model.FileBatch{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if pickupCode != "" {
		query = query.Where("pickup_code LIKE ?", "%"+pickupCode+"%")
	}

	var total int64
	query.Count(&total)

	var batches []model.FileBatch
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&batches).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(ListBatchesResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Data:     batches,
	}))
}

// GetBatch 获取批次详情
// @Summary      获取批次详情
// @Description  根据批次 ID 获取批次信息及关联的文件列表
// @Tags         Admin
// @Security     AdminAuth
// @Param        batch_id  path  string  true  "批次 ID (UUID)"
// @Produce      json
// @Success      200  {object}  model.Response{data=model.FileBatch}
// @Failure      404  {object}  model.Response
// @Router       /api/admin/batches/{batch_id} [get]
func (h *BatchHandler) GetBatch(c *gin.Context) {
	id := c.Param("batch_id")
	var batch model.FileBatch
	if err := bootstrap.DB.Preload("FileItems").First(&batch, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse(model.CodeNotFound, "batch not found"))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(batch))
}

// UpdateBatch 修改批次信息
// @Summary      修改批次信息
// @Description  允许修改备注、过期策略、最大下载次数、状态等
// @Tags         Admin
// @Security     AdminAuth
// @Accept       json
// @Produce      json
// @Param        batch_id  path  string  true  "批次 ID (UUID)"
// @Param        request body UpdateBatchRequest true "修改内容"
// @Success      200  {object}  model.Response{data=model.FileBatch}
// @Failure      400  {object}  model.Response
// @Router       /api/admin/batches/{batch_id} [put]
func (h *BatchHandler) UpdateBatch(c *gin.Context) {
	id := c.Param("batch_id")
	var batch model.FileBatch
	if err := bootstrap.DB.First(&batch, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse(model.CodeNotFound, "batch not found"))
		return
	}

	rawBody, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, "failed to read body"))
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	var input UpdateBatchRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, err.Error()))
		return
	}

	var rawMap map[string]interface{}
	json.Unmarshal(rawBody, &rawMap)

	updates := make(map[string]interface{})
	if input.Remark != nil {
		updates["remark"] = *input.Remark
	}
	if input.ExpireType != nil {
		newType := *input.ExpireType
		updates["expire_type"] = newType

		// 如果类型发生变化，根据新类型清除不相关的配置
		if newType != batch.ExpireType {
			if newType == "download" {
				updates["expire_at"] = nil
			} else if newType == "time" {
				updates["max_downloads"] = 0
			} else if newType == "permanent" {
				updates["expire_at"] = nil
				updates["max_downloads"] = 0
			}
		}
	}

	// 显式提供的值具有最高优先级，但仅在逻辑允许的情况下
	// 例如：如果切换到了 download 类型，用户可以同时提供一个新的 max_downloads
	if _, ok := rawMap["expire_at"]; ok {
		updates["expire_at"] = input.ExpireAt
	}
	if input.MaxDownloads != nil {
		updates["max_downloads"] = *input.MaxDownloads
	}

	// 强制校验：如果最终结果是 permanent，确保限制被清空
	// 这样即使用户在请求中显式传了非零值，也会被修正
	finalType := batch.ExpireType
	if t, ok := updates["expire_type"].(string); ok {
		finalType = t
	}

	if finalType == "permanent" {
		updates["expire_at"] = nil
		updates["max_downloads"] = 0
	} else if finalType == "time" {
		// 如果是时间过期，max_downloads 应该始终为 0
		updates["max_downloads"] = 0
	} else if finalType == "download" {
		// 如果是下载次数过期，expire_at 应该始终为 null
		updates["expire_at"] = nil
	}
	if input.DownloadCount != nil {
		updates["download_count"] = *input.DownloadCount
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}

	if len(updates) > 0 {
		if err := bootstrap.DB.Model(&batch).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, err.Error()))
			return
		}
		// 重新从数据库读取，确保返回的是完整且最新的数据
		bootstrap.DB.First(&batch, "id = ?", id)
	}

	c.JSON(http.StatusOK, model.SuccessResponse(batch))
}

// CleanBatches 手动触发清理过期或已删除的批次
// @Summary      手动触发清理
// @Description  手动扫描并物理删除所有已过期或标记为删除的文件批次及其关联文件
// @Tags         Admin
// @Security     AdminAuth
// @Produce      json
// @Success      200  {object}  model.Response
// @Failure      500  {object}  model.Response
// @Router       /api/admin/batches/clean [post]
func (h *BatchHandler) CleanBatches(c *gin.Context) {
	slog.Info("Admin triggered manual cleanup")
	if err := h.batchService.Cleanup(c.Request.Context()); err != nil {
		slog.Error("Manual cleanup failed", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "cleanup failed: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// DeleteBatch 删除批次
// @Summary      删除批次
// @Description  标记批次为已删除，并物理删除关联的存储文件
// @Tags         Admin
// @Security     AdminAuth
// @Param        batch_id  path  string  true  "批次 ID (UUID)"
// @Produce      json
// @Success      200  {object}  model.Response
// @Failure      500  {object}  model.Response
// @Router       /api/admin/batches/{batch_id} [delete]
func (h *BatchHandler) DeleteBatch(c *gin.Context) {
	id := c.Param("batch_id")

	if err := h.batchService.DeleteBatch(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(map[string]interface{}{}))
}
