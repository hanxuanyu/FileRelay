package public

import (
	"FileRelay/internal/api/middleware"
	"FileRelay/internal/bootstrap"
	"FileRelay/internal/config"
	"FileRelay/internal/model"
	"FileRelay/internal/service"
	"FileRelay/internal/storage"
	"archive/zip"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PickupResponse struct {
	Remark        string           `json:"remark"`
	ExpireAt      *time.Time       `json:"expire_at"`
	ExpireType    string           `json:"expire_type"`
	DownloadCount int              `json:"download_count"`
	MaxDownloads  int              `json:"max_downloads"`
	Type          string           `json:"type"`
	Content       string           `json:"content,omitempty"`
	Files         []model.FileItem `json:"files,omitempty"`
}

type DownloadCountResponse struct {
	DownloadCount int `json:"download_count"`
	MaxDownloads  int `json:"max_downloads"`
}

// DownloadBatch 批量下载文件 (ZIP)
// @Summary      批量下载文件
// @Description  根据取件码将批次内的所有文件打包为 ZIP 格式一次性下载。可选提供带 pickup scope 的 API Token。
// @Tags         Public
// @Security     APITokenAuth
// @Param        pickup_code  path  string  true  "取件码"
// @Produce      application/zip
// @Success      200  {file}    file
// @Failure      404  {object}  model.Response
// @Router       /api/batches/{pickup_code}/download [get]
func (h *PickupHandler) DownloadBatch(c *gin.Context) {
	code := c.Param("pickup_code")
	batch, err := h.batchService.GetBatchByPickupCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse(model.CodeNotFound, "batch not found or expired"))
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"batch_%s.zip\"", code))
	c.Header("Content-Type", "application/zip")

	zw := zip.NewWriter(c.Writer)
	defer zw.Close()

	for _, item := range batch.FileItems {
		reader, err := storage.GlobalStorage.Open(c.Request.Context(), item.StoragePath)
		if err != nil {
			continue // Skip failed files
		}

		f, err := zw.Create(item.OriginalName)
		if err != nil {
			reader.Close()
			continue
		}

		_, _ = io.Copy(f, reader)
		reader.Close()
	}

	// 增加下载次数
	if err := h.batchService.IncrementDownloadCount(batch.ID); err != nil {
		slog.Error("Failed to increment download count", "batch_id", batch.ID, "error", err)
	}
}

type PickupHandler struct {
	batchService *service.BatchService
}

func NewPickupHandler() *PickupHandler {
	return &PickupHandler{
		batchService: service.NewBatchService(),
	}
}

// Pickup 获取批次信息
// @Summary      获取批次信息
// @Description  根据取件码获取文件批次详细信息和文件列表。可选提供带 pickup scope 的 API Token。
// @Tags         Public
// @Security     APITokenAuth
// @Produce      json
// @Param        pickup_code  path  string  true  "取件码"
// @Success      200  {object}  model.Response{data=PickupResponse}
// @Failure      404  {object}  model.Response
// @Router       /api/batches/{pickup_code} [get]
func (h *PickupHandler) Pickup(c *gin.Context) {
	code := c.Param("pickup_code")
	if code == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, "pickup code required"))
		return
	}

	batch, err := h.batchService.GetBatchByPickupCode(code)
	if err != nil {
		middleware.RecordPickupFailure(c.ClientIP())
		c.JSON(http.StatusNotFound, model.ErrorResponse(model.CodeNotFound, "batch not found or expired"))
		return
	}

	if batch.Type == "text" {
		if err := h.batchService.IncrementDownloadCount(batch.ID); err != nil {
			slog.Error("Failed to increment download count for batch", "batch_id", batch.ID, "error", err)
		} else {
			batch.DownloadCount++
		}
	}

	baseURL := getBaseURL(c)

	for i := range batch.FileItems {
		batch.FileItems[i].DownloadURL = fmt.Sprintf("%s/api/files/%s/%s", baseURL, batch.FileItems[i].ID, batch.FileItems[i].OriginalName)
	}

	c.JSON(http.StatusOK, model.SuccessResponse(PickupResponse{
		Remark:        batch.Remark,
		ExpireAt:      batch.ExpireAt,
		ExpireType:    batch.ExpireType,
		DownloadCount: batch.DownloadCount,
		MaxDownloads:  batch.MaxDownloads,
		Type:          batch.Type,
		Content:       batch.Content,
		Files:         batch.FileItems,
	}))
}

// GetDownloadCount 查询下载次数
// @Summary      查询下载次数
// @Description  根据取件码查询当前下载次数和最大允许下载次数。支持已过期的批次。
// @Tags         Public
// @Produce      json
// @Param        pickup_code  path  string  true  "取件码"
// @Success      200  {object}  model.Response{data=DownloadCountResponse}
// @Failure      400  {object}  model.Response
// @Failure      404  {object}  model.Response
// @Router       /api/batches/{pickup_code}/count [get]
func (h *PickupHandler) GetDownloadCount(c *gin.Context) {
	code := c.Param("pickup_code")
	if code == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, "pickup code required"))
		return
	}

	count, max, err := h.batchService.GetDownloadCountByPickupCode(code)
	if err != nil {
		middleware.RecordPickupFailure(c.ClientIP())
		c.JSON(http.StatusNotFound, model.ErrorResponse(model.CodeNotFound, "batch not found"))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(DownloadCountResponse{
		DownloadCount: count,
		MaxDownloads:  max,
	}))
}

