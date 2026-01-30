package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type AdminProductService interface {
	CreateProduct(product *domain.Product) error
	GetProduct(id int64) (*domain.Product, error)
	UpdateProduct(product *domain.Product) error
	ListProducts(offset, limit int, filters map[string]interface{}) ([]*domain.Product, int64, error)
}

type AdminProductHandler struct {
	service AdminProductService
}

func NewAdminProductHandler(service AdminProductService) *AdminProductHandler {
	return &AdminProductHandler{service: service}
}

type AdminProductRequest struct {
	Name          string  `json:"name"`
	Slug          string  `json:"slug"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	ComparePrice  float64 `json:"compare_price"`
	SKU           string  `json:"sku"`
	StockQuantity int     `json:"stock_quantity"`
	Status        int     `json:"status"`
	CategoryID    *int64  `json:"category_id"`
	Images        string  `json:"images"`
}

type AdminProductStatusRequest struct {
	Status int `json:"status"`
}

func (h *AdminProductHandler) Create(c *gin.Context) {
	var req AdminProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if req.Name == "" || req.Slug == "" || req.SKU == "" || req.Price <= 0 {
		respondError(c, http.StatusBadRequest, "missing_required_fields", "Name, slug, sku and price are required")
		return
	}

	product := &domain.Product{
		Name:          req.Name,
		Slug:          req.Slug,
		Description:   req.Description,
		Price:         req.Price,
		ComparePrice:  req.ComparePrice,
		SKU:           req.SKU,
		StockQuantity: req.StockQuantity,
		Status:        req.Status,
		CategoryID:    req.CategoryID,
		Images:        req.Images,
	}

	if err := h.service.CreateProduct(product); err != nil {
		respondError(c, http.StatusInternalServerError, "create_failed", "Failed to create product")
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *AdminProductHandler) List(c *gin.Context) {
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
	filters := map[string]interface{}{}
	if status := c.Query("status"); status != "" {
		filters["status = ?"] = status
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		filters["category_id = ?"] = categoryID
	}
	if sku := c.Query("sku"); sku != "" {
		filters["sku = ?"] = sku
	}

	products, total, err := h.service.ListProducts(offset, limitInt, filters)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "list_failed", "Failed to list products")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"pagination": gin.H{
			"page":  pageInt,
			"limit": limitInt,
			"total": total,
		},
	})
}

func (h *AdminProductHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_id", "Invalid product id")
		return
	}
	product, err := h.service.GetProduct(id)
	if err != nil {
		respondError(c, http.StatusNotFound, "product_not_found", "Product not found")
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *AdminProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_id", "Invalid product id")
		return
	}
	var req AdminProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	product, err := h.service.GetProduct(id)
	if err != nil {
		respondError(c, http.StatusNotFound, "product_not_found", "Product not found")
		return
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Slug != "" {
		product.Slug = req.Slug
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.ComparePrice > 0 {
		product.ComparePrice = req.ComparePrice
	}
	if req.SKU != "" {
		product.SKU = req.SKU
	}
	if req.StockQuantity >= 0 {
		product.StockQuantity = req.StockQuantity
	}
	if req.Status != 0 {
		product.Status = req.Status
	}
	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}
	if req.Images != "" {
		product.Images = req.Images
	}

	if err := h.service.UpdateProduct(product); err != nil {
		respondError(c, http.StatusInternalServerError, "update_failed", "Failed to update product")
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *AdminProductHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_id", "Invalid product id")
		return
	}
	var req AdminProductStatusRequest
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if req.Status == 0 {
		respondError(c, http.StatusBadRequest, "missing_status", "Status is required")
		return
	}

	product, err := h.service.GetProduct(id)
	if err != nil {
		respondError(c, http.StatusNotFound, "product_not_found", "Product not found")
		return
	}

	product.Status = req.Status
	if err := h.service.UpdateProduct(product); err != nil {
		respondError(c, http.StatusInternalServerError, "update_failed", "Failed to update product")
		return
	}

	c.JSON(http.StatusOK, product)
}
