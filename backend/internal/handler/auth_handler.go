package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"zblog-backend/internal/dto"
	"zblog-backend/internal/middleware"
	"zblog-backend/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, dto.BadRequest("username and password are required"))
		return
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusOK, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"token": token,
	}))
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.authService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusOK, dto.NotFound("user not found"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(user))
}
