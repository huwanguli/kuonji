package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"zblog-backend/internal/dto"
	"zblog-backend/internal/service"
)

type TagHandler struct {
	tagService service.TagService
}

func NewTagHandler(tagService service.TagService) *TagHandler {
	return &TagHandler{tagService: tagService}
}

func (h *TagHandler) GetAll(c *gin.Context) {
	tags, err := h.tagService.GetAll()
	if err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(tags))
}

func (h *TagHandler) Create(c *gin.Context) {
	var req dto.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, dto.BadRequest(err.Error()))
		return
	}

	tag, err := h.tagService.Create(req.Name, req.Slug)
	if err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(tag))
}

func (h *TagHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.BadRequest("invalid id"))
		return
	}

	var req dto.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, dto.BadRequest(err.Error()))
		return
	}

	tag, err := h.tagService.Update(uint(id), req.Name, req.Slug)
	if err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(tag))
}

func (h *TagHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.BadRequest("invalid id"))
		return
	}

	if err := h.tagService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessMessage("tag deleted"))
}
