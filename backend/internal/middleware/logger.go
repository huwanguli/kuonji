package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method

		entry := logrus.WithFields(logrus.Fields{
			"status":    status,
			"method":    method,
			"path":      path,
			"query":     query,
			"ip":        c.ClientIP(),
			"latency":   latency.String(),
			"userAgent": c.Request.UserAgent(),
		})

		if status >= 500 {
			entry.Error("request completed")
		} else if status >= 400 {
			entry.Warn("request completed")
		} else {
			entry.Info("request completed")
		}
	}
}
