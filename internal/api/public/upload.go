package public

import (
	"FileRelay/internal/config"
	"FileRelay/internal/model"
	"FileRelay/internal/service"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService *service.UploadService
}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{
		uploadService: service.NewUploadService(),
	}
}

type UploadResponse struct {
	PickupCode string     `json:"pickup_code"`
	ExpireAt   *time.Time `json:"expire_at"`
	BatchID    string     `json:"batch_id"`
}

// Upload 上传文件并生成取件码
// @Summary      上传文件
// @Description  上传一个或多个文件并创建一个提取批次。如果配置了 require_token，则必须提供带 upload scope 的 API Token。
// @Tags         Public
// @Accept       multipart/form-data
// @Produce      json
// @Security     APITokenAuth
// @Param        files       formData  file    true   "文件列表"
// @Param        remark      formData  string  false  "备注"
// @Param        expire_type formData  string  false  "过期类型 (time/download/permanent)"
// @Param        expire_days formData  int     false  "过期天数 (针对 time 类型)"
// @Param        max_downloads formData int    false  "最大下载次数 (针对 download 类型)"
// @Success      200  {object}  model.Response{data=UploadResponse}
// @Failure      400  {object}  model.Response
// @Failure      500  {object}  model.Response
// @Router       /api/batches [post]
func (h *UploadHandler) Upload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, "invalid form"))
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, "no files uploaded"))
		return
	}

	if len(files) > config.GlobalConfig.Upload.MaxBatchFiles {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, "too many files"))
		return
	}

	// 校验单个文件大小
	maxSize := config.GlobalConfig.Upload.MaxFileSizeMB * 1024 * 1024
	for _, file := range files {
		if file.Size > maxSize {
			c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, fmt.Sprintf("文件 %s 超过最大限制 (%dMB)", file.Filename, config.GlobalConfig.Upload.MaxFileSizeMB)))
			return
		}
	}

	remark := c.PostForm("remark")
	expireType := c.PostForm("expire_type") // time / download / permanent
	if expireType == "" {
		expireType = "time"
	}

	var expireValue interface{}
	switch expireType {
	case "time":
		days, _ := strconv.Atoi(c.PostForm("expire_days"))
		if days <= 0 {
			days = config.GlobalConfig.Upload.MaxRetentionDays
		}
		expireValue = days
	case "download":
		max, _ := strconv.Atoi(c.PostForm("max_downloads"))
		if max <= 0 {
			max = 1
		}
		expireValue = max
	}

	batch, err := h.uploadService.CreateBatch(c.Request.Context(), files, remark, expireType, expireValue)
	if err != nil {
		slog.Error("Upload handler failed to create batch", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(UploadResponse{
		PickupCode: batch.PickupCode,
		ExpireAt:   batch.ExpireAt,
		BatchID:    batch.ID,
	}))
}

type UploadTextRequest struct {
	Content      string `json:"content" binding:"required" example:"这是一段长文本内容..."`
	Remark       string `json:"remark" example:"文本备注"`
	ExpireType   string `json:"expire_type" example:"time"`
	ExpireDays   int    `json:"expire_days" example:"7"`
	MaxDownloads int    `json:"max_downloads" example:"5"`
}

// UploadText 发送长文本并生成取件码
// @Summary      发送长文本
// @Description  中转一段长文本内容并创建一个提取批次。如果配置了 require_token，则必须提供带 upload scope 的 API Token。
// @Tags         Public
// @Accept       json
// @Produce      json
// @Security     APITokenAuth
// @Param        request body UploadTextRequest true "文本内容及配置"
// @Success      200  {object}  model.Response{data=UploadResponse}
// @Failure      400  {object}  model.Response
// @Failure      500  {object}  model.Response
// @Router       /api/batches/text [post]
func (h *UploadHandler) UploadText(c *gin.Context) {
	var req UploadTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, err.Error()))
		return
	}

	// 校验文本长度
	maxSize := config.GlobalConfig.Upload.MaxFileSizeMB * 1024 * 1024
	if int64(len(req.Content)) > maxSize {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, fmt.Sprintf("文本内容超过最大限制 (%dMB)", config.GlobalConfig.Upload.MaxFileSizeMB)))
		return
	}

	if req.ExpireType == "" {
		req.ExpireType = "time"
	}

	var expireValue interface{}
	switch req.ExpireType {
	case "time":
		if req.ExpireDays <= 0 {
			req.ExpireDays = config.GlobalConfig.Upload.MaxRetentionDays
		}
		expireValue = req.ExpireDays
	case "download":
		if req.MaxDownloads <= 0 {
			req.MaxDownloads = 1
		}
		expireValue = req.MaxDownloads
	}

	batch, err := h.uploadService.CreateTextBatch(c.Request.Context(), req.Content, req.Remark, req.ExpireType, expireValue)
	if err != nil {
		slog.Error("Upload handler failed to create text batch", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(UploadResponse{
		PickupCode: batch.PickupCode,
		ExpireAt:   batch.ExpireAt,
		BatchID:    batch.ID,
	}))
}
