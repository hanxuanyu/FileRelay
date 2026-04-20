package public

import (
	"FileRelay/internal/bootstrap"
	"FileRelay/internal/model"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

type HealthResponse struct {
	Status   string `json:"status" example:"ok"`
	Uptime   string `json:"uptime" example:"2h30m15s"`
	Database string `json:"database" example:"ok"`
	GoVer    string `json:"go_version" example:"go1.21.0"`
}

// Health 健康检查
// @Summary      健康检查
// @Description  检查服务运行状态，需要有效的 API Token 或管理员 JWT 进行鉴权。
// @Tags         Public
// @Produce      json
// @Security     APITokenAuth
// @Success      200  {object}  model.Response{data=HealthResponse}
// @Failure      401  {object}  model.Response
// @Router       /api/health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	dbStatus := "ok"
	sqlDB, err := bootstrap.DB.DB()
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = "unavailable"
	}

	status := "ok"
	if dbStatus != "ok" {
		status = "degraded"
	}

	c.JSON(http.StatusOK, model.SuccessResponse(HealthResponse{
		Status:   status,
		Uptime:   time.Since(startTime).Round(time.Second).String(),
		Database: dbStatus,
		GoVer:    runtime.Version(),
	}))
}
