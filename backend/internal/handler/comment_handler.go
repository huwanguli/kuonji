package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"zblog-backend/internal/dto"
	"zblog-backend/internal/service"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (h *CommentHandler) Create(c *gin.Context) {
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, dto.BadRequest(err.Error()))
		return
	}

	comment, err := h.commentService.Create(&req, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusOK, dto.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(comment))
}

func (h *CommentHandler) GetList(c *gin.Context) {
	var query dto.CommentListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusOK, dto.BadRequest(err.Error()))
		return
	}
	if query.ArticleID == 0 {
		c.JSON(http.StatusOK, dto.BadRequest("article_id is required"))
		return
	}

	comments, total, err := h.commentService.GetList(&query)
	if err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(dto.CommentListResponse{
		List:  comments,
		Total: total,
	}))
}

func (h *CommentHandler) GetAdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	var status *int
	if s, err := strconv.Atoi(c.Query("status")); err == nil {
		status = &s
	}

	comments, total, err := h.commentService.GetAllForAdmin(page, pageSize, status)
	if err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.PageResult(comments, total, page, pageSize))
}

func (h *CommentHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.BadRequest("invalid id"))
		return
	}

	var req dto.UpdateCommentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, dto.BadRequest(err.Error()))
		return
	}

	if err := h.commentService.UpdateStatus(uint(id), req.Status); err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessMessage("comment status updated"))
}

func (h *CommentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.BadRequest("invalid id"))
		return
	}

	if err := h.commentService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessMessage("comment deleted"))
}
