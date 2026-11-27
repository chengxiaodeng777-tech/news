package system

import (
	systemmodel "base/internal/model/system"
	"base/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (h *Handler) ListArticles(c *gin.Context) {
	var q systemmodel.PageQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Error(c, 400, "invalid query")
		return
	}
	v, err := h.service.ListArticles(c.Request.Context(), q)
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, v)
}

func (h *Handler) ListPublishedArticles(c *gin.Context) {
	// Publication status is enforced by the service; category is an exact match.
	var q struct {
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		Keyword  string `form:"keyword"`
		Category string `form:"category" binding:"max=80"`
	}
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Error(c, 400, "invalid query")
		return
	}
	v, err := h.service.ListPublishedArticles(c.Request.Context(), systemmodel.PageQuery{
		Page: q.Page, PageSize: q.PageSize, Keyword: q.Keyword, Category: q.Category,
	})
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, v)
}

func (h *Handler) GetPublishedArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.Error(c, 400, "invalid id")
		return
	}
	v, err := h.service.GetPublishedArticle(c.Request.Context(), id)
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, v)
}
func (h *Handler) GetArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "invalid id")
		return
	}
	v, err := h.service.GetArticle(c.Request.Context(), id)
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, v)
}
func (h *Handler) CreateArticle(c *gin.Context) {
	var req systemmodel.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid body")
		return
	}
	if err := h.service.CreateArticle(c.Request.Context(), req); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}
func (h *Handler) UpdateArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "invalid id")
		return
	}
	var req systemmodel.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid body")
		return
	}
	if err := h.service.UpdateArticle(c.Request.Context(), id, req); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}
func (h *Handler) DeleteArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "invalid id")
		return
	}
	if err := h.service.DeleteArticle(c.Request.Context(), id); err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, true)
}
