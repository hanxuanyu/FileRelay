package config

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Site     SiteConfig     `yaml:"site" json:"site"`           // 站点设置
	Security SecurityConfig `yaml:"security" json:"security"`   // 安全设置
	Upload   UploadConfig   `yaml:"upload" json:"upload"`       // 上传设置
	Storage  StorageConfig  `yaml:"storage" json:"storage"`     // 存储设置
	APIToken APITokenConfig `yaml:"api_token" json:"api_token"` // API Token 设置
	Database DatabaseConfig `yaml:"database" json:"database"`   // 数据库设置
	Web      WebConfig      `yaml:"web" json:"web"`             // Web 前端设置
	Log      LogConfig      `yaml:"log" json:"log"`             // 日志设置
}

type LogConfig struct {
	Level    string `yaml:"level" json:"level"`         // 日志级别: debug, info, warn, error
	FilePath string `yaml:"file_path" json:"file_path"` // 日志文件路径，为空则仅输出到控制台
}

type WebConfig struct {
	Path string `yaml:"path" json:"path"` // Web 前端资源路径
}

type SiteConfig struct {
	Name        string `yaml:"name" json:"name"`               // 站点名称
	Description string `yaml:"description" json:"description"` // 站点描述
	Logo        string `yaml:"logo" json:"logo"`               // 站点 Logo URL
	BaseURL     string `yaml:"base_url" json:"base_url"`       // 站点外部访问地址 (例如: https://file.example.com)
	Port        int    `yaml:"port" json:"port"`               // 监听端口
}

type SecurityConfig struct {
	AdminPasswordHash string `yaml:"admin_password_hash" json:"admin_password_hash"` // 管理员密码哈希 (bcrypt)
	AdminPassword     string `yaml:"-" json:"admin_password,omitempty"`              // 管理员密码明文 (仅用于更新请求，不保存到文件)
	PickupCodeLength  int    `yaml:"pickup_code_length" json:"pickup_code_length"`   // 取件码长度 (变更后将自动通过右侧补零或截取调整存量数据)
	PickupFailLimit   int    `yaml:"pickup_fail_limit" json:"pickup_fail_limit"`     // 取件失败尝试限制
	JWTSecret         string `yaml:"jwt_secret" json:"jwt_secret"`                   // JWT 签名密钥
}

type UploadConfig struct {
	MaxFileSizeMB    int64 `yaml:"max_file_size_mb" json:"max_file_size_mb"`     // 单个文件最大大小 (MB)
	MaxBatchFiles    int   `yaml:"max_batch_files" json:"max_batch_files"`       // 每个批次最大文件数
	MaxRetentionDays int   `yaml:"max_retention_days" json:"max_retention_days"` // 最大保留天数
	RequireToken     bool  `yaml:"require_token" json:"require_token"`           // 是否强制要求上传 Token
}

type StorageConfig struct {
	Type  string `yaml:"type" json:"type"` // 存储类型: local, webdav, s3
	Local struct {
		Path string `yaml:"path" json:"path"` // 本地存储路径
	} `yaml:"local" json:"local"`
	WebDAV struct {
		URL      string `yaml:"url" json:"url"`           // WebDAV 地址
		Username string `yaml:"username" json:"username"` // WebDAV 用户名
		Password string `yaml:"password" json:"password"` // WebDAV 密码
		Root     string `yaml:"root" json:"root"`         // WebDAV 根目录
	} `yaml:"webdav" json:"webdav"`
	S3 struct {
		Endpoint  string `yaml:"endpoint" json:"endpoint"`     // S3 端点
		Region    string `yaml:"region" json:"region"`         // S3 区域
		AccessKey string `yaml:"access_key" json:"access_key"` // S3 Access Key
		SecretKey string `yaml:"secret_key" json:"secret_key"` // S3 Secret Key
		Bucket    string `yaml:"bucket" json:"bucket"`         // S3 Bucket
		UseSSL    bool   `yaml:"use_ssl" json:"use_ssl"`       // 是否使用 SSL
	} `yaml:"s3" json:"s3"`
}

type APITokenConfig struct {
	Enabled       bool `yaml:"enabled" json:"enabled"`                 // 是否启用 API Token
	AllowAdminAPI bool `yaml:"allow_admin_api" json:"allow_admin_api"` // 是否允许 API Token 访问管理接口
	MaxTokens     int  `yaml:"max_tokens" json:"max_tokens"`           // 最大 Token 数量
}

type DatabaseConfig struct {
	Type     string `yaml:"type" json:"type"`         // 数据库类型: sqlite, mysql, postgres
	Path     string `yaml:"path" json:"path"`         // SQLite 数据库文件路径
	Host     string `yaml:"host" json:"host"`         // 数据库地址
	Port     int    `yaml:"port" json:"port"`         // 数据库端口
	User     string `yaml:"user" json:"user"`         // 数据库用户名
	Password string `yaml:"password" json:"password"` // 数据库密码
	DBName   string `yaml:"dbname" json:"dbname"`     // 数据库名称
	Config   string `yaml:"config" json:"config"`     // 额外配置参数 (DSN)
}

