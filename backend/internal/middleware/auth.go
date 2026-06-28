package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"zblog-backend/internal/config"
	"zblog-backend/internal/utils"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "authorization header required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "authorization header format must be Bearer {token}",
			})
			return
		}

		claims, err := utils.ParseToken(parts[1], config.Cfg.JWT.Secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "invalid or expired token",
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func GetUserID(c *gin.Context) uint {
	id, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return id.(uint)
}
