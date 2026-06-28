package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)
	return path
}

func TestLoad(t *testing.T) {
	yaml := `
server:
  port: 8080
  mode: debug

database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: secret
  dbname: zblog
  charset: utf8mb4
  max_idle_conns: 10
  max_open_conns: 100

jwt:
  secret: test-secret
  expire: 24h

log:
  level: debug
  file: logs/zblog.log
  max_age: 30
  max_size: 100

upload:
  path: uploads/
  max_size: 10
  allowed_exts: [".jpg", ".png"]
`
	path := writeTempConfig(t, yaml)
	cfg, err := Load(path)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "debug", cfg.Server.Mode)
	assert.Equal(t, "127.0.0.1", cfg.Database.Host)
	assert.Equal(t, 3306, cfg.Database.Port)
	assert.Equal(t, "root", cfg.Database.User)
	assert.Equal(t, "secret", cfg.Database.Password)
	assert.Equal(t, "zblog", cfg.Database.DBName)
	assert.Equal(t, "utf8mb4", cfg.Database.Charset)
	assert.Equal(t, 10, cfg.Database.MaxIdleConns)
	assert.Equal(t, 100, cfg.Database.MaxOpenConns)
	assert.Equal(t, "test-secret", cfg.JWT.Secret)
	assert.Greater(t, cfg.JWT.Expire, int64(0))
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "logs/zblog.log", cfg.Log.File)
	assert.Equal(t, 30, cfg.Log.MaxAge)
	assert.Equal(t, 100, cfg.Log.MaxSize)
	assert.Equal(t, "uploads/", cfg.Upload.Path)
	assert.Equal(t, 10, cfg.Upload.MaxSize)
	assert.Equal(t, []string{".jpg", ".png"}, cfg.Upload.AllowedExts)
}

func TestLoadFileNotFound(t *testing.T) {
	_, err := Load("nonexistent.yaml")
	assert.Error(t, err)
}

func TestLoadInvalidYAML(t *testing.T) {
	path := writeTempConfig(t, "invalid: [yaml: broken")
	_, err := Load(path)
	assert.Error(t, err)
}
