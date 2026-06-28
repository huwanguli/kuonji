package handler

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"zblog-backend/internal/config"
	"zblog-backend/internal/dto"
	"zblog-backend/internal/utils"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, dto.BadRequest("file is required"))
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allow := false
	for _, e := range config.Cfg.Upload.AllowedExts {
		if e == ext {
			allow = true
			break
		}
	}
	if !allow {
		c.JSON(http.StatusOK, dto.BadRequest(fmt.Sprintf("file extension %s is not allowed", ext)))
		return
	}

	if file.Size > int64(config.Cfg.Upload.MaxSize)*1024*1024 {
		c.JSON(http.StatusOK, dto.BadRequest(fmt.Sprintf("file size exceeds %dMB", config.Cfg.Upload.MaxSize)))
		return
	}

	today := time.Now().Format("2006-01-02")
	subDir := filepath.Join(config.Cfg.Upload.Path, today)
	if err := os.MkdirAll(subDir, 0755); err != nil {
		logrus.WithError(err).Error("create upload directory failed")
		c.JSON(http.StatusOK, dto.InternalError("upload failed"))
		return
	}

	newName := utils.GenerateSlug(strings.TrimSuffix(file.Filename, ext)) + ext

	if err := c.SaveUploadedFile(file, filepath.Join(subDir, newName)); err != nil {
		logrus.WithError(err).Error("upload file failed")
		c.JSON(http.StatusOK, dto.InternalError("upload failed"))
		return
	}

	urlPath := path.Join("/", config.Cfg.Upload.Path, today)
	url := urlPath + "/" + newName
	c.JSON(http.StatusOK, dto.Success(gin.H{
		"url": url,
	}))
}
