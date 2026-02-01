package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddressValidationHandler struct{}

func NewAddressValidationHandler() *AddressValidationHandler {
	return &AddressValidationHandler{}
}

func (h *AddressValidationHandler) Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"valid": true})
}
