package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefaultConfig(t *testing.T) {
	cfg := GetDefaultConfig()
	assert.Equal(t, "data/storage_data", cfg.Storage.Local.Path)
	assert.Equal(t, "data/file_relay.db", cfg.Database.Path)
	assert.Equal(t, "data/logs/app.log", cfg.Log.FilePath)
}

func TestOverrideWithEnv(t *testing.T) {
	// 设置环境变量
	os.Setenv("FR_SITE_NAME", "EnvSiteName")
	os.Setenv("FR_SITE_PORT", "9999")
	os.Setenv("FR_DB_TYPE", "mysql")
	os.Setenv("FR_DB_PATH", "custom_db_path.db")
	os.Setenv("FR_UPLOAD_MAX_SIZE", "500")
	os.Setenv("FR_LOG_FILE_PATH", "custom_log_path.log")

	cfg := GetDefaultConfig()
	cfg.overrideWithEnv()

	assert.Equal(t, "EnvSiteName", cfg.Site.Name)
	assert.Equal(t, 9999, cfg.Site.Port)
	assert.Equal(t, "mysql", cfg.Database.Type)
	assert.Equal(t, "custom_db_path.db", cfg.Database.Path)
	assert.Equal(t, int64(500), cfg.Upload.MaxFileSizeMB)
	assert.Equal(t, "custom_log_path.log", cfg.Log.FilePath)

	// 清理环境变量
	os.Unsetenv("FR_SITE_NAME")
	os.Unsetenv("FR_SITE_PORT")
	os.Unsetenv("FR_DB_TYPE")
	os.Unsetenv("FR_DB_PATH")
	os.Unsetenv("FR_UPLOAD_MAX_SIZE")
	os.Unsetenv("FR_LOG_FILE_PATH")
}