var (
	GlobalConfig *Config
	ConfigPath   string
	configLock   sync.RWMutex
)

func LoadConfig(path string) error {
	configLock.Lock()
	defer configLock.Unlock()

	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 创建默认配置
		cfg := GetDefaultConfig()
		data, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}
		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			return err
		}
		cfg.overrideWithEnv()
		GlobalConfig = cfg
		ConfigPath = path
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	cfg.overrideWithEnv()
	GlobalConfig = &cfg
	ConfigPath = path
	return nil
}

func (cfg *Config) overrideWithEnv() {
	// Site settings
	if val := os.Getenv("FR_SITE_NAME"); val != "" {
		cfg.Site.Name = val
	}
	if val := os.Getenv("FR_SITE_BASE_URL"); val != "" {
		cfg.Site.BaseURL = val
	}
	if val := os.Getenv("FR_SITE_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			cfg.Site.Port = port
		}
	}

	// Security settings
	if val := os.Getenv("FR_SECURITY_JWT_SECRET"); val != "" {
		cfg.Security.JWTSecret = val
	}

	// Upload settings
	if val := os.Getenv("FR_UPLOAD_MAX_SIZE"); val != "" {
		if size, err := strconv.ParseInt(val, 10, 64); err == nil {
			cfg.Upload.MaxFileSizeMB = size
		}
	}
	if val := os.Getenv("FR_UPLOAD_RETENTION_DAYS"); val != "" {
		if days, err := strconv.Atoi(val); err == nil {
			cfg.Upload.MaxRetentionDays = days
		}
	}

	// Database settings
	if val := os.Getenv("FR_DB_TYPE"); val != "" {
		cfg.Database.Type = val
	}
	if val := os.Getenv("FR_DB_PATH"); val != "" {
		cfg.Database.Path = val
	}
	if val := os.Getenv("FR_DB_HOST"); val != "" {
		cfg.Database.Host = val
	}
	if val := os.Getenv("FR_DB_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			cfg.Database.Port = port
		}
	}
	if val := os.Getenv("FR_DB_USER"); val != "" {
		cfg.Database.User = val
	}
	if val := os.Getenv("FR_DB_PASSWORD"); val != "" {
		cfg.Database.Password = val
	}
	if val := os.Getenv("FR_DB_NAME"); val != "" {
		cfg.Database.DBName = val
	}

	// Storage settings
	if val := os.Getenv("FR_STORAGE_TYPE"); val != "" {
		cfg.Storage.Type = val
	}
	if val := os.Getenv("FR_STORAGE_LOCAL_PATH"); val != "" {
		cfg.Storage.Local.Path = val
	}

	// Log settings
	if val := os.Getenv("FR_LOG_LEVEL"); val != "" {
		cfg.Log.Level = val
	}
	if val := os.Getenv("FR_LOG_FILE_PATH"); val != "" {
		cfg.Log.FilePath = val
	}

	// Web settings
	if val := os.Getenv("FR_WEB_PATH"); val != "" {
		cfg.Web.Path = val
	}
}

func GetDefaultConfig() *Config {
	return &Config{
		Site: SiteConfig{
			Name:        "文件暂存柜",
			Description: "临时文件中转服务",
			Logo:        "/favicon.png",
			BaseURL:     "",
			Port:        8080,
		},
		Security: SecurityConfig{
			AdminPasswordHash: "$2a$10$Bm0TEmU4uj.bVHYiIPFBheUkcdg6XHpsanLvmpoAtgU1UnKbo9.vy", // 默认密码: admin123
			PickupCodeLength:  6,
			PickupFailLimit:   5,
			JWTSecret:         "file-relay-secret",
		},
		Upload: UploadConfig{
			MaxFileSizeMB:    100,
			MaxBatchFiles:    20,
			MaxRetentionDays: 30,
			RequireToken:     false,
		},
		Storage: StorageConfig{
			Type: "local",
			Local: struct {
				Path string `yaml:"path" json:"path"`
			}{
				Path: "data/storage_data",
			},
		},
		APIToken: APITokenConfig{
			Enabled:       true,
			AllowAdminAPI: true,
			MaxTokens:     20,
		},
		Database: DatabaseConfig{
			Type: "sqlite",
			Path: "data/file_relay.db",
		},
		Web: WebConfig{
			Path: "web",
		},
		Log: LogConfig{
			Level:    "info",
			FilePath: "data/logs/app.log",
		},
	}
}

func SaveConfig() error {
	configLock.RLock()
	defer configLock.RUnlock()

	data, err := yaml.Marshal(GlobalConfig)
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigPath, data, 0644)
}

func UpdateGlobalConfig(cfg *Config) {
	configLock.Lock()
	defer configLock.Unlock()
	GlobalConfig = cfg
}
