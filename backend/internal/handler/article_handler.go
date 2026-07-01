package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"zblog-backend/internal/dto"
	"zblog-backend/internal/model"
	"zblog-backend/internal/service"
)

func translateBindErr(err error) string {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "Title"):
		return "请填写文章标题"
	case strings.Contains(msg, "ContentMD"):
		return "请填写文章正文"
	default:
		return "请求参数有误"
	}
}

type ArticleHandler struct {
	articleService service.ArticleService
}

func NewArticleHandler(articleService service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: articleService}
}

func (h *ArticleHandler) GetPublicList(c *gin.Context) {
	var query dto.ArticleListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusOK, dto.BadRequest(err.Error()))
		return
	}
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	published := 1
	query.Status = &published

	articles, total, err := h.articleService.GetList(&query)
	if err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.PageResult(articles, total, query.Page, query.PageSize))
}

func (h *ArticleHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	detail, err := h.articleService.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusOK, dto.NotFound("article not found"))
		return
	}
	if a, ok := detail.Article.(*model.Article); !ok || a.Status != 1 {
		c.JSON(http.StatusOK, dto.NotFound("article not found"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(detail))
}

func (h *ArticleHandler) GetAdminDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.BadRequest("invalid id"))
		return
	}

	article, err := h.articleService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, dto.NotFound("article not found"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(article))
}

func (h *ArticleHandler) GetSeriesList(c *gin.Context) {
	series, err := h.articleService.GetSeriesList()
	if err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(series))
}

func (h *ArticleHandler) GetAnnouncements(c *gin.Context) {
	announcements, err := h.articleService.GetAnnouncements()
	if err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(announcements))
}

func (h *ArticleHandler) GetAdminList(c *gin.Context) {
	var query dto.ArticleListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusOK, dto.BadRequest(err.Error()))
		return
	}
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	articles, total, err := h.articleService.GetList(&query)
	if err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.PageResult(articles, total, query.Page, query.PageSize))
}

func (h *ArticleHandler) Create(c *gin.Context) {
	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, dto.BadRequest(translateBindErr(err)))
		return
	}

	article, err := h.articleService.Create(&req)
	if err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(article))
}

func (h *ArticleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.BadRequest("invalid id"))
		return
	}

	var req dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, dto.BadRequest(translateBindErr(err)))
		return
	}

	article, err := h.articleService.Update(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(article))
}

func (h *ArticleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.BadRequest("invalid id"))
		return
	}

	if err := h.articleService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessMessage("article deleted"))
}
