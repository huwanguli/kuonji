package logger

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "test.log")
	Init("debug", logPath, 7, 10)
	defer Close()

	assert.Equal(t, logrus.DebugLevel, logrus.GetLevel())

	logrus.Info("test info message")

	info, err := os.Stat(logPath)
	assert.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))
}

func TestInitDefaultLevel(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "test.log")
	Init("invalid", logPath, 7, 10)
	defer Close()

	assert.Equal(t, logrus.InfoLevel, logrus.GetLevel())
}

func TestInitFileOutput(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "test.log")
	Init("debug", logPath, 7, 10)
	defer Close()

	logrus.Debug("debug output")

	data, err := os.ReadFile(logPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "debug output")
}
