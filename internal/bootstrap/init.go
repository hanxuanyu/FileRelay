package bootstrap

import (
	"FileRelay/internal/config"
	"FileRelay/internal/model"
	"FileRelay/internal/storage"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"path/filepath"
	"runtime"

	"log/slog"
	"os"
	"strings"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitLog() {
	level := slog.LevelInfo
	switch strings.ToLower(config.GlobalConfig.Log.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	var handlers []slog.Handler

	// 1. 控制台处理器 (简化格式)
	handlers = append(handlers, &ConsoleHandler{
		out:   os.Stdout,
		level: level,
	})

	// 2. 文件处理器 (结构化 Text 格式)
	if config.GlobalConfig.Log.FilePath != "" {
		logDir := filepath.Dir(config.GlobalConfig.Log.FilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Printf("Warning: Failed to create log directory: %v\n", err)
		} else {
			file, err := os.OpenFile(config.GlobalConfig.Log.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("Warning: Failed to open log file: %v\n", err)
			} else {
				fileHandler := slog.NewTextHandler(file, &slog.HandlerOptions{
					Level:     level,
					AddSource: true,
					ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
						if a.Key == slog.TimeKey {
							return slog.String(a.Key, a.Value.Time().Format("2006-01-02 15:04:05.000"))
						}
						if a.Key == slog.SourceKey {
							source := a.Value.Any().(*slog.Source)
							return slog.String(a.Key, fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line))
						}
						return a
					},
				})
				handlers = append(handlers, fileHandler)
			}
		}
	}

	var finalHandler slog.Handler
	if len(handlers) == 1 {
		finalHandler = handlers[0]
	} else {
		finalHandler = &MultiHandler{handlers: handlers}
	}

	logger := slog.New(finalHandler).With("service", "filerelay")
	slog.SetDefault(logger)
}

func InitDB() {
	var err error
	cfg := config.GlobalConfig.Database

	DB, err = ConnectDB(cfg)
	if err != nil {
		slog.Error("Failed to initialize database", "type", cfg.Type, "error", err)
		os.Exit(1)
	}

	slog.Info("Database initialized and migrated", "type", cfg.Type)

	// 初始化存储
	if err := ReloadStorage(); err != nil {
		slog.Error("Failed to initialize storage", "error", err)
		os.Exit(1)
	}

	// 初始化管理员 (如果不存在)
	initAdmin()
}

func ConnectDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch strings.ToLower(cfg.Type) {
	case "mysql":
		params := cfg.Config
		if params == "" {
			params = "parseTime=True&loc=Local&charset=utf8mb4"
		} else {
			if !strings.Contains(strings.ToLower(params), "parsetime") {
				params += "&parseTime=True"
			}
			if !strings.Contains(strings.ToLower(params), "loc=") {
				params += "&loc=Local"
			}
			if !strings.Contains(strings.ToLower(params), "charset") {
				params += "&charset=utf8mb4"
			}
		}
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, params)
		dialector = mysql.Open(dsn)
	case "postgres", "postgresql":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d %s",
			cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.Config)
		dialector = postgres.Open(dsn)
	case "sqlite", "sqlite3":
		fallthrough
	default:
		dbPath := cfg.Path
		if dbPath == "" {
			dbPath = "data/file_relay.db"
		}
		dialector = sqlite.Open(dbPath)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移
	err = db.AutoMigrate(
		&model.FileBatch{},
		&model.FileItem{},
		&model.APIToken{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ReloadDB(newCfg config.DatabaseConfig) error {
	newDB, err := ConnectDB(newCfg)
	if err != nil {
		return err
	}

	if DB != nil {
		// 检查是否真的是不同的数据库，避免自迁移导致冲突
		// 这里简单判断类型或连接串是否变化，或者直接让用户决定
		// 为了安全，我们只在连接参数确实变化时才尝试迁移
		slog.Info("Starting data migration to new database...")
		if err := MigrateData(DB, newDB); err != nil {
			slog.Error("Data migration failed", "error", err)
			// 迁移失败不一定需要中断，但需要记录
		} else {
			slog.Info("Data migration completed successfully")
		}
	}

	DB = newDB
	return nil
}

func MigrateData(sourceDB, targetDB *gorm.DB) error {
	// 迁移 APIToken
	var tokens []model.APIToken
	if err := sourceDB.Find(&tokens).Error; err == nil && len(tokens) > 0 {
		if err := targetDB.Save(&tokens).Error; err != nil {
			slog.Warn("Failed to migrate APITokens", "error", err)
		}
	}

	// 迁移 FileBatch (分批处理以节省内存)
	var batches []model.FileBatch
	err := sourceDB.Model(&model.FileBatch{}).FindInBatches(&batches, 100, func(tx *gorm.DB, batch int) error {
		if err := targetDB.Save(&batches).Error; err != nil {
			return err
		}
		return nil
	}).Error
	if err != nil {
		slog.Warn("Failed to migrate FileBatches", "error", err)
	}

	// 迁移 FileItem (分批处理以节省内存)
	var items []model.FileItem
	err = sourceDB.Model(&model.FileItem{}).FindInBatches(&items, 100, func(tx *gorm.DB, batch int) error {
		if err := targetDB.Save(&items).Error; err != nil {
			return err
		}
		return nil
	}).Error
	if err != nil {
		slog.Warn("Failed to migrate FileItems", "error", err)
	}

	return nil
}

func ReloadStorage() error {
	storageType := config.GlobalConfig.Storage.Type
	switch storageType {
	case "local":
		storage.GlobalStorage = storage.NewLocalStorage(config.GlobalConfig.Storage.Local.Path)
	case "webdav":
		cfg := config.GlobalConfig.Storage.WebDAV
		storage.GlobalStorage = storage.NewWebDAVStorage(cfg.URL, cfg.Username, cfg.Password, cfg.Root)
	case "s3":
		cfg := config.GlobalConfig.Storage.S3
		s3Storage, err := storage.NewS3Storage(context.Background(), cfg.Endpoint, cfg.Region, cfg.AccessKey, cfg.SecretKey, cfg.Bucket, cfg.UseSSL)
		if err != nil {
			return err
		}
		storage.GlobalStorage = s3Storage
	default:
		return fmt.Errorf("unsupported storage type: %s", storageType)
	}
	slog.Info("Storage initialized", "type", storageType)
	return nil
}

func initAdmin() {
	passwordHash := config.GlobalConfig.Security.AdminPasswordHash
	if passwordHash == "" {
		// 生成随机密码
		password := generateRandomPassword(12)
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("Failed to generate password hash", "error", err)
			os.Exit(1)
		}
		passwordHash = string(hash)
		fmt.Printf("**************************************************\n")
		fmt.Printf("NO ADMIN PASSWORD CONFIGURED. GENERATED RANDOM PASSWORD:\n")
		fmt.Printf("Password: %s\n", password)
		fmt.Printf("Please save this password or configure admin_password_hash in config.yaml\n")
		fmt.Printf("**************************************************\n")

		// 将生成的哈希保存回配置文件
		config.GlobalConfig.Security.AdminPasswordHash = passwordHash
		if err := config.SaveConfig(); err != nil {
			slog.Warn("Failed to save generated password hash to config", "error", err)
		}
	}
	slog.Info("Admin authentication initialized")
}

func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "admin123" // 退路
		}
		b[i] = charset[num.Int64()]
	}
	return string(b)
}

