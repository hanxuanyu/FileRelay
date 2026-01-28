package middleware

import (
	"FileRelay/internal/config"
	"FileRelay/internal/model"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	pickupFailures = make(map[string]int)
	failureMutex   sync.Mutex
)

func PickupRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		failureMutex.Lock()
		count, exists := pickupFailures[key]
		failureMutex.Unlock()

		if exists && count >= config.GlobalConfig.Security.PickupFailLimit {
			slog.Warn("Pickup rate limit exceeded", "ip", key, "count", count)
			c.JSON(http.StatusTooManyRequests, model.ErrorResponse(model.CodeTooManyRequests, "Too many failed attempts. Please try again later."))
			c.Abort()
			return
		}

		c.Next()
	}
}

func RecordPickupFailure(ip string) {
	key := ip
	failureMutex.Lock()
	pickupFailures[key]++

	// 仅在第一次失败时启动清除记录的计时器
	if pickupFailures[key] == 1 {
		go func() {
			// 设置 1 分钟后清除记录 (简单实现)
			time.Sleep(1 * time.Hour)
			failureMutex.Lock()
			delete(pickupFailures, key)
			slog.Info("Pickup failure record cleared", "ip", key)
			failureMutex.Unlock()
		}()
	}

	failureMutex.Unlock()
}
