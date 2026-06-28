package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"zblog-backend/internal/dto"
	"zblog-backend/internal/service"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.categoryService.GetAll()
	if err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(categories))
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, dto.BadRequest(err.Error()))
		return
	}

	cat, err := h.categoryService.Create(req.Name, req.Slug, req.Description)
	if err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(cat))
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.BadRequest("invalid id"))
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, dto.BadRequest(err.Error()))
		return
	}

	cat, err := h.categoryService.Update(uint(id), req.Name, req.Slug, req.Description)
	if err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(cat))
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.BadRequest("invalid id"))
		return
	}

	if err := h.categoryService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessMessage("category deleted"))
}
