package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
)

type CurrencyService interface {
	Convert(amount float64, base, target string) (float64, error)
}

type PublicProductService interface {
	ListProducts(offset, limit int, filters map[string]interface{}) ([]*domain.Product, int64, error)
	GetProduct(id int64) (*domain.Product, error)
}

type ProductHandler struct {
	service         PublicProductService
	currencyService CurrencyService
}

func NewProductHandler(service PublicProductService, currencyService CurrencyService) *ProductHandler {
	return &ProductHandler{
		service:         service,
		currencyService: currencyService,
	}
}

func (h *ProductHandler) List(c *gin.Context) {
	currency := c.DefaultQuery("currency", "CNY")
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
	filters := map[string]interface{}{}
	filters["status = ?"] = 1

	products, total, err := h.service.ListProducts(offset, limitInt, filters)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "list_failed", "Failed to list products")
		return
	}

	convertedProducts := make([]gin.H, len(products))
	for i, product := range products {
		convertedPrice, _ := h.currencyService.Convert(product.Price, "CNY", currency)
		convertedComparePrice, _ := h.currencyService.Convert(product.ComparePrice, "CNY", currency)
		convertedProducts[i] = gin.H{
			"id":             product.ID,
			"name":           product.Name,
			"slug":           product.Slug,
			"description":    product.Description,
			"price":          convertedPrice,
			"compare_price":  convertedComparePrice,
			"sku":            product.SKU,
			"stock_quantity": product.StockQuantity,
			"category_id":    product.CategoryID,
			"images":         product.Images,
			"currency":       currency,
			"locale":         locale,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"products": convertedProducts,
		"pagination": gin.H{
			"page":  pageInt,
			"limit": limitInt,
			"total": total,
		},
	})
}

func (h *ProductHandler) Get(c *gin.Context) {
	currency := c.DefaultQuery("currency", "CNY")
	locale := c.DefaultQuery("locale", "zh-CN")

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_id", "Invalid product id")
		return
	}

	product, err := h.service.GetProduct(id)
	if err != nil {
		respondError(c, http.StatusNotFound, "product_not_found", "Product not found")
		return
	}

	convertedPrice, _ := h.currencyService.Convert(product.Price, "CNY", currency)
	convertedComparePrice, _ := h.currencyService.Convert(product.ComparePrice, "CNY", currency)

	c.JSON(http.StatusOK, gin.H{
		"id":             product.ID,
		"name":           product.Name,
		"slug":           product.Slug,
		"description":    product.Description,
		"price":          convertedPrice,
		"compare_price":  convertedComparePrice,
		"sku":            product.SKU,
		"stock_quantity": product.StockQuantity,
		"category_id":    product.CategoryID,
		"images":         product.Images,
		"currency":       currency,
		"locale":         locale,
	})
}
