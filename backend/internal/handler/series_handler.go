package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"zblog-backend/internal/dto"
	"zblog-backend/internal/service"
)

type SeriesHandler struct {
	seriesService service.SeriesService
}

func NewSeriesHandler(seriesService service.SeriesService) *SeriesHandler {
	return &SeriesHandler{seriesService: seriesService}
}

func (h *SeriesHandler) Create(c *gin.Context) {
	var req dto.CreateSeriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, dto.BadRequest("请输入系列名称"))
		return
	}

	series, err := h.seriesService.Create(req.Name, req.Cover, req.Description)
	if err != nil {
		c.JSON(http.StatusOK, dto.BadRequest(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(series))
}

func (h *SeriesHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.BadRequest("无效 ID"))
		return
	}

	var req dto.CreateSeriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, dto.BadRequest("请输入系列名称"))
		return
	}

	series, err := h.seriesService.Update(uint(id), req.Name, req.Cover, req.Description)
	if err != nil {
		c.JSON(http.StatusOK, dto.BadRequest(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(series))
}

func (h *SeriesHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, dto.BadRequest("无效 ID"))
		return
	}

	if err := h.seriesService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessMessage("已删除"))
}

func (h *SeriesHandler) GetAll(c *gin.Context) {
	list, err := h.seriesService.GetAllWithCount()
	if err != nil {
		c.JSON(http.StatusOK, dto.InternalError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(list))
}

func (h *SeriesHandler) GetDetail(c *gin.Context) {
	name := c.Param("name")
	detail, err := h.seriesService.GetDetail(name)
	if err != nil {
		c.JSON(http.StatusOK, dto.NotFound("系列不存在"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(detail))
}
