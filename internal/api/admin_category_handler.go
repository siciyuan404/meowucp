package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type AdminCategoryService interface {
	CreateCategory(category *domain.Category) error
	UpdateCategory(category *domain.Category) error
	GetCategory(id int64) (*domain.Category, error)
	ListCategories(offset, limit int) ([]*domain.Category, int64, error)
}

type AdminCategoryHandler struct {
	service AdminCategoryService
}

func NewAdminCategoryHandler(service AdminCategoryService) *AdminCategoryHandler {
	return &AdminCategoryHandler{service: service}
}

type AdminCategoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	ParentID    *int64 `json:"parent_id"`
	SortOrder   int    `json:"sort_order"`
	Status      int    `json:"status"`
}

func (h *AdminCategoryHandler) Create(c *gin.Context) {
	var req AdminCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if req.Name == "" || req.Slug == "" {
		respondError(c, http.StatusBadRequest, "missing_required_fields", "Name and slug are required")
		return
	}
	category := &domain.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    req.ParentID,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	}
	if err := h.service.CreateCategory(category); err != nil {
		respondError(c, http.StatusInternalServerError, "create_failed", "Failed to create category")
		return
	}
	c.JSON(http.StatusCreated, category)
}

func (h *AdminCategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_id", "Invalid category id")
		return
	}
	var req AdminCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	category, err := h.service.GetCategory(id)
	if err != nil {
		respondError(c, http.StatusNotFound, "category_not_found", "Category not found")
		return
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.SortOrder != 0 {
		category.SortOrder = req.SortOrder
	}
	if req.Status != 0 {
		category.Status = req.Status
	}

	if err := h.service.UpdateCategory(category); err != nil {
		respondError(c, http.StatusInternalServerError, "update_failed", "Failed to update category")
		return
	}
	c.JSON(http.StatusOK, category)
}

func (h *AdminCategoryHandler) List(c *gin.Context) {
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
	items, total, err := h.service.ListCategories(offset, limitInt)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "list_failed", "Failed to list categories")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"categories": items,
		"pagination": gin.H{
			"page":  pageInt,
			"limit": limitInt,
			"total": total,
		},
	})
}
