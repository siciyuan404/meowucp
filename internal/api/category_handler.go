package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type PublicCategoryService interface {
	ListCategories(offset, limit int) ([]*domain.Category, int64, error)
	GetCategory(id int64) (*domain.Category, error)
}

type CategoryHandler struct {
	service PublicCategoryService
}

func NewCategoryHandler(service PublicCategoryService, _ interface{}) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

func (h *CategoryHandler) List(c *gin.Context) {
	locale := c.DefaultQuery("locale", "zh-CN")
	limit := c.DefaultQuery("limit", "20")
	page := c.DefaultQuery("page", "1")

	limitInt := parseInt(limit)
	pageInt := parseInt(page)
	if pageInt < 1 {
		pageInt = 1
	}
	if limitInt < 1 {
		limitInt = 20
	}

	offset := (pageInt - 1) * limitInt

	categories, total, err := h.service.ListCategories(offset, limitInt)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "list_failed", "Failed to list categories")
		return
	}

	convertedCategories := make([]gin.H, len(categories))
	for i, category := range categories {
		convertedCategories[i] = gin.H{
			"id":          category.ID,
			"name":        category.Name,
			"slug":        category.Slug,
			"description": category.Description,
			"parent_id":   category.ParentID,
			"sort_order":  category.SortOrder,
			"locale":      locale,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": convertedCategories,
		"pagination": gin.H{
			"page":  pageInt,
			"limit": limitInt,
			"total": total,
		},
	})
}

func (h *CategoryHandler) Get(c *gin.Context) {
	locale := c.DefaultQuery("locale", "zh-CN")

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_id", "Invalid category id")
		return
	}

	category, err := h.service.GetCategory(id)
	if err != nil {
		respondError(c, http.StatusNotFound, "category_not_found", "Category not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          category.ID,
		"name":        category.Name,
		"slug":        category.Slug,
		"description": category.Description,
		"parent_id":   category.ParentID,
		"sort_order":  category.SortOrder,
		"locale":      locale,
	})
}