// ConsoleHandler 实现 simplified 控制台日志
type ConsoleHandler struct {
	out   io.Writer
	level slog.Leveler
	attrs []slog.Attr
}

func (h *ConsoleHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.level.Level()
}

func (h *ConsoleHandler) Handle(_ context.Context, r slog.Record) error {
	var b strings.Builder

	// 时间: 15:04:05.000
	b.WriteString(r.Time.Format("15:04:05.000"))
	b.WriteString(" ")

	// 级别: INFO
	level := r.Level.String()
	b.WriteString(fmt.Sprintf("%-5s", level))
	b.WriteString(" ")

	// 源码: [main.go:123]
	pc := r.PC
	if pc == 0 {
		// 如果 Logger 没采集 PC (自定义 Handler 默认情况)，我们手动采集
		// 跳过: 1:runtime.Callers, 2:Handle, 3:slog.(*Logger).log, 4:slog.(*Logger).Info
		var pcs [1]uintptr
		runtime.Callers(4, pcs[:])
		pc = pcs[0]
	}

	if pc != 0 {
		fs := runtime.CallersFrames([]uintptr{pc})
		f, _ := fs.Next()
		b.WriteString(fmt.Sprintf("[%s:%d] ", filepath.Base(f.File), f.Line))
	}

	// 消息
	b.WriteString(r.Message)

	// 属性: 仅输出值
	writeAttr := func(a slog.Attr) {
		if a.Key == "service" {
			return
		}
		b.WriteString(" ")
		val := a.Value.Resolve().String()
		if strings.Contains(val, " ") {
			b.WriteString(fmt.Sprintf("%q", val))
		} else {
			b.WriteString(val)
		}
	}

	for _, a := range h.attrs {
		writeAttr(a)
	}
	r.Attrs(func(a slog.Attr) bool {
		writeAttr(a)
		return true
	})

	b.WriteString("\n")
	_, err := h.out.Write([]byte(b.String()))
	return err
}

func (h *ConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)
	return &ConsoleHandler{
		out:   h.out,
		level: h.level,
		attrs: newAttrs,
	}
}

func (h *ConsoleHandler) WithGroup(name string) slog.Handler {
	return h // 简化版暂不支持分组
}

// MultiHandler 实现多路分发
type MultiHandler struct {
	handlers []slog.Handler
}

func (h *MultiHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for _, hh := range h.handlers {
		if hh.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

func (h *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, hh := range h.handlers {
		if hh.Enabled(ctx, r.Level) {
			_ = hh.Handle(ctx, r)
		}
	}
	return nil
}

func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, hh := range h.handlers {
		newHandlers[i] = hh.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: newHandlers}
}

func (h *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, hh := range h.handlers {
		newHandlers[i] = hh.WithGroup(name)
	}
	return &MultiHandler{handlers: newHandlers}
}