func getBaseURL(c *gin.Context) string {
	// 优先使用配置中的 BaseURL
	if config.GlobalConfig.Site.BaseURL != "" {
		return strings.TrimSuffix(config.GlobalConfig.Site.BaseURL, "/")
	}

	// 自动检测逻辑
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	} else {
		// 检查常用的代理协议头 (优先)
		// 增加对用户提供的 :scheme (可能被某些代理转为普通 header) 的支持
		// 增加对 X-Forwarded-Proto 可能存在的逗号分隔列表的处理
		checkHeaders := []struct {
			name   string
			values []string
		}{
			{"X-Forwarded-Proto", []string{"https"}},
			{"X-Forwarded-Protocol", []string{"https"}},
			{"X-Url-Scheme", []string{"https"}},
			{"Front-End-Https", []string{"on", "https"}},
			{"X-Forwarded-Ssl", []string{"on", "https"}},
			{":scheme", []string{"https"}},
			{"X-Scheme", []string{"https"}},
		}

		found := false
		for _, h := range checkHeaders {
			val := c.GetHeader(h.name)
			if val == "" {
				continue
			}
			// 处理可能的逗号分隔列表 (如 X-Forwarded-Proto: https, http)
			firstVal := strings.TrimSpace(strings.ToLower(strings.Split(val, ",")[0]))
			for _, target := range h.values {
				if firstVal == target {
					scheme = "https"
					found = true
					break
				}
			}
			if found {
				break
			}
		}

		// 检查 Forwarded 头部 (RFC 7239)
		if !found {
			if forwarded := c.GetHeader("Forwarded"); forwarded != "" {
				if strings.Contains(strings.ToLower(forwarded), "proto=https") {
					scheme = "https"
					found = true
				}
			}
		}

		// 启发式判断：如果上述头部都没有，但 Referer 是 https，则认为也是 https
		// 这在同域 API 请求时非常可靠
		if !found {
			if referer := c.GetHeader("Referer"); strings.HasPrefix(strings.ToLower(referer), "https://") {
				scheme = "https"
			}
		}
	}

	host := c.Request.Host
	if forwardedHost := c.GetHeader("X-Forwarded-Host"); forwardedHost != "" {
		// 处理可能的逗号分隔列表
		host = strings.TrimSpace(strings.Split(forwardedHost, ",")[0])
	}

	return fmt.Sprintf("%s://%s", scheme, host)
}

// DownloadFile 下载单个文件
// @Summary      下载单个文件
// @Description  根据文件 ID 下载单个文件。支持直观的文件名结尾以方便下载工具识别。可选提供带 pickup scope 的 API Token。
// @Tags         Public
// @Security     APITokenAuth
// @Param        file_id   path  string  true  "文件 ID (UUID)"
// @Param        filename  path  string  false "文件名"
// @Produce      application/octet-stream
// @Success      200  {file}    file
// @Failure      404  {object}  model.Response
// @Failure      410  {object}  model.Response
// @Router       /api/files/{file_id}/{filename} [get]
// @Router       /api/files/{file_id}/download [get]
func (h *PickupHandler) DownloadFile(c *gin.Context) {
	fileID := c.Param("file_id")

	var item model.FileItem
	if err := bootstrap.DB.First(&item, "id = ?", fileID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse(model.CodeNotFound, "file not found"))
		return
	}

	var batch model.FileBatch
	if err := bootstrap.DB.First(&batch, "id = ?", item.BatchID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse(model.CodeNotFound, "batch not found"))
		return
	}

	if h.batchService.IsExpired(&batch) {
		h.batchService.MarkAsExpired(&batch)
		// 按照需求，如果不存在（已在上面处理）或达到上限，返回 404
		if batch.ExpireType == "download" && batch.DownloadCount >= batch.MaxDownloads {
			c.JSON(http.StatusNotFound, model.ErrorResponse(model.CodeNotFound, "file not found or download limit reached"))
		} else {
			c.JSON(http.StatusGone, model.ErrorResponse(model.CodeGone, "batch expired"))
		}
		return
	}

	// 打开文件
	reader, err := storage.GlobalStorage.Open(c.Request.Context(), item.StoragePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, "failed to open file"))
		return
	}
	defer reader.Close()

	// 增加下载次数
	if err := h.batchService.IncrementDownloadCount(batch.ID); err != nil {
		// 记录错误但不中断下载过程
		slog.Error("Failed to increment download count for batch", "batch_id", batch.ID, "error", err)
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", item.OriginalName))
	c.Header("Content-Type", item.MimeType)
	c.Header("Content-Length", strconv.FormatInt(item.Size, 10))

	// 如果是 HEAD 请求，只返回 Header
	if c.Request.Method == http.MethodHead {
		return
	}

	if _, err := io.Copy(c.Writer, reader); err != nil {
		slog.Error("Error during file download", "file_id", item.ID, "error", err)
	}
}
