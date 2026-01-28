package main

import (
	_ "FileRelay/docs"
	"FileRelay/internal/api/admin"
	"FileRelay/internal/api/middleware"
	"FileRelay/internal/api/public"
	"FileRelay/internal/bootstrap"
	"FileRelay/internal/config"
	"FileRelay/internal/model"
	"FileRelay/internal/task"
	"context"
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed all:web
var webFS embed.FS

// @title           文件暂存柜 API
// @version         1.0
// @description     自托管的文件暂存柜后端系统 API 文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /

// @securityDefinitions.apikey AdminAuth
// @in header
// @name Authorization
// @description Type "Bearer <JWT-Token>" or "Bearer <API-Token>" to authenticate. API Token must have 'admin' scope.

// @securityDefinitions.apikey APITokenAuth
// @in header
// @name Authorization
// @description Type "Bearer <API-Token>" to authenticate. Required scope depends on the endpoint.

func main() {
	// 注册常用 MIME 类型
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".woff", "font/woff")
	mime.AddExtensionType(".woff2", "font/woff2")
	mime.AddExtensionType(".svg", "image/svg+xml")

	// 解析命令行参数
	configPath := flag.String("config", "config/config.yaml", "path to config file")
	flag.Parse()

	// 1. 加载配置
	if err := config.LoadConfig(*configPath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 初始化日志
	bootstrap.InitLog()

	port := config.GlobalConfig.Site.Port
	if port <= 0 {
		port = 8080
	}

	printBanner(*configPath, port)

	// 2. 初始化
	bootstrap.InitDB()

	// 3. 启动清理任务
	cleaner := task.NewCleaner()
	go cleaner.Start(context.Background())

	// 4. 设置路由
	r := gin.New()
	r.Use(gin.Recovery())

	// 自定义 Gin 日志中间件，使用 slog
	r.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				slog.Error("Request error", "error", e, "path", path)
			}
		} else {
			attributes := []any{
				"status", c.Writer.Status(),
				"method", c.Request.Method,
				"path", path,
				"ip", c.ClientIP(),
				"latency", latency,
			}
			if query != "" {
				attributes = append(attributes, "query", query)
			}
			slog.Info("Request", attributes...)
		}
	})

	// 配置更完善的 CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization", "Accept", "X-Requested-With")
	corsConfig.AllowMethods = append(corsConfig.AllowMethods, "OPTIONS")
	r.Use(cors.New(corsConfig))

	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 公共接口
	uploadHandler := public.NewUploadHandler()
	pickupHandler := public.NewPickupHandler()
	publicConfigHandler := public.NewConfigHandler()

	api := r.Group("/api")
	{
		api.GET("/config", publicConfigHandler.GetPublicConfig)
		// 统一使用 /batches 作为资源路径
		api.POST("/batches", middleware.UploadAuth(), uploadHandler.Upload)
		api.POST("/batches/text", middleware.UploadAuth(), uploadHandler.UploadText)
		api.GET("/batches/:pickup_code", middleware.PickupRateLimit(), middleware.APITokenAuth(model.ScopePickup, true), pickupHandler.Pickup)
		api.GET("/batches/:pickup_code/count", middleware.PickupRateLimit(), pickupHandler.GetDownloadCount)
		api.GET("/batches/:pickup_code/download", middleware.APITokenAuth(model.ScopePickup, true), pickupHandler.DownloadBatch)
		// 文件下载路由，支持直观的文件名结尾
		api.GET("/files/:file_id/:filename", middleware.APITokenAuth(model.ScopePickup, true), pickupHandler.DownloadFile)
		// 保持旧路由兼容性
		api.GET("/files/:file_id/download", middleware.APITokenAuth(model.ScopePickup, true), pickupHandler.DownloadFile)

		// 保持旧路由兼容性 (可选，但为了平滑过渡通常建议保留一段时间或直接更新)
		// 这里根据需求“调整不符合规范的”，我将直接采用新路由
	}

	// 管理员接口
	authHandler := admin.NewAuthHandler()
	batchHandler := admin.NewBatchHandler()
	tokenHandler := admin.NewTokenHandler()
	configHandler := admin.NewConfigHandler()

	api.POST("/admin/login", authHandler.Login)

	adm := api.Group("/admin")
	adm.Use(middleware.AdminAuth())
	{
		adm.GET("/config", configHandler.GetConfig)
		adm.PUT("/config", configHandler.UpdateConfig)

		adm.GET("/batches", batchHandler.ListBatches)
		adm.POST("/batches/clean", batchHandler.CleanBatches)
		adm.GET("/batches/:batch_id", batchHandler.GetBatch)
		adm.PUT("/batches/:batch_id", batchHandler.UpdateBatch)
		adm.DELETE("/batches/:batch_id", batchHandler.DeleteBatch)

		adm.GET("/api-tokens", tokenHandler.ListTokens)
		adm.POST("/api-tokens", tokenHandler.CreateToken)
		adm.DELETE("/api-tokens/:id", tokenHandler.DeleteToken)
		adm.POST("/api-tokens/:id/revoke", tokenHandler.RevokeToken)
		adm.POST("/api-tokens/:id/recover", tokenHandler.RecoverToken)
	}

	// 静态资源服务 (放在最后，确保 API 路由优先)
	webSub, _ := fs.Sub(webFS, "web")

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 如果请求的是 API 或 Swagger，则不处理静态资源 (让其返回 404)
		// 注意：此处不排除 /admin，因为 /admin 通常是前端 SPA 的路由地址
		if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/swagger") {
			return
		}

		// 辅助函数：尝试从外部或嵌入服务文件
		serveFile := func(relPath string, allowExternal bool) bool {
			// 1. 优先尝试外部路径
			if allowExternal && config.GlobalConfig.Web.Path != "" {
				fullPath := filepath.Join(config.GlobalConfig.Web.Path, relPath)
				if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
					c.File(fullPath)
					return true
				}
			}

			// 2. 尝试嵌入式文件
			f, err := webSub.Open(relPath)
			if err == nil {
				defer f.Close()
				stat, err := f.Stat()
				if err == nil && !stat.IsDir() {
					// 使用 http.ServeContent 避免 c.FileFromFS 重定向问题
					if rs, ok := f.(io.ReadSeeker); ok {
						// 显式设置 Content-Type，防止某些环境下识别失败
						ext := filepath.Ext(relPath)
						ctype := mime.TypeByExtension(ext)
						if ctype != "" {
							c.Header("Content-Type", ctype)
						}
						http.ServeContent(c.Writer, c.Request, stat.Name(), stat.ModTime(), rs)
						return true
					}
				}
			}
			return false
		}

		// 1. 尝试直接请求的文件 (如果是 / 则尝试 index.html)
		requestedPath := strings.TrimPrefix(path, "/")
		if requestedPath == "" {
			requestedPath = "index.html"
		}

		if serveFile(requestedPath, true) {
			return
		}

		// 2. SPA 支持：对于非文件请求（没有后缀或不包含点），尝试返回 index.html
		// 如果请求的是 assets 目录下的文件（包含点）却没找到，不应该回退到 index.html
		isAsset := strings.Contains(requestedPath, ".")
		if !isAsset || requestedPath == "index.html" {
			if serveFile("index.html", true) {
				return
			}
		}

		// 最终找不到则返回 404
		c.Status(http.StatusNotFound)
	})

	// 5. 运行
	localIP := getLocalIP()
	fmt.Println("FileRelay service is ready:")
	fmt.Printf("  - Local:   http://localhost:%d\n", port)
	fmt.Printf("  - LAN:     http://%s:%d\n", localIP, port)
	fmt.Printf("  - API Doc: http://localhost:%d/swagger/index.html\n", port)
	fmt.Println()

	slog.Info("FileRelay service is ready",
		"local_url", fmt.Sprintf("http://localhost:%d", port),
		"lan_url", fmt.Sprintf("http://%s:%d", localIP, port),
		"docs_url", fmt.Sprintf("http://localhost:%d/swagger/index.html", port),
	)

	r.Run(fmt.Sprintf(":%d", port))
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func printBanner(configPath string, port int) {
	fmt.Println(`
   _____ _ _      _____      _             
  |  ___(_) | ___|  _  \ ___| | __ _ _   _ 
  | |_  | | |/ _ \ |_) / _ \ |/ _' | | | |
  |  _| | | |  __/  _ <  __/ | (_| | |_| |
  |_|   |_|_|\___|_| \_\___|_|\__,_|\__, |
                                    |___/ 
	`)
	fmt.Printf("  Config:       %s\n", configPath)
	fmt.Printf("  Port:         %d\n", port)
	fmt.Printf("  Database:     %s\n", config.GlobalConfig.Database.Type)
	fmt.Printf("  Storage:      %s\n", config.GlobalConfig.Storage.Type)

	webPath := config.GlobalConfig.Web.Path
	if webPath != "" {
		if info, err := os.Stat(webPath); err == nil && info.IsDir() {
			fmt.Printf("  Web Assets:   External (%s)\n", webPath)
		} else {
			fmt.Printf("  Web Assets:   Embedded (External path %s not found)\n", webPath)
		}
	} else {
		fmt.Println("  Web Assets:   Embedded")
	}
	fmt.Println()

	slog.Info("Starting FileRelay service",
		"config", configPath,
		"port", port,
		"db_type", config.GlobalConfig.Database.Type,
		"storage_type", config.GlobalConfig.Storage.Type,
	)
}
